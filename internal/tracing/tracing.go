// Package tracing configures and initializes OpenTelemetry tracing with Uptrace.
package tracing

import (
	"context"
	"log"
	"os"

	"github.com/danstis/lotto-numbers/internal/version"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func SetupTracing() (trace.Tracer, func(ctx context.Context)) {
	// Retrieve the uptrace environment variables.
	dsn := os.Getenv("UPTRACE_DSN")
	if dsn == "" {
		log.Panicf("UPTRACE_DSN environment variable not set")
	}
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	// Create a new tracer
	// Note: the service name must match the name of the service in Uptrace
	// Note: the DSN must be set in the UPTRACE_DSN environment variable
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(dsn),
		uptrace.WithServiceName("lotto-numbers"),
		uptrace.WithServiceVersion(version.Version),
		uptrace.WithDeploymentEnvironment(environment),
	)

	tracer := otel.Tracer("lotto-numbers")

	return tracer, func(ctx context.Context) {
		err := uptrace.Shutdown(ctx)
		if err != nil {
			// handle error
			log.Println(err)
		}
	}
}
