package service

import "context"

// ConfigServiceInterface defines the interface for config service operations
type ConfigServiceInterface interface {
	GetConfiguration(ctx context.Context, params ClientParams) (*Configuration, error)
}
