-- +goose Up

-- Create entry_points table
CREATE TABLE IF NOT EXISTS entry_points (
    id INT AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(50) NOT NULL UNIQUE,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_entry_points_key ON entry_points(`key`);

-- +goose Down
DROP INDEX IF EXISTS idx_entry_points_key ON entry_points;
DROP TABLE IF EXISTS entry_points; 