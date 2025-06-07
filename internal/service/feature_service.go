package service

import (
	"context"
	"slices"
	"strconv"
	
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
func (s *DefaultFeatureService) GetFeaturesForVersion(ctx context.Context, targetVersion string, packageName string) (*domain.FeatureResponse, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// Validate input
	if targetVersion == "" {
		return nil, domain.NewValidationError("GetFeaturesForVersion", "target version cannot be empty", nil)
	}
	
	// Get releases up to target version
	availableReleases, err := s.repository.GetReleasesUpToVersion(ctx, targetVersion)
	if err != nil {
		return nil, domain.NewServiceError("GetFeaturesForVersion", "failed to get releases up to version", err).
			WithContext("targetVersion", targetVersion)
	}
	
	if len(availableReleases) == 0 {
		return nil, domain.NewNotFoundError("GetFeaturesForVersion", "no releases found up to version").
			WithContext("targetVersion", targetVersion)
	}
	
	// Find oldest version
	oldestVersion, err := s.repository.GetOldestVersion(ctx)
	if err != nil {
		return nil, domain.NewServiceError("GetFeaturesForVersion", "failed to get oldest version", err)
	}
	
	// Build response
	response := &domain.FeatureResponse{
		FromVersion: oldestVersion,
		ToVersion:   targetVersion,
		Changes:     make([]domain.Change, 0),
		PackageInfo: make(map[string][]domain.PackageChange),
	}
	
	// Collect all changes from available releases using modern Go patterns
	allChanges := make(map[string][]domain.Change)
	allPackageInfo := make(map[string]map[string][]domain.PackageChange)
	
	for _, release := range availableReleases {
		// Group changes by version using slices.Clone for safety
		allChanges[release.Version] = slices.Clone(release.Changes)
		
		// Group package changes by version
		allPackageInfo[release.Version] = make(map[string][]domain.PackageChange)
		
		if packageName != "" {
			// Filter for specific package
			if pkgChanges, exists := release.Packages[packageName]; exists {
				allPackageInfo[release.Version][packageName] = slices.Clone(pkgChanges)
			}
		} else {
			// Include all packages using maps.Copy for efficiency
			for pkg, changes := range release.Packages {
				allPackageInfo[release.Version][pkg] = slices.Clone(changes)
			}
		}
	}
	
	// Flatten for JSON response using modern slice operations
	totalChangesEstimate := 0
	for _, versionChanges := range allChanges {
		totalChangesEstimate += len(versionChanges)
	}
	response.Changes = make([]domain.Change, 0, totalChangesEstimate)
	
	// Collect all changes efficiently
	for _, versionChanges := range allChanges {
		response.Changes = append(response.Changes, versionChanges...)
	}
	
	// Collect package information efficiently
	for _, versionPackages := range allPackageInfo {
		for pkg, changes := range versionPackages {
			existing, exists := response.PackageInfo[pkg]
			if !exists {
				response.PackageInfo[pkg] = slices.Clone(changes)
			} else {
				response.PackageInfo[pkg] = append(existing, changes...)
			}
		}
	}
	
	// Store version-specific data for formatted output
	response.VersionChanges = allChanges
	response.VersionPackages = allPackageInfo
	
	// Generate summary using modern string formatting
	summary := s.generateSummary(targetVersion, oldestVersion, packageName, response)
	response.Summary = summary
	
	return response, nil
}

// generateSummary creates an appropriate summary message
func (s *DefaultFeatureService) generateSummary(targetVersion, oldestVersion, packageName string, response *domain.FeatureResponse) string {
	totalChanges := len(response.Changes)
	totalPackages := len(response.PackageInfo)
	
	if packageName != "" {
		if totalPackages > 0 {
			return "Features available for package '" + packageName + "' in your Go " + targetVersion + " project (from Go " + oldestVersion + ")"
		}
		return "No features found for package '" + packageName + "' in your Go " + targetVersion + " project"
	}
	
	// Use more efficient string building for complex formatting
	if totalChanges == 0 && totalPackages == 0 {
		return "No Go features found in your Go " + targetVersion + " project"
	}
	
	// Build summary with available data
	summary := "All Go features available in your project (Go " + targetVersion + ")"
	if totalChanges > 0 || totalPackages > 0 {
		summary += ": " + strconv.Itoa(totalChanges) + " changes across " + strconv.Itoa(totalPackages) + " packages from Go " + oldestVersion
	}
	
	return summary
}