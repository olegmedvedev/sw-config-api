package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// EntryPointRepository handles database operations for entry points
type EntryPointRepository struct {
	db    *sqlx.DB
	query *sqlx.Stmt
}

// NewEntryPointRepository creates a new entry point repository
func NewEntryPointRepository(db *sqlx.DB) (*EntryPointRepository, error) {
	query := "SELECT `key`, url FROM entry_points"
	stmt, err := db.PreparexContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare entry points query: %w", err)
	}

	return &EntryPointRepository{
		db:    db,
		query: stmt,
	}, nil
}

// Get retrieves all entry points as a map of key to URL
func (r *EntryPointRepository) Get(ctx context.Context) (map[string]string, error) {
	rows, err := r.query.QueryxContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query entry points: %w", err)
	}
	defer rows.Close()

	entryPoints := make(map[string]string)
	for rows.Next() {
		var key, url string
		if err := rows.Scan(&key, &url); err != nil {
			return nil, fmt.Errorf("failed to scan entry point: %w", err)
		}
		entryPoints[key] = url
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over entry points: %w", err)
	}

	return entryPoints, nil
}
