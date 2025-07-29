package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sw-config-api/internal/storage"

	"github.com/Masterminds/semver"
)

// ClientParams represents parameters for configuration retrieval from client
type ClientParams struct {
	Platform           string
	AppVersion         string
	AssetsVersion      string
	DefinitionsVersion string
}

// ConfigService handles business logic for configuration operations
type ConfigService struct {
	assetRepository           ResourceRepo
	definitionRepository      ResourceRepo
	assetURLRepository        URLRepo
	definitionURLRepository   URLRepo
	platformVersionRepository PlatformVersionRepository
}

// NewConfigService creates a new config service
func NewConfigService(
	assetRepository ResourceRepo,
	definitionRepository ResourceRepo,
	assetURLRepository URLRepo,
	definitionURLRepository URLRepo,
	platformVersionRepository PlatformVersionRepository,
) *ConfigService {
	return &ConfigService{
		assetRepository:           assetRepository,
		definitionRepository:      definitionRepository,
		assetURLRepository:        assetURLRepository,
		definitionURLRepository:   definitionURLRepository,
		platformVersionRepository: platformVersionRepository,
	}
}

// GetConfiguration retrieves configuration for the given parameters
func (s *ConfigService) GetConfiguration(ctx context.Context, params ClientParams) (*Configuration, error) {
	// Get platform version information
	platformVersion, err := s.platformVersionRepository.GetPlatformVersion(ctx, params.Platform)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{
				Platform: params.Platform,
			}
		}
		return nil, err // Return original error for database issues
	}

	var asset *storage.Resource
	var definition *storage.Resource

	// Handle assets version selection
	if params.AssetsVersion != "" {
		// Client explicitly specified assetsVersion - try to get exact version
		asset, err = s.assetRepository.GetResource(ctx, params.Platform, params.AssetsVersion)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("specified assets version %s not found: %w", params.AssetsVersion, &NotFoundError{
					Platform:   params.Platform,
					AppVersion: params.AppVersion,
				})
			}
			return nil, err // Return original error for database issues
		}
		// Validate that specified assets version is compatible with app version
		if !isAssetsCompatible(params.AppVersion, asset.Version) {
			return nil, fmt.Errorf("specified assets version %s is not compatible with app version %s: %w", asset.Version, params.AppVersion, &NotFoundError{
				Platform:   params.Platform,
				AppVersion: params.AppVersion,
			})
		}
	} else {
		// No explicit assetsVersion - find compatible version
		asset, err = s.assetRepository.GetCompatibleResource(ctx, params.Platform, params.AppVersion)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("no compatible assets version found: %w", &NotFoundError{
					Platform:   params.Platform,
					AppVersion: params.AppVersion,
				})
			}
			return nil, err // Return original error for database issues
		}
	}

	// Handle definitions version selection
	if params.DefinitionsVersion != "" {
		// Client explicitly specified definitionsVersion - try to get exact version
		definition, err = s.definitionRepository.GetResource(ctx, params.Platform, params.DefinitionsVersion)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("specified definitions version %s not found: %w", params.DefinitionsVersion, &NotFoundError{
					Platform:   params.Platform,
					AppVersion: params.AppVersion,
				})
			}
			return nil, err // Return original error for database issues
		}
		// Validate that specified definitions version is compatible with app version
		if !isDefinitionsCompatible(params.AppVersion, definition.Version) {
			return nil, fmt.Errorf("specified definitions version %s is not compatible with app version %s: %w", definition.Version, params.AppVersion, &NotFoundError{
				Platform:   params.Platform,
				AppVersion: params.AppVersion,
			})
		}
	} else {
		// No explicit definitionsVersion - find compatible version
		definition, err = s.definitionRepository.GetCompatibleResource(ctx, params.Platform, params.AppVersion)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("no compatible definitions version found: %w", &NotFoundError{
					Platform:   params.Platform,
					AppVersion: params.AppVersion,
				})
			}
			return nil, err // Return original error for database issues
		}
	}

	// Get asset URLs
	assetURLs, err := s.assetURLRepository.ListURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset URLs: %w", err)
	}

	// Get definition URLs
	definitionURLs, err := s.definitionURLRepository.ListURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get definition URLs: %w", err)
	}

	// Build configuration
	config := &Configuration{
		Version: VersionInfo{
			Required: platformVersion.RequiredVersion,
			Store:    platformVersion.StoreVersion,
		},
		Assets: Resource{
			Version: asset.Version,
			Hash:    asset.Hash,
			Urls:    assetURLs,
		},
		Definitions: Resource{
			Version: definition.Version,
			Hash:    definition.Hash,
			Urls:    definitionURLs,
		},
	}

	return config, nil
}

// isAssetsCompatible checks if assets version is compatible with app version
// Assets are compatible if MAJOR version matches (MAJOR.MINOR.PATCH)
func isAssetsCompatible(appVersion, assetsVersion string) bool {
	appVer, err := semver.NewVersion(appVersion)
	if err != nil {
		return false
	}

	assetsVer, err := semver.NewVersion(assetsVersion)
	if err != nil {
		return false
	}

	return appVer.Major() == assetsVer.Major()
}

// isDefinitionsCompatible checks if definitions version is compatible with app version
// Definitions are compatible if MAJOR and MINOR versions match (MAJOR.MINOR.PATCH)
func isDefinitionsCompatible(appVersion, definitionsVersion string) bool {
	appVer, err := semver.NewVersion(appVersion)
	if err != nil {
		return false
	}

	defVer, err := semver.NewVersion(definitionsVersion)
	if err != nil {
		return false
	}

	return appVer.Major() == defVer.Major() && appVer.Minor() == defVer.Minor()
}
