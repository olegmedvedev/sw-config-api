package service

// Configuration represents the complete configuration for a client
type Configuration struct {
	Version     VersionInfo
	Assets      Resource
	Definitions Resource
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
