package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// URLRepositoryImpl implements service.URLRepo
type URLRepositoryImpl struct {
	db           *sqlx.DB
	listURLsStmt *sqlx.Stmt
	tableName    string
}

// NewURLRepository creates a new URL repository
func NewURLRepository(ctx context.Context, db *sqlx.DB, tableName string) (*URLRepositoryImpl, error) {
	// Prepare statement for listing URLs
	listURLsStmt, err := db.PreparexContext(ctx,
		fmt.Sprintf("SELECT url FROM %s", tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare listURLs statement: %w", err)
	}

	return &URLRepositoryImpl{
		db:           db,
		listURLsStmt: listURLsStmt,
		tableName:    tableName,
	}, nil
}

// ListURLs retrieves all URLs
func (r *URLRepositoryImpl) ListURLs(ctx context.Context) ([]string, error) {
	var urls []struct {
		URL string `db:"url"`
	}
	err := r.listURLsStmt.SelectContext(ctx, &urls)
	if err != nil {
		return nil, err
	}

	// Convert to string slice
	result := make([]string, len(urls))
	for i, url := range urls {
		result[i] = url.URL
	}

	return result, nil
}
