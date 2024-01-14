package handlers

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

//go:embed web/index.html web/assets/*
var static embed.FS

// SetupRoutes configures the routes for the application.
func SetupRoutes(ctx context.Context, r *mux.Router) {
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "handlers.SetupRoutes")
	defer span.End()

	// Create a sub filesystem for the web directory.
	webFs, err := fs.Sub(static, "web")
	if err != nil {
		span.RecordError(err)
		log.Fatalf("Unable to create sub filesystem: %v", err)
	}
	// Create a sub filesystem for the assets directory.
	assetsFs, err := fs.Sub(static, "web/assets")
	if err != nil {
		span.RecordError(err)
		log.Fatalf("Unable to create sub filesystem for assets: %v", err)
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve the index.html file from the sub filesystem
		fileServer := http.FileServer(http.FS(webFs))
		fileServer.ServeHTTP(w, r)
	})
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFs))))
	r.HandleFunc("/numbers", GetLotteryNumbers).Methods("GET")
	r.HandleFunc("/version", VersionHandler).Methods("GET")
}
