package storage

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/tenkoh/recent-go-mcp/internal/version"
)

func TestEmbeddedReleaseRepository_WithMockFS(t *testing.T) {
	// Create mock file system with test data using testing/fstest
	testJSON := `{
		"version": "1.21",
		"release_date": "2023-08-08T00:00:00Z",
		"summary": "Test release",
		"changes": [],
		"packages": {}
	}`

	// Create mock filesystem using fstest.MapFS
	mockFS := fstest.MapFS{
		"data/releases/go1.21.json": &fstest.MapFile{
			Data: []byte(testJSON),
		},
	}

	// Create repository with mock filesystem
	comparator := version.NewSemanticVersionComparator()

	// This demonstrates how we can now inject any fs.FS implementation
	repo, err := NewEmbeddedReleaseRepository(mockFS, comparator)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test the repository
	ctx := context.Background()
	releases, err := repo.GetAllReleases(ctx)
	if err != nil {
		t.Fatalf("Failed to get releases: %v", err)
	}

	if len(releases) != 1 {
		t.Errorf("Expected 1 release, got %d", len(releases))
	}

	if releases[0].Version != "1.21" {
		t.Errorf("Expected version 1.21, got %s", releases[0].Version)
	}
}
