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

// getUpdatesForVersion returns all changes from the specified version to the latest
func getUpdatesForVersion(targetVersion string, packageName string) (*UpdateResponse, error) {
	var fromRelease *GoRelease
	var relevantReleases []GoRelease
	
	// Find the target version and collect all newer releases
	for _, release := range goReleases {
		if compareVersions(release.Version, targetVersion) > 0 {
			relevantReleases = append(relevantReleases, release)
		} else if release.Version == targetVersion {
			fromRelease = &release
			break
		}
	}
	
	if fromRelease == nil {
		return nil, fmt.Errorf("version %s not found", targetVersion)
	}
	
	if len(relevantReleases) == 0 {
		return &UpdateResponse{
			FromVersion: targetVersion,
			ToVersion:   targetVersion,
			Summary:     "No updates available - you're using the latest version!",
			Changes:     []Change{},
		}, nil
	}
	
	// Build response
	response := &UpdateResponse{
		FromVersion: targetVersion,
		ToVersion:   relevantReleases[0].Version, // Latest version
		Changes:     []Change{},
		PackageInfo: make(map[string][]PackageChange),
	}
	
	// Collect all changes from newer releases
	for _, release := range relevantReleases {
		response.Changes = append(response.Changes, release.Changes...)
		
		// If specific package requested, filter package changes
		if packageName != "" {
			if pkgChanges, exists := release.Packages[packageName]; exists {
				response.PackageInfo[packageName] = append(response.PackageInfo[packageName], pkgChanges...)
			}
		} else {
			// Include all package changes
			for pkg, changes := range release.Packages {
				response.PackageInfo[pkg] = append(response.PackageInfo[pkg], changes...)
			}
		}
	}
	
	// Generate summary
	if packageName != "" {
		if len(response.PackageInfo[packageName]) > 0 {
			response.Summary = fmt.Sprintf("Found %d updates for package '%s' from Go %s to %s",
				len(response.PackageInfo[packageName]), packageName, targetVersion, response.ToVersion)
		} else {
			response.Summary = fmt.Sprintf("No updates found for package '%s' from Go %s to %s",
				packageName, targetVersion, response.ToVersion)
		}
	} else {
		response.Summary = fmt.Sprintf("Found %d general changes and %d package updates from Go %s to %s",
			len(response.Changes), len(response.PackageInfo), targetVersion, response.ToVersion)
	}
	
	return response, nil
}