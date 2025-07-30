package service

// Configuration represents the complete configuration for a client
type Configuration struct {
	Version           VersionInfo
	BackendEntryPoint BackendService
	Assets            Resource
	Definitions       Resource
	Notifications     BackendService
}

// VersionInfo represents version information for a platform
type VersionInfo struct {
	Required string
	Store    string
}

// Resource represents a resource with version, hash and URLs
type Resource struct {
	Version string
	Hash    string
	Urls    []string
}

// BackendService represents a backend service configuration
type BackendService struct {
	JsonRpcUrl string `json:"jsonrpc_url"`
}
