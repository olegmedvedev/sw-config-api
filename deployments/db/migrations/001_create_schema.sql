-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS assets (
    id INTEGER PRIMARY KEY,
    platform TEXT NOT NULL,
    version TEXT NOT NULL,
    hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(platform, version)
);

CREATE TABLE IF NOT EXISTS asset_urls (
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS definitions (
    id INTEGER PRIMARY KEY,
    platform TEXT NOT NULL,
    version TEXT NOT NULL,
    hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(platform, version)
);

CREATE TABLE IF NOT EXISTS definition_urls (
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_assets_platform_version ON assets(platform, version);
CREATE INDEX IF NOT EXISTS idx_definitions_platform_version ON definitions(platform, version);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_definitions_platform_version;
DROP INDEX IF EXISTS idx_assets_platform_version;
DROP TABLE IF EXISTS definition_urls;
DROP TABLE IF EXISTS definitions;
DROP TABLE IF EXISTS asset_urls;
DROP TABLE IF EXISTS assets;
-- +goose StatementEnd 