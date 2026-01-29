package openapi

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/ansonallard/go_utils/openapi/ierr"
	"github.com/ansonallard/go_utils/openapi/middleware/openapi"
	"github.com/ansonallard/go_utils/openapi/request"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	defaultIPv4OpenAddress = "0.0.0.0"
)

func ServeOpenAPI(ctx context.Context, config serveOpenAPIConfig) error {
	log := zerolog.Ctx(ctx)

	loader := openapi3.NewLoader()

	var openAPISpec *openapi3.T
	var err error
	if config.openAPISpecFilePath != "" {
		openAPISpec, err = loader.LoadFromFile(config.openAPISpecFilePath)
	} else if config.openAPISpecData != nil {
		openAPISpec, err = loader.LoadFromData(config.openAPISpecData)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error loading OpenAPI spec")
		return err
	}

	// Validate the OpenAPI spec itself
	err = openAPISpec.Validate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error validating swagger spec")
		return err
	}

	// Create router from OpenAPI spec
	router, err := gorillamux.NewRouter(openAPISpec)
	if err != nil {
		log.Error().Err(err).Msg("Error creating router")
		return err
	}

	// Create Gin router
	ginMode := gin.DebugMode
	if !config.isDevMode {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())

	if config.authZMiddleware != nil {
		ginRouter.Use(openapi.ValidationMiddleware(router, config.authZMiddleware.AuthorizeCaller))
	}

	// Validate that top level struct contains all required OpenAPI operation IDs
	if err = openapi.ValidateStructAndOpenAPI(openAPISpec, config.serviceController); err != nil {
		log.Error().Err(err).Msg("Failed to ValidateStructAndOpenAPI")
		return err
	}

	ginRouter.Any("/*path", func(c *gin.Context) {
		log.Info().Interface("request", c.Request).Msg("API Request")
		startTime := time.Now().UTC()

		route, pathParams, err := router.FindRoute(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error finding route: %v", err)})
			return
		}
		log.Info().Interface("route", route).Interface("pathParams", pathParams).Msg("Route and path params")

		firstSuccessfulResponseCode, err := openapi.GetFirstSuccessfulStatusCode(route.Operation.Responses)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		topLevelStructReflected := reflect.ValueOf(config.serviceController)
		method := topLevelStructReflected.MethodByName(openapi.ConvertOperationIdToPascalCase(route.Operation.OperationID))

		iRequest := request.NewRequest(&request.RequestConfig{
			QueryParams: c.Request.URL.Query(),
			Headers:     c.Request.Header,
			PathParams:  pathParams,
			RequestBody: c.Request.Body,
		})

		values := []reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(iRequest)}
		result := method.Call(values)

		// All top level methods must either return an error
		// or a successful response and error
		var methodResult any
		switch len(result) {
		case 1:
			err, ok := result[0].Interface().(error)
			if ok {
				errorHandler(ctx, err, c)
				return
			}
		case 2:
			methodResult = result[0].Interface()
			err, ok := result[1].Interface().(error)
			if ok {
				errorHandler(ctx, err, c)
				return
			}
		}

		log.Info().Int("status", firstSuccessfulResponseCode).
			Interface("response", methodResult).TimeDiff("latency", time.Now().UTC(), startTime).
			Str("httpMethod", route.Method).Str("path", route.Path).
			Msg("API Response")
		c.JSON(firstSuccessfulResponseCode, methodResult)
	})

	port := config.port
	log.Info().Uint16("port", port).Msgf("Server starting on :%d", port)
	if err := ginRouter.Run(fmt.Sprintf("%s:%d", defaultIPv4OpenAddress, port)); err != nil {
		log.Error().Err(err).Msg("Failed to run gin router")
		return err
	}
	return nil
}

func errorHandler(ctx context.Context, err error, c *gin.Context) {
	switch err.(type) {
	case *ierr.UnAuthorizedError:
		abortWithStatusResponse(ctx, http.StatusUnauthorized, err, c)
	case *ierr.NotFoundError:
		abortWithStatusResponse(ctx, http.StatusNotFound, err, c)
	case *ierr.ConflictError:
		abortWithStatusResponse(ctx, http.StatusConflict, err, c)
	default:
		abortWithStatusResponse(ctx, http.StatusInternalServerError, err, c)
	}
}

func abortWithStatusResponse(ctx context.Context, code int, err error, c *gin.Context) {
	log := zerolog.Ctx(ctx)
	log.Warn().Err(err).Int("status", code).Interface("request", c.Request).Msg("API Response Error")

	c.AbortWithStatusJSON(code, map[string]string{"message": err.Error()})
}
