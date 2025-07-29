package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ResourceRepositoryImpl implements service.ResourceRepo
type ResourceRepositoryImpl struct {
	db                        *sqlx.DB
	getResourceStmt           *sqlx.Stmt
	getCompatibleResourceStmt *sqlx.Stmt
	tableName                 string
	compatibility             VersionCompatibility
}

// NewResourceRepository creates a new resource repository
func NewResourceRepository(ctx context.Context, db *sqlx.DB, tableName string, compatibility VersionCompatibility) (*ResourceRepositoryImpl, error) {
	// Prepare statement for getting exact resource
	getResourceStmt, err := db.PreparexContext(ctx,
		fmt.Sprintf("SELECT platform, version, hash FROM %s WHERE platform = ? AND version = ?", tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getResource statement: %w", err)
	}

	// Prepare statement for compatible resources based on compatibility level
	var getCompatibleResourceStmt *sqlx.Stmt
	switch compatibility {
	case MajorOnly:
		getCompatibleResourceStmt, err = db.PreparexContext(ctx,
			fmt.Sprintf(`SELECT platform, version, hash FROM %s
			 WHERE platform = ?
			 AND SUBSTRING_INDEX(version, '.', 1) = SUBSTRING_INDEX(?, '.', 1)
			 ORDER BY version DESC
			 LIMIT 1`, tableName))
	case MajorMinor:
		getCompatibleResourceStmt, err = db.PreparexContext(ctx,
			fmt.Sprintf(`SELECT platform, version, hash FROM %s
			 WHERE platform = ?
			 AND CONCAT(SUBSTRING_INDEX(version, '.', 2)) = CONCAT(SUBSTRING_INDEX(?, '.', 2))
			 ORDER BY version DESC
			 LIMIT 1`, tableName))
	default:
		return nil, fmt.Errorf("unsupported compatibility level: %v", compatibility)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getCompatibleResource statement: %w", err)
	}

	return &ResourceRepositoryImpl{
		db:                        db,
		getResourceStmt:           getResourceStmt,
		getCompatibleResourceStmt: getCompatibleResourceStmt,
		tableName:                 tableName,
		compatibility:             compatibility,
	}, nil
}

// GetResource retrieves a resource by platform and version
func (r *ResourceRepositoryImpl) GetResource(ctx context.Context, platform, version string) (*Resource, error) {
	var resource Resource
	err := r.getResourceStmt.GetContext(ctx, &resource, platform, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err // Return sql.ErrNoRows for "not found" case
		}
		return nil, err // Return original error for database issues
	}
	return &resource, nil
}

// GetCompatibleResource retrieves a compatible resource by platform and app version
func (r *ResourceRepositoryImpl) GetCompatibleResource(ctx context.Context, platform, appVersion string) (*Resource, error) {
	var resource Resource
	err := r.getCompatibleResourceStmt.GetContext(ctx, &resource, platform, appVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err // Return sql.ErrNoRows for "not found" case
		}
		return nil, err // Return original error for database issues
	}
	return &resource, nil
}
