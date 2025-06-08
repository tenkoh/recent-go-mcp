package domain

import "time"

// GoRelease represents a Go version release with its updates
type GoRelease struct {
	Version     string                     `json:"version"`
	ReleaseDate time.Time                  `json:"release_date"`
	Summary     string                     `json:"summary"`
	Changes     []Change                   `json:"changes"`
	Packages    map[string][]PackageChange `json:"packages"`
}

// Change represents a general change in a Go release
type Change struct {
	Category    string `json:"category"` // "language", "runtime", "toolchain", "performance"
	Description string `json:"description"`
	Impact      string `json:"impact"` // "breaking", "enhancement", "deprecation", "new"
}

// PackageChange represents changes specific to a standard library package
type PackageChange struct {
	Function    string `json:"function,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Example     string `json:"example,omitempty"`
}

// FeatureResponse represents the response containing features available up to a version
type FeatureResponse struct {
	FromVersion string                     `json:"from_version"`
	ToVersion   string                     `json:"to_version"`
	Summary     string                     `json:"summary"`
	Changes     []Change                   `json:"changes"`
	PackageInfo map[string][]PackageChange `json:"package_info,omitempty"`
	// Version-specific data for formatted output
	VersionChanges  map[string][]Change                   `json:"-"`
	VersionPackages map[string]map[string][]PackageChange `json:"-"`
}
