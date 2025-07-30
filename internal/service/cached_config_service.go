package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"sw-config-api/internal/cache"
)

// CachedConfigService wraps ConfigService with caching
type CachedConfigService struct {
	configService *ConfigService
	cache         cache.Interface
	ttl           time.Duration
	logger        *slog.Logger
}

// NewCachedConfigService creates a new cached config service
func NewCachedConfigService(configService *ConfigService, cache cache.Interface, ttl time.Duration, logger *slog.Logger) *CachedConfigService {
	return &CachedConfigService{
		configService: configService,
		cache:         cache,
		ttl:           ttl,
		logger:        logger,
	}
}

// GetConfiguration retrieves configuration with caching
func (s *CachedConfigService) GetConfiguration(ctx context.Context, params ClientParams) (*Configuration, error) {
	// Generate cache key based on parameters
	cacheKey := s.generateCacheKey(params)

	// Try to get from cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		var config Configuration
		if err := json.Unmarshal(cached, &config); err == nil {
			return &config, nil
		}
		// If unmarshal fails, continue to get fresh data
	}

	// If not in cache, get from underlying service
	config, err := s.configService.GetConfiguration(ctx, params)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(config); err == nil {
		if err := s.cache.Set(cacheKey, data, s.ttl); err != nil {
			// Log cache error but don't fail the request
			s.logger.Error("failed to cache configuration",
				"error", err.Error(),
				"cache_key", cacheKey,
			)
		}
	}

	return config, nil
}

// generateCacheKey creates a unique cache key based on request parameters
// Format: config:{platform}:{appVersion}:{assetsVersion}:{definitionsVersion}
func (s *CachedConfigService) generateCacheKey(params ClientParams) string {
	var builder strings.Builder

	// Build key with required parameters
	builder.WriteString("config:")
	builder.WriteString(params.Platform)
	builder.WriteString(":")
	builder.WriteString(params.AppVersion)

	// Add optional assets version
	builder.WriteString(":")
	if params.AssetsVersion != "" {
		builder.WriteString(params.AssetsVersion)
	}

	// Add optional definitions version
	builder.WriteString(":")
	if params.DefinitionsVersion != "" {
		builder.WriteString(params.DefinitionsVersion)
	}

	return builder.String()
}
