# E2E Tests

End-to-End tests for full API integration.

## Test Cases

### 1. Successful Configuration Retrieval

Tests successful configuration retrieval for different platforms and versions.

| Test Case | Platform | App Version | Description | Expected |
|-----------|----------|-------------|-------------|----------|
| `android_latest` | android | `14.8.447` | Android with latest version | 200 OK |
| `ios_latest` | ios | `14.6.743` | iOS with latest version | 200 OK |
| `android_older_version` | android | `13.2.528` | Android with older version | 200 OK |
| `ios_older_version` | ios | `12.4.328` | iOS with older version | 200 OK |

**Response validation:**
- Contains `version` field
- Contains `backend_entry_point` field  
- Contains `assets` field
- Contains `definitions` field
- Contains `notifications` field

### 2. Configuration with Explicit Versions

Tests configuration retrieval with explicit assets and definitions versions.

| Test Case | App Version | Platform | Assets Version | Definitions Version | Expected |
|-----------|-------------|----------|----------------|-------------------|----------|
| `explicit_versions_valid` | `14.8.447` | android | `14.8.447` | `14.8.98` | 200 OK |
| `explicit_assets_only` | `14.8.447` | android | `14.8.447` | - | 200 OK |
| `explicit_definitions_only` | `14.8.447` | android | - | `14.8.98` | 200 OK |

### 3. Parameter Validation

Tests validation of required parameters and SemVer format.

| Test Case | Query Parameters | Expected | Description |
|-----------|------------------|----------|-------------|
| `missing_app_version` | `platform=android` | 400 Bad Request | Missing appVersion |
| `missing_platform` | `appVersion=14.8.447` | 400 Bad Request | Missing platform |
| `invalid_semver_format` | `appVersion=invalid&platform=android` | 400 Bad Request | Invalid SemVer format |
| `invalid_semver_with_v` | `appVersion=v14.8.447&platform=android` | 400 Bad Request | SemVer with v prefix |

**Error response validation:**
- Contains `error` object
- Contains `code` field (400)
- Contains `message` field

### 4. Not Found Scenarios

Tests scenarios where configuration is not found.

| Test Case | App Version | Platform | Expected | Description |
|-----------|-------------|----------|----------|-------------|
| `non_existent_platform` | `14.8.447` | `non_existent_platform` | 404 Not Found | Non-existent platform |
| `very_old_version` | `10.0.0` | android | 404 Not Found | Very old version |
| `future_version` | `20.0.0` | android | 404 Not Found | Future version |

**Error response validation:**
- Contains `error` object
- Contains `code` field (404)
- Contains `message` field with "Configuration not found"

### 5. SemVer Compatibility

Tests version compatibility rules for assets (MAJOR only) and definitions (MAJOR.MINOR).

| Test Case | App Version | Platform | Assets Version | Definitions Version | Expected | Description |
|-----------|-------------|----------|----------------|-------------------|----------|-------------|
| `major_compatibility_assets` | `14.8.447` | android | `14.0.357` | - | 200 OK | MAJOR compatibility for assets |
| `major_minor_compatibility_definitions` | `14.8.447` | android | - | `14.8.98` | 200 OK | MAJOR.MINOR compatibility for definitions |
| `incompatible_major_assets` | `14.8.447` | android | `13.2.528` | - | 404 Not Found | Incompatible MAJOR version for assets |
| `incompatible_minor_definitions` | `14.8.447` | android | - | `14.1.822` | 404 Not Found | Incompatible MINOR version for definitions |

## Data Sources

Test data is based on real data from database migrations:

### Assets Versions
- **Android**: `13.2.528`, `13.9.519`, `14.7.159`, `14.8.40`, `14.2.723`, `14.8.447`, `13.5.244`, `14.7.919`, `14.0.357`, `13.5.275`
- **iOS**: `14.4.861`, `13.1.783`, `14.6.743`, `12.4.328`, `13.2.906`, `14.0.415`, `12.0.631`, `13.5.397`, `12.3.817`, `12.6.496`

### Definitions Versions  
- **Android**: `14.8.98`, `12.3.567`, `12.8.199`, `12.0.177`, `13.6.610`, `13.2.20`, `12.4.155`, `14.1.822`, `12.6.962`, `13.2.296`
- **iOS**: `12.3.807`, `14.2.281`, `13.3.481`, `13.5.693`, `14.5.580`, `12.4.454`, `13.0.435`, `12.4.623`, `12.8.819`, `12.1.761`

### Platform Versions
- **Required version**: `12.2.423`
- **Store version**: `13.7.556`

## Compatibility Rules

- **Assets**: MAJOR version compatibility (e.g., `14.x.x` compatible with `14.y.z`)
- **Definitions**: MAJOR.MINOR version compatibility (e.g., `14.8.x` compatible with `14.8.y`)

## Running Tests

```bash
# Start infrastructure
docker-compose up -d

# Run tests
cd tests
go test -v

# With custom URL
E2E_BASE_URL=http://localhost:8080 go test -v
``` 