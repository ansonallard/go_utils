package middleware

import (
	"bytes"
	"encoding/json"
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

		var incomingLog *zerolog.Event
		if len(requestBody) > 0 && json.Valid(requestBody) {
			incomingLog = log.Info().RawJSON("requestBody", requestBody)
		} else if len(requestBody) > 0 {
			incomingLog = log.Info().Str("requestBody", string(requestBody))
		} else {
			incomingLog = log.Info()
		}

		incomingLog.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Interface("headers", c.Request.Header).
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
