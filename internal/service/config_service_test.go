package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	serviceErrors "sw-config-api/internal/errors"
	"sw-config-api/internal/storage"
)

// Mock repositories
type MockResourceRepo struct {
	mock.Mock
}

func (m *MockResourceRepo) GetResource(ctx context.Context, platform, version string) (*storage.Resource, error) {
	args := m.Called(ctx, platform, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.Resource), args.Error(1)
}

func (m *MockResourceRepo) GetCompatibleResource(ctx context.Context, platform, appVersion string) (*storage.Resource, error) {
	args := m.Called(ctx, platform, appVersion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.Resource), args.Error(1)
}

type MockURLRepo struct {
	mock.Mock
}

func (m *MockURLRepo) ListURLs(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type MockPlatformVersionRepository struct {
	mock.Mock
}

func (m *MockPlatformVersionRepository) GetPlatformVersion(ctx context.Context, platform string) (*storage.PlatformVersion, error) {
	args := m.Called(ctx, platform)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.PlatformVersion), args.Error(1)
}

type MockEntryPointRepository struct {
	mock.Mock
}

func (m *MockEntryPointRepository) Get(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestConfigService_GetConfiguration_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:   "android",
		AppVersion: "13.6.956",
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets
	mockAssetRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(&storage.Resource{
		Version: "13.6.956",
		Hash:    "abc123",
	}, nil)

	// Mock definitions
	mockDefinitionRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(&storage.Resource{
		Version: "13.6.956",
		Hash:    "def456",
	}, nil)

	// Mock URLs
	mockAssetURLRepo.On("ListURLs", ctx).Return([]string{"https://cdn.example.com/assets"}, nil)
	mockDefinitionURLRepo.On("ListURLs", ctx).Return([]string{"https://cdn.example.com/definitions"}, nil)

	// Mock entry points
	mockEntryPointRepo.On("Get", ctx).Return(map[string]string{
		"backend_entry_point": "api.application.com/jsonrpc/v2",
		"notifications":       "notifications.application.com/jsonrpc/v1",
	}, nil)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "13.6.0", config.Version.Required)
	assert.Equal(t, "13.6.956", config.Version.Store)
	assert.Equal(t, "13.6.956", config.Assets.Version)
	assert.Equal(t, "abc123", config.Assets.Hash)
	assert.Equal(t, []string{"https://cdn.example.com/assets"}, config.Assets.Urls)
	assert.Equal(t, "13.6.956", config.Definitions.Version)
	assert.Equal(t, "def456", config.Definitions.Hash)
	assert.Equal(t, []string{"https://cdn.example.com/definitions"}, config.Definitions.Urls)
	assert.Equal(t, "api.application.com/jsonrpc/v2", config.BackendEntryPoint.JsonRpcUrl)
	assert.Equal(t, "notifications.application.com/jsonrpc/v1", config.Notifications.JsonRpcUrl)

	mockAssetRepo.AssertExpectations(t)
	mockDefinitionRepo.AssertExpectations(t)
	mockAssetURLRepo.AssertExpectations(t)
	mockDefinitionURLRepo.AssertExpectations(t)
	mockPlatformVersionRepo.AssertExpectations(t)
	mockEntryPointRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_WithExplicitVersions(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:           "android",
		AppVersion:         "13.6.956",
		AssetsVersion:      "13.6.955",
		DefinitionsVersion: "13.6.954",
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets with explicit version
	mockAssetRepo.On("GetResource", ctx, "android", "13.6.955").Return(&storage.Resource{
		Version: "13.6.955",
		Hash:    "abc123",
	}, nil)

	// Mock definitions with explicit version
	mockDefinitionRepo.On("GetResource", ctx, "android", "13.6.954").Return(&storage.Resource{
		Version: "13.6.954",
		Hash:    "def456",
	}, nil)

	// Mock URLs
	mockAssetURLRepo.On("ListURLs", ctx).Return([]string{"https://cdn.example.com/assets"}, nil)
	mockDefinitionURLRepo.On("ListURLs", ctx).Return([]string{"https://cdn.example.com/definitions"}, nil)

	// Mock entry points
	mockEntryPointRepo.On("Get", ctx).Return(map[string]string{
		"backend_entry_point": "api.application.com/jsonrpc/v2",
		"notifications":       "notifications.application.com/jsonrpc/v1",
	}, nil)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "13.6.955", config.Assets.Version)
	assert.Equal(t, "13.6.954", config.Definitions.Version)

	mockAssetRepo.AssertExpectations(t)
	mockDefinitionRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_PlatformNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:   "unknown",
		AppVersion: "13.6.956",
	}

	// Mock platform version not found
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "unknown").Return(nil, sql.ErrNoRows)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)

	var notFoundErr *serviceErrors.NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
	assert.Equal(t, "unknown", notFoundErr.Platform)

	mockPlatformVersionRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_AssetsNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:   "android",
		AppVersion: "13.6.956",
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets not found
	mockAssetRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(nil, sql.ErrNoRows)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)

	var notFoundErr *serviceErrors.NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
	assert.Equal(t, "android", notFoundErr.Platform)
	assert.Equal(t, "13.6.956", notFoundErr.AppVersion)

	mockPlatformVersionRepo.AssertExpectations(t)
	mockAssetRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_DefinitionsNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:   "android",
		AppVersion: "13.6.956",
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets found
	mockAssetRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(&storage.Resource{
		Version: "13.6.956",
		Hash:    "abc123",
	}, nil)

	// Mock definitions not found
	mockDefinitionRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(nil, sql.ErrNoRows)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)

	var notFoundErr *serviceErrors.NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
	assert.Equal(t, "android", notFoundErr.Platform)
	assert.Equal(t, "13.6.956", notFoundErr.AppVersion)

	mockPlatformVersionRepo.AssertExpectations(t)
	mockAssetRepo.AssertExpectations(t)
	mockDefinitionRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_IncompatibleAssetsVersion(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:      "android",
		AppVersion:    "13.6.956",
		AssetsVersion: "14.0.0", // Incompatible major version
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets with incompatible version
	mockAssetRepo.On("GetResource", ctx, "android", "14.0.0").Return(&storage.Resource{
		Version: "14.0.0",
		Hash:    "abc123",
	}, nil)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "specified assets version 14.0.0 is not compatible with app version 13.6.956")

	mockPlatformVersionRepo.AssertExpectations(t)
	mockAssetRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_IncompatibleDefinitionsVersion(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:           "android",
		AppVersion:         "13.6.956",
		DefinitionsVersion: "13.5.0", // Incompatible minor version
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets
	mockAssetRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(&storage.Resource{
		Version: "13.6.956",
		Hash:    "abc123",
	}, nil)

	// Mock definitions with incompatible version
	mockDefinitionRepo.On("GetResource", ctx, "android", "13.5.0").Return(&storage.Resource{
		Version: "13.5.0",
		Hash:    "def456",
	}, nil)

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "specified definitions version 13.5.0 is not compatible with app version 13.6.956")

	mockPlatformVersionRepo.AssertExpectations(t)
	mockAssetRepo.AssertExpectations(t)
	mockDefinitionRepo.AssertExpectations(t)
}

func TestConfigService_GetConfiguration_EntryPointsError(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockAssetRepo := &MockResourceRepo{}
	mockDefinitionRepo := &MockResourceRepo{}
	mockAssetURLRepo := &MockURLRepo{}
	mockDefinitionURLRepo := &MockURLRepo{}
	mockPlatformVersionRepo := &MockPlatformVersionRepository{}
	mockEntryPointRepo := &MockEntryPointRepository{}

	service := NewConfigService(
		mockAssetRepo,
		mockDefinitionRepo,
		mockAssetURLRepo,
		mockDefinitionURLRepo,
		mockPlatformVersionRepo,
		mockEntryPointRepo,
	)

	params := ClientParams{
		Platform:   "android",
		AppVersion: "13.6.956",
	}

	// Mock platform version
	mockPlatformVersionRepo.On("GetPlatformVersion", ctx, "android").Return(&storage.PlatformVersion{
		RequiredVersion: "13.6.0",
		StoreVersion:    "13.6.956",
	}, nil)

	// Mock assets
	mockAssetRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(&storage.Resource{
		Version: "13.6.956",
		Hash:    "abc123",
	}, nil)

	// Mock definitions
	mockDefinitionRepo.On("GetCompatibleResource", ctx, "android", "13.6.956").Return(&storage.Resource{
		Version: "13.6.956",
		Hash:    "def456",
	}, nil)

	// Mock URLs
	mockAssetURLRepo.On("ListURLs", ctx).Return([]string{"https://cdn.example.com/assets"}, nil)
	mockDefinitionURLRepo.On("ListURLs", ctx).Return([]string{"https://cdn.example.com/definitions"}, nil)

	// Mock entry points error
	mockEntryPointRepo.On("Get", ctx).Return(nil, errors.New("database error"))

	// Act
	config, err := service.GetConfiguration(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "failed to get entry points")

	mockPlatformVersionRepo.AssertExpectations(t)
	mockAssetRepo.AssertExpectations(t)
	mockDefinitionRepo.AssertExpectations(t)
	mockAssetURLRepo.AssertExpectations(t)
	mockDefinitionURLRepo.AssertExpectations(t)
	mockEntryPointRepo.AssertExpectations(t)
}

func TestConfigService_CompatibilityChecks(t *testing.T) {
	tests := []struct {
		name                  string
		appVersion            string
		assetsVersion         string
		definitionsVersion    string
		assetsCompatible      bool
		definitionsCompatible bool
	}{
		{
			name:             "Compatible assets - same major",
			appVersion:       "13.6.956",
			assetsVersion:    "13.5.100",
			assetsCompatible: true,
		},
		{
			name:             "Incompatible assets - different major",
			appVersion:       "13.6.956",
			assetsVersion:    "14.0.0",
			assetsCompatible: false,
		},
		{
			name:                  "Compatible definitions - same major and minor",
			appVersion:            "13.6.956",
			definitionsVersion:    "13.6.100",
			definitionsCompatible: true,
		},
		{
			name:                  "Incompatible definitions - different minor",
			appVersion:            "13.6.956",
			definitionsVersion:    "13.5.0",
			definitionsCompatible: false,
		},
		{
			name:                  "Incompatible definitions - different major",
			appVersion:            "13.6.956",
			definitionsVersion:    "14.0.0",
			definitionsCompatible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.assetsVersion != "" {
				result := isAssetsCompatible(tt.appVersion, tt.assetsVersion)
				assert.Equal(t, tt.assetsCompatible, result)
			}

			if tt.definitionsVersion != "" {
				result := isDefinitionsCompatible(tt.appVersion, tt.definitionsVersion)
				assert.Equal(t, tt.definitionsCompatible, result)
			}
		})
	}
}
