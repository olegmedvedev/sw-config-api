-- +goose Up

-- Add version component columns to assets table
ALTER TABLE assets 
ADD COLUMN major INT AFTER version,
ADD COLUMN minor INT AFTER major,
ADD COLUMN patch INT AFTER minor;

-- Add version component columns to definitions table
ALTER TABLE definitions 
ADD COLUMN major INT AFTER version,
ADD COLUMN minor INT AFTER major,
ADD COLUMN patch INT AFTER minor;

-- Update existing data with parsed version components
UPDATE assets SET 
    major = CAST(SUBSTRING_INDEX(version, '.', 1) AS UNSIGNED),
    minor = CAST(SUBSTRING_INDEX(SUBSTRING_INDEX(version, '.', 2), '.', -1) AS UNSIGNED),
    patch = CAST(SUBSTRING_INDEX(version, '.', -1) AS UNSIGNED);

UPDATE definitions SET 
    major = CAST(SUBSTRING_INDEX(version, '.', 1) AS UNSIGNED),
    minor = CAST(SUBSTRING_INDEX(SUBSTRING_INDEX(version, '.', 2), '.', -1) AS UNSIGNED),
    patch = CAST(SUBSTRING_INDEX(version, '.', -1) AS UNSIGNED);

-- Make version components NOT NULL after populating data
ALTER TABLE assets 
MODIFY COLUMN major INT NOT NULL,
MODIFY COLUMN minor INT NOT NULL,
MODIFY COLUMN patch INT NOT NULL;

ALTER TABLE definitions 
MODIFY COLUMN major INT NOT NULL,
MODIFY COLUMN minor INT NOT NULL,
MODIFY COLUMN patch INT NOT NULL;

-- Create indexes for efficient version component queries
CREATE INDEX idx_assets_platform_major_minor_patch ON assets(platform, major, minor, patch);

CREATE INDEX idx_definitions_platform_major_minor_patch ON definitions(platform, major, minor, patch);

-- +goose Down

-- Drop indexes
DROP INDEX IF EXISTS idx_definitions_platform_major_minor_patch ON definitions;

DROP INDEX IF EXISTS idx_assets_platform_major_minor_patch ON assets;

-- Drop version component columns
ALTER TABLE definitions 
DROP COLUMN patch,
DROP COLUMN minor,
DROP COLUMN major;

ALTER TABLE assets 
DROP COLUMN patch,
DROP COLUMN minor,
DROP COLUMN major; 