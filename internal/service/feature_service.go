package service

import (
	"fmt"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// DefaultFeatureService implements FeatureService
type DefaultFeatureService struct {
	repository domain.ReleaseRepository
	comparator domain.VersionComparator
}

// NewFeatureService creates a new feature service
func NewFeatureService(repository domain.ReleaseRepository, comparator domain.VersionComparator) domain.FeatureService {
	return &DefaultFeatureService{
		repository: repository,
		comparator: comparator,
	}
}

// GetFeaturesForVersion returns all features available from the oldest version up to the specified version
func (s *DefaultFeatureService) GetFeaturesForVersion(targetVersion string, packageName string) (*domain.FeatureResponse, error) {
	// Get releases up to target version
	availableReleases, err := s.repository.GetReleasesUpToVersion(targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get releases up to version %s: %w", targetVersion, err)
	}
	
	if len(availableReleases) == 0 {
		return nil, fmt.Errorf("no releases found up to version %s", targetVersion)
	}
	
	// Find oldest version
	oldestVersion, err := s.repository.GetOldestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest version: %w", err)
	}
	
	// Build response
	response := &domain.FeatureResponse{
		FromVersion: oldestVersion,
		ToVersion:   targetVersion,
		Changes:     []domain.Change{},
		PackageInfo: make(map[string][]domain.PackageChange),
	}
	
	// Collect all changes from available releases
	allChanges := make(map[string][]domain.Change)
	allPackageInfo := make(map[string]map[string][]domain.PackageChange)
	
	for _, release := range availableReleases {
		// Group changes by version
		allChanges[release.Version] = release.Changes
		
		// Group package changes by version
		allPackageInfo[release.Version] = make(map[string][]domain.PackageChange)
		
		if packageName != "" {
			// Filter for specific package
			if pkgChanges, exists := release.Packages[packageName]; exists {
				allPackageInfo[release.Version][packageName] = pkgChanges
			}
		} else {
			// Include all packages
			for pkg, changes := range release.Packages {
				allPackageInfo[release.Version][pkg] = changes
			}
		}
	}
	
	// Store structured data for formatted output
	response.Changes = []domain.Change{} // Will be used for JSON, but formatted text will be different
	response.PackageInfo = make(map[string][]domain.PackageChange)
	
	// Flatten for JSON response
	for _, versionChanges := range allChanges {
		response.Changes = append(response.Changes, versionChanges...)
	}
	
	for _, versionPackages := range allPackageInfo {
		for pkg, changes := range versionPackages {
			response.PackageInfo[pkg] = append(response.PackageInfo[pkg], changes...)
		}
	}
	
	// Store version-specific data for formatted output
	response.VersionChanges = allChanges
	response.VersionPackages = allPackageInfo
	
	// Generate summary
	totalChanges := len(response.Changes)
	totalPackages := len(response.PackageInfo)
	
	if packageName != "" {
		if totalPackages > 0 {
			response.Summary = fmt.Sprintf("Features available for package '%s' in your Go %s project (from Go %s)",
				packageName, targetVersion, oldestVersion)
		} else {
			response.Summary = fmt.Sprintf("No features found for package '%s' in your Go %s project",
				packageName, targetVersion)
		}
	} else {
		response.Summary = fmt.Sprintf("All Go features available in your project (Go %s): %d changes across %d packages from Go %s",
			targetVersion, totalChanges, totalPackages, oldestVersion)
	}
	
	return response, nil
}