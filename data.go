package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed data/releases/go1.23.json
var go123Data []byte

//go:embed data/releases/go1.22.json
var go122Data []byte

//go:embed data/releases/go1.21.json
var go121Data []byte

var goReleases []GoRelease

func init() {
	// Load all embedded release files
	releaseFiles := map[string][]byte{
		"1.23": go123Data,
		"1.22": go122Data,
		"1.21": go121Data,
	}
	
	goReleases = make([]GoRelease, 0, len(releaseFiles))
	
	for version, data := range releaseFiles {
		var release GoRelease
		if err := json.Unmarshal(data, &release); err != nil {
			panic(fmt.Sprintf("failed to load release data for Go %s: %v", version, err))
		}
		goReleases = append(goReleases, release)
	}
	
	// Sort releases by version in descending order (newest first)
	sort.Slice(goReleases, func(i, j int) bool {
		return compareVersions(goReleases[i].Version, goReleases[j].Version) > 0
	})
}

// compareVersions compares two version strings (e.g., "1.23" vs "1.22")
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	// Simple version comparison for Go versions like "1.23"
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		var n1, n2 int
		fmt.Sscanf(parts1[i], "%d", &n1)
		fmt.Sscanf(parts2[i], "%d", &n2)
		
		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}
	
	if len(parts1) > len(parts2) {
		return 1
	} else if len(parts1) < len(parts2) {
		return -1
	}
	
	return 0
}

// getFeaturesForVersion returns all features available from the oldest version up to the specified version
func getFeaturesForVersion(targetVersion string, packageName string) (*UpdateResponse, error) {
	var targetRelease *GoRelease
	var availableReleases []GoRelease
	
	// Find the target version
	for _, release := range goReleases {
		if release.Version == targetVersion {
			targetRelease = &release
			break
		}
	}
	
	if targetRelease == nil {
		return nil, fmt.Errorf("version %s not found", targetVersion)
	}
	
	// Collect all releases from oldest to target version (inclusive)
	for _, release := range goReleases {
		if compareVersions(release.Version, targetVersion) <= 0 {
			availableReleases = append(availableReleases, release)
		}
	}
	
	// Sort releases from oldest to newest
	sort.Slice(availableReleases, func(i, j int) bool {
		return compareVersions(availableReleases[i].Version, availableReleases[j].Version) < 0
	})
	
	if len(availableReleases) == 0 {
		return nil, fmt.Errorf("no releases found up to version %s", targetVersion)
	}
	
	// Find oldest version
	oldestVersion := availableReleases[0].Version
	
	// Build response
	response := &UpdateResponse{
		FromVersion: oldestVersion,
		ToVersion:   targetVersion,
		Changes:     []Change{},
		PackageInfo: make(map[string][]PackageChange),
	}
	
	// Collect all changes from available releases
	allChanges := make(map[string][]Change)
	allPackageInfo := make(map[string]map[string][]PackageChange)
	
	for _, release := range availableReleases {
		// Group changes by version
		allChanges[release.Version] = release.Changes
		
		// Group package changes by version
		allPackageInfo[release.Version] = make(map[string][]PackageChange)
		
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
	response.Changes = []Change{} // Will be used for JSON, but formatted text will be different
	response.PackageInfo = make(map[string][]PackageChange)
	
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