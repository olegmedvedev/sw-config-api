package storage

// VersionCompatibility defines how version compatibility should be checked
type VersionCompatibility int

const (
	MajorOnly  VersionCompatibility = iota // Only MAJOR version must match
	MajorMinor                             // MAJOR and MINOR versions must match
)

// Resource represents a generic resource in the database (asset, definition, etc.)
type Resource struct {
	Version string `db:"version"`
	Hash    string `db:"hash"`
}

// PlatformVersion represents platform version information in the database
type PlatformVersion struct {
	RequiredVersion string `db:"required_version"`
	StoreVersion    string `db:"store_version"`
}
