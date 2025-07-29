package service

import (
	"context"

	"sw-config-api/internal/storage"
)

// ResourceRepo interface for resource operations (assets, definitions, etc.)
type ResourceRepo interface {
	GetResource(ctx context.Context, platform, version string) (*storage.Resource, error)
	GetCompatibleResource(ctx context.Context, platform, appVersion string) (*storage.Resource, error)
}

// URLRepo interface for URL operations (asset URLs, definition URLs, etc.)
type URLRepo interface {
	ListURLs(ctx context.Context) ([]string, error)
}

// PlatformVersionRepository interface for platform version operations
type PlatformVersionRepository interface {
	GetPlatformVersion(ctx context.Context, platform string) (*storage.PlatformVersion, error)
}
