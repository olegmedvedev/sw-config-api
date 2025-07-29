-- +goose Up
-- +goose StatementBegin

-- Import asset URLs from assets_urls.json
INSERT OR IGNORE INTO asset_urls (url)
SELECT json_extract(value, '$') as url
FROM json_each(
    json_extract(
        readfile('deployments/db/data/assets_urls.json'), 
        '$.assets_urls'
    )
);

-- Import Android assets from assets.json
INSERT OR REPLACE INTO assets (platform, version, hash)
SELECT 
    'android' as platform,
    json_extract(value, '$.version') as version,
    json_extract(value, '$.hash') as hash
FROM json_each(
    json_extract(
        readfile('deployments/db/data/assets.json'), 
        '$.android'
    )
);

-- Import iOS assets from assets.json
INSERT OR REPLACE INTO assets (platform, version, hash)
SELECT 
    'ios' as platform,
    json_extract(value, '$.version') as version,
    json_extract(value, '$.hash') as hash
FROM json_each(
    json_extract(
        readfile('deployments/db/data/assets.json'), 
        '$.ios'
    )
);

-- Import definition URLs from definitions_urls.json
INSERT OR IGNORE INTO definition_urls (url)
SELECT json_extract(value, '$') as url
FROM json_each(
    json_extract(
        readfile('deployments/db/data/definitions_urls.json'), 
        '$.definitions_urls'
    )
);

-- Import Android definitions from definitions.json
INSERT OR REPLACE INTO definitions (platform, version, hash)
SELECT 
    'android' as platform,
    json_extract(value, '$.version') as version,
    json_extract(value, '$.hash') as hash
FROM json_each(
    json_extract(
        readfile('deployments/db/data/definitions.json'), 
        '$.android'
    )
);

-- Import iOS definitions from definitions.json
INSERT OR REPLACE INTO definitions (platform, version, hash)
SELECT 
    'ios' as platform,
    json_extract(value, '$.version') as version,
    json_extract(value, '$.hash') as hash
FROM json_each(
    json_extract(
        readfile('deployments/db/data/definitions.json'), 
        '$.ios'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM definition_urls;
DELETE FROM definitions;
DELETE FROM asset_urls;
DELETE FROM assets;
-- +goose StatementEnd 