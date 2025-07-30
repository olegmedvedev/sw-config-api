package service

import (
	"context"
	"log/slog"

	"sw-config-api/internal/api"
)

// Handler handles API requests and business logic
type Handler struct {
	configService ConfigServiceInterface
	logger        *slog.Logger
}

// NewHandler creates a new handler with config service
func NewHandler(configService ConfigServiceInterface, logger *slog.Logger) *Handler {
	return &Handler{
		configService: configService,
		logger:        logger,
	}
}

// ConfigGet implements GET /config operation.
//
// Get configuration for client.
//
// GET /config
func (h *Handler) ConfigGet(ctx context.Context, params api.ConfigGetParams) (api.ConfigGetRes, error) {
	// Determine which versions to use for assets and definitions
	assetsVersion := string(params.AppVersion)
	if params.AssetsVersion.IsSet() {
		if assetsVer, ok := params.AssetsVersion.Get(); ok {
			assetsVersion = string(assetsVer)
		}
	}

	definitionsVersion := string(params.AppVersion)
	if params.DefinitionsVersion.IsSet() {
		if defVer, ok := params.DefinitionsVersion.Get(); ok {
			definitionsVersion = string(defVer)
		}
	}

	// Create config parameters
	clientParams := ClientParams{
		Platform:           params.Platform,
		AppVersion:         string(params.AppVersion),
		AssetsVersion:      assetsVersion,
		DefinitionsVersion: definitionsVersion,
	}

	// Get configuration from business logic layer
	config, err := h.configService.GetConfiguration(ctx, clientParams)
	if err != nil {
		// Check if it's a "not found" error and return appropriate API response
		if IsNotFoundError(err) {
			// Log the not found error
			h.logger.Warn("Configuration not found",
				"error", err.Error(),
				"platform", clientParams.Platform,
				"appVersion", clientParams.AppVersion,
			)

			return &api.ConfigGetNotFound{
				Error: api.NewOptConfigGetNotFoundError(api.ConfigGetNotFoundError{
					Code:    api.NewOptInt(404),
					Message: api.NewOptString("Configuration not found"),
				}),
			}, nil
		}
		// Return internal server error for other errors
		return nil, err
	}

	// Convert business model to API model
	apiConfig := &api.Config{
		Version: api.NewOptVersion(api.Version{
			Required: api.NewOptSemVer(api.SemVer(config.Version.Required)),
			Store:    api.NewOptSemVer(api.SemVer(config.Version.Store)),
		}),
		Assets: api.NewOptResource(api.Resource{
			Version: api.NewOptSemVer(api.SemVer(config.Assets.Version)),
			Hash:    api.NewOptString(config.Assets.Hash),
			Urls:    config.Assets.Urls,
		}),
		Definitions: api.NewOptResource(api.Resource{
			Version: api.NewOptSemVer(api.SemVer(config.Definitions.Version)),
			Hash:    api.NewOptString(config.Definitions.Hash),
			Urls:    config.Definitions.Urls,
		}),
	}

	return apiConfig, nil
}
