-- +goose Up

-- Insert platform versions data
INSERT INTO platform_versions (platform, required_version, store_version) VALUES
('android', '12.2.423', '13.7.556'),
('ios', '12.2.423', '13.7.556');

-- +goose Down
DELETE FROM platform_versions WHERE platform IN ('android', 'ios'); 