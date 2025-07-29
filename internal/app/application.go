package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sw-config-api/internal/api"
	"sw-config-api/internal/middleware"
	"sw-config-api/internal/service"
	"sw-config-api/internal/storage"

	"github.com/jmoiron/sqlx"
)

type Application struct {
	logger     *slog.Logger
	db         *sqlx.DB
	handler    *service.Handler
	apiServer  *api.Server
	httpServer *http.Server
}

func New(ctx context.Context, config *Config) (*Application, error) {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize storage
	storageConfig := &storage.Config{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
	}

	db, err := storage.New(storageConfig)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	assetRepository, err := storage.NewResourceRepository(ctx, db, "assets", storage.MajorOnly)
	if err != nil {
		return nil, err
	}

	definitionRepository, err := storage.NewResourceRepository(ctx, db, "definitions", storage.MajorMinor)
	if err != nil {
		return nil, err
	}

	assetURLRepository, err := storage.NewURLRepository(ctx, db, "asset_urls")
	if err != nil {
		return nil, err
	}

	definitionURLRepository, err := storage.NewURLRepository(ctx, db, "definition_urls")
	if err != nil {
		return nil, err
	}

	platformVersionRepository, err := storage.NewPlatformVersionRepository(ctx, db)
	if err != nil {
		return nil, err
	}

	// Initialize config service
	configService := service.NewConfigService(
		assetRepository,
		definitionRepository,
		assetURLRepository,
		definitionURLRepository,
		platformVersionRepository,
	)

	// Initialize handler with config service
	handler := service.NewHandler(configService)

	// Create API server with custom error handler and logging middleware
	apiServer, err := api.NewServer(
		handler,
		api.WithErrorHandler(api.CustomErrorHandler),
		api.WithMiddleware(
			middleware.LoggingMiddleware(logger),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP server wrapper for graceful shutdown
	httpServer := &http.Server{
		Addr:         config.ServerAddr,
		Handler:      apiServer,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Application{
		logger:     logger,
		db:         db,
		handler:    handler,
		apiServer:  apiServer,
		httpServer: httpServer,
	}, nil
}

func (app *Application) Start() error {
	go func() {
		app.logger.Info("starting HTTP server", "addr", app.httpServer.Addr)
		if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	return nil
}

func (app *Application) Shutdown(ctx context.Context) error {
	app.logger.Info("shutting down server...")

	if err := app.httpServer.Shutdown(ctx); err != nil {
		app.logger.Error("server forced to shutdown", "error", err)
		return err
	}

	if err := app.db.Close(); err != nil {
		app.logger.Error("failed to close database", "error", err)
		return err
	}

	app.logger.Info("server exited")
	return nil
}

func (app *Application) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
