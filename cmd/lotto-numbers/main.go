package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danstis/lotto-numbers/internal/handlers"
	"github.com/danstis/lotto-numbers/internal/middleware"
	"github.com/danstis/lotto-numbers/internal/tracing"
	"github.com/danstis/lotto-numbers/internal/version"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
)

// Main entry point for the app
func main() {
	log.Printf("Lotto Numbers started - v%v", version.Version)

	// Start Tracing.
	ctx := context.Background()
	tracer, shutdown := tracing.SetupTracing()
	defer shutdown(ctx)

	ctx, main := tracer.Start(ctx, "startup")
	defer main.End()

	// Set up the HTTP server using Gorilla Mux
	r := mux.NewRouter()

	// Wrap the router with the logging and tracing middlewares
	r.Use(otelmux.Middleware("lotto-numbers"))
	r.Use(middleware.LoggingMiddleware)

	handlers.SetupRoutes(ctx, r)

	// Retrieve the port number from an environment variable or use 8080 as default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8091"
	}
	main.SetAttributes(attribute.String("port", port))

	main.End()

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
