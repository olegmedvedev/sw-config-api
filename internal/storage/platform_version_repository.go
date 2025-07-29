package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// PlatformVersionRepositoryImpl implements service.PlatformVersionRepository
type PlatformVersionRepositoryImpl struct {
	db *sqlx.DB
}

// NewPlatformVersionRepository creates a new platform version repository
func NewPlatformVersionRepository(ctx context.Context, db *sqlx.DB) (*PlatformVersionRepositoryImpl, error) {
	return &PlatformVersionRepositoryImpl{
		db: db,
	}, nil
}

// GetPlatformVersion retrieves platform version information by platform
func (r *PlatformVersionRepositoryImpl) GetPlatformVersion(ctx context.Context, platform string) (*PlatformVersion, error) {
	var platformVersion PlatformVersion
	err := r.db.GetContext(ctx, &platformVersion,
		"SELECT required_version, store_version FROM platform_versions WHERE platform = ?", platform)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err // Return sql.ErrNoRows for "not found" case
		}
		return nil, err // Return original error for database issues
	}
	return &platformVersion, nil
}
