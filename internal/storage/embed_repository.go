package storage

import (
	"context"
	"encoding/json"
	"io/fs"
	"path"
	"slices"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// FullFS combines all the filesystem interfaces we need
type FullFS interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

// EmbeddedReleaseRepository implements ReleaseRepository using embedded JSON files
type EmbeddedReleaseRepository struct {
	releases   []*domain.GoRelease
	comparator domain.VersionComparator
	fs         FullFS
}

// NewEmbeddedReleaseRepository creates a new repository with embedded data
func NewEmbeddedReleaseRepository(filesystem FullFS, comparator domain.VersionComparator) (domain.ReleaseRepository, error) {
	repo := &EmbeddedReleaseRepository{
		fs:         filesystem,
		comparator: comparator,
	}
	
	if err := repo.loadReleases(); err != nil {
		return nil, domain.NewRepositoryError("initialization", "failed to load embedded releases", err)
	}
	
	return repo, nil
}

// loadReleases loads and sorts all embedded release data
func (r *EmbeddedReleaseRepository) loadReleases() error {
	// Read all JSON files from the embedded filesystem
	entries, err := r.fs.ReadDir("data/releases")
	if err != nil {
		return domain.NewRepositoryError("loadReleases", "failed to read embedded release directory", err)
	}
	
	r.releases = make([]*domain.GoRelease, 0, len(entries))
	
	for _, entry := range entries {
		if !entry.IsDir() && path.Ext(entry.Name()) == ".json" {
			filePath := path.Join("data/releases", entry.Name())
			
			data, err := r.fs.ReadFile(filePath)
			if err != nil {
				return domain.NewRepositoryError("loadReleases", "failed to read embedded file", err).
					WithContext("file", filePath)
			}
			
			var release domain.GoRelease
			if err := json.Unmarshal(data, &release); err != nil {
				return domain.NewRepositoryError("loadReleases", "failed to unmarshal release data", err).
					WithContext("file", filePath)
			}
			
			r.releases = append(r.releases, &release)
		}
	}
	
	if len(r.releases) == 0 {
		return domain.NewRepositoryError("loadReleases", "no release data found in embedded filesystem", nil)
	}
	
	// Sort releases by version in descending order (newest first) using slices.SortFunc
	slices.SortFunc(r.releases, func(a, b *domain.GoRelease) int {
		return -r.comparator.Compare(a.Version, b.Version) // Negative for descending order
	})
	
	return nil
}

// GetAllReleases returns all available Go releases
func (r *EmbeddedReleaseRepository) GetAllReleases(ctx context.Context) ([]*domain.GoRelease, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// Return a copy of the slice using slices.Clone (Go 1.21+)
	return slices.Clone(r.releases), nil
}

// GetReleaseByVersion returns a specific release by version
func (r *EmbeddedReleaseRepository) GetReleaseByVersion(ctx context.Context, version string) (*domain.GoRelease, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// Find release using slices utilities
	idx := slices.IndexFunc(r.releases, func(release *domain.GoRelease) bool {
		return release.Version == version
	})
	
	if idx == -1 {
		return nil, domain.NewNotFoundError("GetReleaseByVersion", "release not found").
			WithContext("version", version)
	}
	
	return r.releases[idx], nil
}

// GetReleasesUpToVersion returns all releases from oldest up to the specified version
func (r *EmbeddedReleaseRepository) GetReleasesUpToVersion(ctx context.Context, targetVersion string) ([]*domain.GoRelease, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// First, verify the target version exists
	if _, err := r.GetReleaseByVersion(ctx, targetVersion); err != nil {
		return nil, err
	}
	
	// Filter releases up to target version using slices utilities
	filtered := make([]*domain.GoRelease, 0, len(r.releases))
	for _, release := range r.releases {
		if r.comparator.Compare(release.Version, targetVersion) <= 0 {
			filtered = append(filtered, release)
		}
	}
	
	// Sort from oldest to newest for chronological display using slices.SortFunc
	slices.SortFunc(filtered, func(a, b *domain.GoRelease) int {
		return r.comparator.Compare(a.Version, b.Version)
	})
	
	return filtered, nil
}

// GetOldestVersion returns the oldest available version
func (r *EmbeddedReleaseRepository) GetOldestVersion(ctx context.Context) (string, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	
	if len(r.releases) == 0 {
		return "", domain.NewRepositoryError("GetOldestVersion", "no releases available", nil)
	}
	
	// Find minimum using slices.MinFunc (more efficient than manual loop)
	oldest := slices.MinFunc(r.releases, func(a, b *domain.GoRelease) int {
		return r.comparator.Compare(a.Version, b.Version)
	})
	
	return oldest.Version, nil
}

// GetLatestVersion returns the latest available version
func (r *EmbeddedReleaseRepository) GetLatestVersion(ctx context.Context) (string, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	
	if len(r.releases) == 0 {
		return "", domain.NewRepositoryError("GetLatestVersion", "no releases available", nil)
	}
	
	// Find maximum using slices.MaxFunc (more efficient than manual loop)
	latest := slices.MaxFunc(r.releases, func(a, b *domain.GoRelease) int {
		return r.comparator.Compare(a.Version, b.Version)
	})
	
	return latest.Version, nil
}