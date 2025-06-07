package version

import (
	"go/version"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// GoVersionComparator implements version comparison using Go's official go/version package
type GoVersionComparator struct{}

// NewSemanticVersionComparator creates a new Go version comparator
// Note: Renamed to maintain interface compatibility, but now uses go/version internally
func NewSemanticVersionComparator() domain.VersionComparator {
	return &GoVersionComparator{}
}

// Compare compares two Go version strings using the official go/version package
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func (c *GoVersionComparator) Compare(v1, v2 string) int {
	// Convert to go/version format (e.g., "1.22" -> "go1.22")
	goV1 := normalizeGoVersion(v1)
	goV2 := normalizeGoVersion(v2)
	
	// Use the official Go version comparison
	return version.Compare(goV1, goV2)
}

// normalizeGoVersion converts version strings to go/version package format
// Examples: "1.22" -> "go1.22", "go1.22" -> "go1.22", "1.22.1" -> "go1.22.1"
func normalizeGoVersion(v string) string {
	if v == "" {
		return "go1.0"
	}
	
	// If it already starts with "go", return as-is
	if len(v) >= 2 && v[:2] == "go" {
		return v
	}
	
	// Add "go" prefix to version number
	return "go" + v
}