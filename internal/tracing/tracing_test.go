package tracing

import (
	"context"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestSetupTracing(t *testing.T) {
	// Set the required environment variables
	os.Setenv("UPTRACE_DSN", "test_dsn")
	os.Setenv("ENVIRONMENT", "test_environment")

	// Call the function
	tracer, shutdown := SetupTracing()

	// Check if a tracer was returned
	if tracer == nil || tracer != otel.Tracer("lotto-numbers") {
		t.Errorf("Expected a non-nil tracer, got %v", tracer)
	}

	// Check if a shutdown function was returned
	if shutdown == nil {
		t.Errorf("Expected a non-nil shutdown function")
	}

	// Call the shutdown function
	shutdown(context.Background())
}
