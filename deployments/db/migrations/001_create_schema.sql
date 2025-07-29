-- +goose Up

-- Create assets table
CREATE TABLE IF NOT EXISTS assets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_platform_version (platform, version)
);

-- Create asset_urls table
CREATE TABLE IF NOT EXISTS asset_urls (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create definitions table
CREATE TABLE IF NOT EXISTS definitions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_platform_version (platform, version)
);

-- Create definition_urls table
CREATE TABLE IF NOT EXISTS definition_urls (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_assets_platform_version ON assets(platform, version);
CREATE INDEX idx_definitions_platform_version ON definitions(platform, version);

-- +goose Down
DROP INDEX IF EXISTS idx_definitions_platform_version ON definitions;
DROP INDEX IF EXISTS idx_assets_platform_version ON assets;
DROP TABLE IF EXISTS definition_urls;
DROP TABLE IF EXISTS definitions;
DROP TABLE IF EXISTS asset_urls;
DROP TABLE IF EXISTS assets; 