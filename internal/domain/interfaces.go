package domain

// ReleaseRepository handles access to Go release data
type ReleaseRepository interface {
	// GetAllReleases returns all available Go releases
	GetAllReleases() ([]GoRelease, error)
	
	// GetReleaseByVersion returns a specific release by version
	GetReleaseByVersion(version string) (*GoRelease, error)
	
	// GetReleasesUpToVersion returns all releases from oldest up to the specified version
	GetReleasesUpToVersion(targetVersion string) ([]GoRelease, error)
	
	// GetOldestVersion returns the oldest available version
	GetOldestVersion() (string, error)
	
	// GetLatestVersion returns the latest available version
	GetLatestVersion() (string, error)
}

// VersionComparator handles version comparison logic
type VersionComparator interface {
	// Compare compares two version strings
	// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
	Compare(v1, v2 string) int
}

// FeatureService provides business logic for feature retrieval
type FeatureService interface {
	// GetFeaturesForVersion returns all features available up to the specified version
	GetFeaturesForVersion(targetVersion string, packageName string) (*FeatureResponse, error)
}

// ResponseFormatter handles formatting of responses
type ResponseFormatter interface {
	// FormatAsText formats a FeatureResponse as human-readable text
	FormatAsText(response *FeatureResponse, version string, packageName string) string
}