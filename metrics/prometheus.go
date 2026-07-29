package metrics

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

// Run Prometheus Metrics Server in a go-routine (non-blocking)
func StartMetricsServer(ctx context.Context) {
	go func(ctx context.Context) {
		log := zerolog.Ctx(ctx)

		metricsRouter := gin.New()
		metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
		var metricsPort uint16 = 9101

		log.Info().Uint16("port", metricsPort).Msgf("Server starting on :%d", metricsPort)
		if err := metricsRouter.Run(fmt.Sprintf(":%d", metricsPort)); err != nil {
			log.Fatal().Err(err).Msg("Failed to start metrics server")
		}
	}(ctx)
}
