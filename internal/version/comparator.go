package version

import (
	"fmt"
	"strings"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// SemanticVersionComparator implements version comparison for semantic versioning
type SemanticVersionComparator struct{}

// NewSemanticVersionComparator creates a new semantic version comparator
func NewSemanticVersionComparator() domain.VersionComparator {
	return &SemanticVersionComparator{}
}

// Compare compares two version strings (e.g., "1.23" vs "1.22")
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func (c *SemanticVersionComparator) Compare(v1, v2 string) int {
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