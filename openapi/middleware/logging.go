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

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := zerolog.Ctx(c.Request.Context())
		startTime := time.Now().UTC()

		// Read and restore request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap the response writer to capture body
		rbw := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = rbw

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
			Interface("headers", c.Writer.Header()).
			RawJSON("responseBody", rbw.body.Bytes()).
			Msg("API Response")
	}
}
