package service

import (
	"fmt"
	"testing"
	"time"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// mockRepository implements domain.ReleaseRepository for testing
type mockRepository struct {
	releases []*domain.GoRelease
}

func (m *mockRepository) GetAllReleases() ([]*domain.GoRelease, error) {
	return m.releases, nil
}

func (m *mockRepository) GetReleaseByVersion(version string) (*domain.GoRelease, error) {
	for _, release := range m.releases {
		if release.Version == version {
			return release, nil
		}
	}
	return nil, fmt.Errorf("version not found: %s", version)
}

func (m *mockRepository) GetReleasesUpToVersion(targetVersion string) ([]*domain.GoRelease, error) {
	var result []*domain.GoRelease
	for _, release := range m.releases {
		if release.Version <= targetVersion { // Simple string comparison for testing
			result = append(result, release)
		}
	}
	return result, nil
}

func (m *mockRepository) GetOldestVersion() (string, error) {
	if len(m.releases) == 0 {
		return "", fmt.Errorf("no releases")
	}
	return m.releases[0].Version, nil
}

func (m *mockRepository) GetLatestVersion() (string, error) {
	if len(m.releases) == 0 {
		return "", fmt.Errorf("no releases")
	}
	return m.releases[len(m.releases)-1].Version, nil
}

// mockComparator implements domain.VersionComparator for testing
type mockComparator struct{}

func (c *mockComparator) Compare(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}
	if v1 < v2 {
		return -1
	}
	return 1
}

func TestDefaultFeatureService_GetFeaturesForVersion(t *testing.T) {
	// Setup test data
	testReleases := []*domain.GoRelease{
		&domain.GoRelease{
			Version:     "1.21",
			ReleaseDate: time.Date(2023, 8, 8, 0, 0, 0, 0, time.UTC),
			Summary:     "Go 1.21 release",
			Changes: []domain.Change{
				{
					Category:    "language",
					Description: "New built-in functions",
					Impact:      "new",
				},
			},
			Packages: map[string][]domain.PackageChange{
				"slices": {
					{
						Function:    "Sort",
						Description: "Sort a slice",
						Impact:      "new",
					},
				},
			},
		},
		&domain.GoRelease{
			Version:     "1.22",
			ReleaseDate: time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),
			Summary:     "Go 1.22 release",
			Changes: []domain.Change{
				{
					Category:    "language",
					Description: "For-range over integers",
					Impact:      "new",
				},
			},
			Packages: map[string][]domain.PackageChange{
				"net/http": {
					{
						Function:    "ServeMux",
						Description: "Enhanced routing",
						Impact:      "enhancement",
					},
				},
			},
		},
	}
	
	repo := &mockRepository{releases: testReleases}
	comparator := &mockComparator{}
	service := NewFeatureService(repo, comparator)
	
	t.Run("successful feature retrieval", func(t *testing.T) {
		response, err := service.GetFeaturesForVersion("1.22", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if response.ToVersion != "1.22" {
			t.Errorf("expected ToVersion 1.22, got %s", response.ToVersion)
		}
		
		if len(response.Changes) != 2 {
			t.Errorf("expected 2 changes, got %d", len(response.Changes))
		}
		
		if len(response.PackageInfo) != 2 {
			t.Errorf("expected 2 packages, got %d", len(response.PackageInfo))
		}
	})
	
	t.Run("package filtering", func(t *testing.T) {
		response, err := service.GetFeaturesForVersion("1.22", "net/http")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(response.PackageInfo) != 1 {
			t.Errorf("expected 1 package after filtering, got %d", len(response.PackageInfo))
		}
		
		if _, exists := response.PackageInfo["net/http"]; !exists {
			t.Error("expected net/http package in response")
		}
	})
}