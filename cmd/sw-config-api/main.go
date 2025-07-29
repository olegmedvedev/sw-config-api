package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"sw-config-api/internal/app"
)

func main() {
	// Load configuration from environment variables
	config, err := app.LoadConfig(context.Background())
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create and initialize application
	application, err := app.New(context.Background(), config)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start the application
	if err := application.Start(); err != nil {
		slog.Error("failed to start application", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown signal
	application.WaitForShutdown()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown application", "error", err)
		os.Exit(1)
	}
}
