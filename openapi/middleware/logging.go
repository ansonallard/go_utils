package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

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

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Interface("headers", c.Request.Header).
			RawJSON("requestBody", requestBody).
			Msg("API Request")

		c.Next()

		log.Info().
			Int("status", c.Writer.Status()).
			TimeDiff("latency", time.Now().UTC(), startTime).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Msg("API Response")
	}
}
