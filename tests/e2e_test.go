package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultBaseURL = "http://localhost:8080"
	defaultTimeout = 30 * time.Second
)

// getBaseURL returns URL for tests
func getBaseURL() string {
	if url := os.Getenv("E2E_BASE_URL"); url != "" {
		return url
	}
	return defaultBaseURL
}

// getHTTPClient creates HTTP client with timeout
func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimeout,
	}
}

// TestSuccessfulConfigRetrieval tests successful configuration retrieval
func TestSuccessfulConfigRetrieval(t *testing.T) {
	client := getHTTPClient()
	baseURL := getBaseURL()

	testCases := []struct {
		name        string
		appVersion  string
		platform    string
		description string
	}{
		{
			name:        "android_latest",
			appVersion:  "14.8.447",
			platform:    "android",
			description: "Android with latest version",
		},
		{
			name:        "ios_latest",
			appVersion:  "14.5.580",
			platform:    "ios",
			description: "iOS with latest version",
		},
		{
			name:        "android_older_version",
			appVersion:  "13.2.528",
			platform:    "android",
			description: "Android with older version",
		},
		{
			name:        "ios_older_version",
			appVersion:  "12.4.328",
			platform:    "ios",
			description: "iOS with older version",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/config?appVersion=%s&platform=%s", baseURL, tc.appVersion, tc.platform)

			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK for %s", tc.description)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var configResponse map[string]interface{}
			err = json.Unmarshal(body, &configResponse)
			require.NoError(t, err)

			// Check response structure
			assert.Contains(t, configResponse, "version")
			assert.Contains(t, configResponse, "backend_entry_point")
			assert.Contains(t, configResponse, "assets")
			assert.Contains(t, configResponse, "definitions")
			assert.Contains(t, configResponse, "notifications")
		})
	}
}

// TestConfigWithExplicitVersions tests configuration retrieval with explicit versions
func TestConfigWithExplicitVersions(t *testing.T) {
	client := getHTTPClient()
	baseURL := getBaseURL()

	testCases := []struct {
		name               string
		appVersion         string
		platform           string
		assetsVersion      string
		definitionsVersion string
		expectedStatus     int
		description        string
	}{
		{
			name:               "explicit_versions_valid",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "14.8.447",
			definitionsVersion: "14.8.98",
			expectedStatus:     http.StatusOK,
			description:        "Explicit versions - valid",
		},
		{
			name:               "explicit_assets_only",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "14.8.447",
			definitionsVersion: "",
			expectedStatus:     http.StatusOK,
			description:        "Only explicit assets version",
		},
		{
			name:               "explicit_definitions_only",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "",
			definitionsVersion: "14.8.98",
			expectedStatus:     http.StatusOK,
			description:        "Only explicit definitions version",
		},
		{
			name:               "assets_version_exact_match",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "14.8.447",
			definitionsVersion: "",
			expectedStatus:     http.StatusOK,
			description:        "Assets version exact match",
		},
		{
			name:               "assets_version_incompatible_major",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "13.2.528",
			definitionsVersion: "",
			expectedStatus:     http.StatusNotFound,
			description:        "Assets version incompatible major",
		},
		{
			name:               "assets_and_definitions_versions",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "14.8.447",
			definitionsVersion: "14.8.98",
			expectedStatus:     http.StatusOK,
			description:        "Both assets and definitions versions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/config?appVersion=%s&platform=%s", baseURL, tc.appVersion, tc.platform)

			if tc.assetsVersion != "" {
				url += fmt.Sprintf("&assetsVersion=%s", tc.assetsVersion)
			}
			if tc.definitionsVersion != "" {
				url += fmt.Sprintf("&definitionsVersion=%s", tc.definitionsVersion)
			}

			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected %d for %s", tc.expectedStatus, tc.description)

			if tc.expectedStatus == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var configResponse map[string]interface{}
				err = json.Unmarshal(body, &configResponse)
				require.NoError(t, err)

				// Check that response contains all required fields
				assert.Contains(t, configResponse, "version")
				assert.Contains(t, configResponse, "backend_entry_point")
				assert.Contains(t, configResponse, "assets")
				assert.Contains(t, configResponse, "definitions")
				assert.Contains(t, configResponse, "notifications")
			}
		})
	}
}

// TestValidationErrors tests parameter validation
func TestValidationErrors(t *testing.T) {
	client := getHTTPClient()
	baseURL := getBaseURL()

	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		description    string
	}{
		{
			name:           "missing_app_version",
			queryParams:    "platform=android",
			expectedStatus: http.StatusBadRequest,
			description:    "Missing appVersion",
		},
		{
			name:           "missing_platform",
			queryParams:    "appVersion=14.8.447",
			expectedStatus: http.StatusBadRequest,
			description:    "Missing platform",
		},
		{
			name:           "invalid_semver_format",
			queryParams:    "appVersion=invalid&platform=android",
			expectedStatus: http.StatusBadRequest,
			description:    "Invalid SemVer format",
		},
		{
			name:           "invalid_semver_with_v",
			queryParams:    "appVersion=v14.8.447&platform=android",
			expectedStatus: http.StatusBadRequest,
			description:    "SemVer with v prefix",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := baseURL + "/config"
			if tc.queryParams != "" {
				url += "?" + tc.queryParams
			}

			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected %d for %s", tc.expectedStatus, tc.description)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var errorResponse map[string]interface{}
			err = json.Unmarshal(body, &errorResponse)
			require.NoError(t, err)

			// Check error structure
			assert.Contains(t, errorResponse, "error")
			errorObj := errorResponse["error"].(map[string]interface{})
			assert.Contains(t, errorObj, "code")
			assert.Contains(t, errorObj, "message")
			assert.Equal(t, float64(tc.expectedStatus), errorObj["code"])
		})
	}
}

// TestNotFoundScenarios tests "not found" scenarios
func TestNotFoundScenarios(t *testing.T) {
	client := getHTTPClient()
	baseURL := getBaseURL()

	testCases := []struct {
		name           string
		appVersion     string
		platform       string
		expectedStatus int
		description    string
	}{
		{
			name:           "non_existent_platform",
			appVersion:     "14.8.447",
			platform:       "non_existent_platform",
			expectedStatus: http.StatusNotFound,
			description:    "Non-existent platform",
		},
		{
			name:           "very_old_version",
			appVersion:     "10.0.0",
			platform:       "android",
			expectedStatus: http.StatusNotFound,
			description:    "Very old version",
		},
		{
			name:           "future_version",
			appVersion:     "20.0.0",
			platform:       "android",
			expectedStatus: http.StatusNotFound,
			description:    "Future version",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/config?appVersion=%s&platform=%s", baseURL, tc.appVersion, tc.platform)

			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected %d for %s", tc.expectedStatus, tc.description)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var errorResponse map[string]interface{}
			err = json.Unmarshal(body, &errorResponse)
			require.NoError(t, err)

			// Check error structure
			assert.Contains(t, errorResponse, "error")
			errorObj := errorResponse["error"].(map[string]interface{})
			assert.Contains(t, errorObj, "code")
			assert.Contains(t, errorObj, "message")
			assert.Equal(t, float64(tc.expectedStatus), errorObj["code"])
			assert.Equal(t, "Configuration not found", errorObj["message"])
		})
	}
}

// TestSemVerCompatibility tests version compatibility
func TestSemVerCompatibility(t *testing.T) {
	client := getHTTPClient()
	baseURL := getBaseURL()

	testCases := []struct {
		name               string
		appVersion         string
		platform           string
		assetsVersion      string
		definitionsVersion string
		expectedStatus     int
		description        string
	}{
		{
			name:               "major_compatibility_assets",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "14.0.357", // Should find compatible version 14.x.x
			definitionsVersion: "",
			expectedStatus:     http.StatusOK,
			description:        "MAJOR compatibility for assets",
		},
		{
			name:               "major_minor_compatibility_definitions",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "",
			definitionsVersion: "14.8.98", // Should find exact version 14.8.x
			expectedStatus:     http.StatusOK,
			description:        "MAJOR.MINOR compatibility for definitions",
		},
		{
			name:               "incompatible_major_assets",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "13.2.528", // Incompatible MAJOR version
			definitionsVersion: "",
			expectedStatus:     http.StatusNotFound,
			description:        "Incompatible MAJOR version for assets",
		},
		{
			name:               "incompatible_minor_definitions",
			appVersion:         "14.8.447",
			platform:           "android",
			assetsVersion:      "",
			definitionsVersion: "14.1.822", // Incompatible MINOR version
			expectedStatus:     http.StatusNotFound,
			description:        "Incompatible MINOR version for definitions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/config?appVersion=%s&platform=%s", baseURL, tc.appVersion, tc.platform)

			if tc.assetsVersion != "" {
				url += fmt.Sprintf("&assetsVersion=%s", tc.assetsVersion)
			}
			if tc.definitionsVersion != "" {
				url += fmt.Sprintf("&definitionsVersion=%s", tc.definitionsVersion)
			}

			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected %d for %s", tc.expectedStatus, tc.description)
		})
	}
}
