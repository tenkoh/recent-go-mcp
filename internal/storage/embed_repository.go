package storage

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// EmbeddedReleaseRepository implements ReleaseRepository using embedded JSON files
type EmbeddedReleaseRepository struct {
	releases   []domain.GoRelease
	comparator domain.VersionComparator
	fs         embed.FS
}

// NewEmbeddedReleaseRepository creates a new repository with embedded data
func NewEmbeddedReleaseRepository(fs embed.FS, comparator domain.VersionComparator) (domain.ReleaseRepository, error) {
	repo := &EmbeddedReleaseRepository{
		fs:         fs,
		comparator: comparator,
	}
	
	if err := repo.loadReleases(); err != nil {
		return nil, fmt.Errorf("failed to load embedded releases: %w", err)
	}
	
	return repo, nil
}

// loadReleases loads and sorts all embedded release data
func (r *EmbeddedReleaseRepository) loadReleases() error {
	// Read all JSON files from the embedded filesystem
	entries, err := r.fs.ReadDir("data/releases")
	if err != nil {
		return fmt.Errorf("failed to read embedded release directory: %w", err)
	}
	
	r.releases = make([]domain.GoRelease, 0, len(entries))
	
	for _, entry := range entries {
		if !entry.IsDir() && path.Ext(entry.Name()) == ".json" {
			filePath := path.Join("data/releases", entry.Name())
			
			data, err := r.fs.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read embedded file %s: %w", filePath, err)
			}
			
			var release domain.GoRelease
			if err := json.Unmarshal(data, &release); err != nil {
				return fmt.Errorf("failed to unmarshal release data from %s: %w", filePath, err)
			}
			
			r.releases = append(r.releases, release)
		}
	}
	
	if len(r.releases) == 0 {
		return fmt.Errorf("no release data found in embedded filesystem")
	}
	
	// Sort releases by version in descending order (newest first)
	sort.Slice(r.releases, func(i, j int) bool {
		return r.comparator.Compare(r.releases[i].Version, r.releases[j].Version) > 0
	})
	
	return nil
}

// GetAllReleases returns all available Go releases
func (r *EmbeddedReleaseRepository) GetAllReleases() ([]domain.GoRelease, error) {
	// Return a copy to prevent external modification
	releases := make([]domain.GoRelease, len(r.releases))
	copy(releases, r.releases)
	return releases, nil
}

// GetReleaseByVersion returns a specific release by version
func (r *EmbeddedReleaseRepository) GetReleaseByVersion(version string) (*domain.GoRelease, error) {
	for _, release := range r.releases {
		if release.Version == version {
			// Return a copy to prevent external modification
			releaseCopy := release
			return &releaseCopy, nil
		}
	}
	return nil, fmt.Errorf("release not found for version %s", version)
}

// GetReleasesUpToVersion returns all releases from oldest up to the specified version
func (r *EmbeddedReleaseRepository) GetReleasesUpToVersion(targetVersion string) ([]domain.GoRelease, error) {
	// First, verify the target version exists
	if _, err := r.GetReleaseByVersion(targetVersion); err != nil {
		return nil, err
	}
	
	var result []domain.GoRelease
	
	// Collect all releases up to and including the target version
	for _, release := range r.releases {
		if r.comparator.Compare(release.Version, targetVersion) <= 0 {
			result = append(result, release)
		}
	}
	
	// Sort from oldest to newest for chronological display
	sort.Slice(result, func(i, j int) bool {
		return r.comparator.Compare(result[i].Version, result[j].Version) < 0
	})
	
	return result, nil
}

// GetOldestVersion returns the oldest available version
func (r *EmbeddedReleaseRepository) GetOldestVersion() (string, error) {
	if len(r.releases) == 0 {
		return "", fmt.Errorf("no releases available")
	}
	
	oldest := r.releases[0].Version
	for _, release := range r.releases {
		if r.comparator.Compare(release.Version, oldest) < 0 {
			oldest = release.Version
		}
	}
	
	return oldest, nil
}

// GetLatestVersion returns the latest available version
func (r *EmbeddedReleaseRepository) GetLatestVersion() (string, error) {
	if len(r.releases) == 0 {
		return "", fmt.Errorf("no releases available")
	}
	
	latest := r.releases[0].Version
	for _, release := range r.releases {
		if r.comparator.Compare(release.Version, latest) > 0 {
			latest = release.Version
		}
	}
	
	return latest, nil
}