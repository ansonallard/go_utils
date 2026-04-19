package logging

import (
	"context"
	"io"
	"os"
	"time"

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
