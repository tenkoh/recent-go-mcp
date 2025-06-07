package domain

import "context"

// ReleaseRepository handles access to Go release data
type ReleaseRepository interface {
	// GetAllReleases returns all available Go releases
	// Note: Returned pointers should be treated as read-only to maintain data integrity
	GetAllReleases(ctx context.Context) ([]*GoRelease, error)
	
	// GetReleaseByVersion returns a specific release by version
	// Note: Returned pointer should be treated as read-only to maintain data integrity
	GetReleaseByVersion(ctx context.Context, version string) (*GoRelease, error)
	
	// GetReleasesUpToVersion returns all releases from oldest up to the specified version
	// Note: Returned pointers should be treated as read-only to maintain data integrity
	GetReleasesUpToVersion(ctx context.Context, targetVersion string) ([]*GoRelease, error)
	
	// GetOldestVersion returns the oldest available version
	GetOldestVersion(ctx context.Context) (string, error)
	
	// GetLatestVersion returns the latest available version
	GetLatestVersion(ctx context.Context) (string, error)
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
	GetFeaturesForVersion(ctx context.Context, targetVersion string, packageName string) (*FeatureResponse, error)
}

// ResponseFormatter handles formatting of responses
type ResponseFormatter interface {
	// FormatAsText formats a FeatureResponse as human-readable text
	FormatAsText(response *FeatureResponse, version string, packageName string) string
}