-- +goose Up

-- Import asset URLs from assets_urls.json
INSERT IGNORE INTO asset_urls (url)
SELECT JSON_UNQUOTE(JSON_EXTRACT(value, '$')) as url
FROM JSON_TABLE(
    JSON_EXTRACT(
        CONVERT(LOAD_FILE('/var/lib/mysql-files/data/assets_urls.json') USING utf8), 
        '$.assets_urls'
    ),
    '$[*]' COLUMNS (value JSON PATH '$')
) AS urls;

-- Import Android assets from assets.json
INSERT INTO assets (platform, version, hash)
SELECT 
    'android' as platform,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.version')) as version,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.hash')) as hash
FROM JSON_TABLE(
    JSON_EXTRACT(
        CONVERT(LOAD_FILE('/var/lib/mysql-files/data/assets.json') USING utf8), 
        '$.android'
    ),
    '$[*]' COLUMNS (value JSON PATH '$')
) AS android_assets
ON DUPLICATE KEY UPDATE
    hash = VALUES(hash);

-- Import iOS assets from assets.json
INSERT INTO assets (platform, version, hash)
SELECT 
    'ios' as platform,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.version')) as version,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.hash')) as hash
FROM JSON_TABLE(
    JSON_EXTRACT(
        CONVERT(LOAD_FILE('/var/lib/mysql-files/data/assets.json') USING utf8), 
        '$.ios'
    ),
    '$[*]' COLUMNS (value JSON PATH '$')
) AS ios_assets
ON DUPLICATE KEY UPDATE
    hash = VALUES(hash);

-- Import definition URLs from definitions_urls.json
INSERT IGNORE INTO definition_urls (url)
SELECT JSON_UNQUOTE(JSON_EXTRACT(value, '$')) as url
FROM JSON_TABLE(
    JSON_EXTRACT(
        CONVERT(LOAD_FILE('/var/lib/mysql-files/data/definitions_urls.json') USING utf8), 
        '$.definitions_urls'
    ),
    '$[*]' COLUMNS (value JSON PATH '$')
) AS urls;

-- Import Android definitions from definitions.json
INSERT INTO definitions (platform, version, hash)
SELECT 
    'android' as platform,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.version')) as version,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.hash')) as hash
FROM JSON_TABLE(
    JSON_EXTRACT(
        CONVERT(LOAD_FILE('/var/lib/mysql-files/data/definitions.json') USING utf8), 
        '$.android'
    ),
    '$[*]' COLUMNS (value JSON PATH '$')
) AS android_definitions
ON DUPLICATE KEY UPDATE
    hash = VALUES(hash);

-- Import iOS definitions from definitions.json
INSERT INTO definitions (platform, version, hash)
SELECT 
    'ios' as platform,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.version')) as version,
    JSON_UNQUOTE(JSON_EXTRACT(value, '$.hash')) as hash
FROM JSON_TABLE(
    JSON_EXTRACT(
        CONVERT(LOAD_FILE('/var/lib/mysql-files/data/definitions.json') USING utf8), 
        '$.ios'
    ),
    '$[*]' COLUMNS (value JSON PATH '$')
) AS ios_definitions
ON DUPLICATE KEY UPDATE
    hash = VALUES(hash);

-- +goose Down
DELETE FROM definition_urls;
DELETE FROM definitions;
DELETE FROM asset_urls;
DELETE FROM assets; 