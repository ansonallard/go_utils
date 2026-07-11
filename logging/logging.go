package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func ZeroLogConfiguration(logFile *os.File, logLevel *zerolog.Level, serviceName, serviceVersion string) context.Context {
	if logLevel != nil {
		zerolog.SetGlobalLevel(*logLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var writer io.Writer
	if logFile != nil {
		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = os.Stdout
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(writer).With().
		Timestamp().
		Str("serviceName", serviceName).
		Str("serviceVersion", serviceVersion).
		Logger()

	ctx := context.Background()

	// Attach the Logger to the context.Context
	ctx = logger.WithContext(ctx)
	return ctx
}

func InjectLogger(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := log.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
	}
}

// LoggingMiddleware - remove body/headers from response log
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := zerolog.Ctx(c.Request.Context())
		startTime := time.Now().UTC()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		var incomingLog *zerolog.Event
		if len(requestBody) > 0 && json.Valid(requestBody) {
			var compacted bytes.Buffer
			if err := json.Compact(&compacted, requestBody); err == nil {
				incomingLog = log.Info().RawJSON("requestBody", compacted.Bytes())
			} else {
				incomingLog = log.Info().Str("requestBody", string(requestBody))
			}
		} else if len(requestBody) > 0 {
			incomingLog = log.Info().Str("requestBody", string(requestBody))
		} else {
			incomingLog = log.Info()
		}

		incomingLog.
			Str("method", c.Request.Method).
			Str("pattern", c.Request.Pattern).
			Str("path", c.Request.URL.Path).
			Interface("headers", c.Request.Header).
			Msg("API Request")

		c.Next()

		log.Info().
			Int("status", c.Writer.Status()).
			TimeDiff("latency", time.Now().UTC(), startTime).
			Str("method", c.Request.Method).
			Str("pattern", c.Request.Pattern).
			Str("path", c.Request.URL.Path).
			Msg("API Response")
	}
}

func RecoveryMiddleware(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Interface("panic", r).
					Str("path", c.Request.URL.Path).
					Str("pattern", c.Request.Pattern).
					Str("method", c.Request.Method).
					Msg("panic recovered")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
