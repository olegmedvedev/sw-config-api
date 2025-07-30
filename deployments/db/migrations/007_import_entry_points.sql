-- +goose Up

-- Insert initial entry points
INSERT INTO entry_points (`key`, url) VALUES 
('backend_entry_point', 'api.application.com/jsonrpc/v2'),
('notifications', 'notifications.application.com/jsonrpc/v1')
ON DUPLICATE KEY UPDATE url = VALUES(url);

-- +goose Down
DELETE FROM entry_points WHERE `key` IN ('backend_entry_point', 'notifications'); 