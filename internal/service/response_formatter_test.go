package service

import (
	"strings"
	"testing"
	
	"github.com/tenkoh/recent-go-mcp/internal/domain"
	"github.com/tenkoh/recent-go-mcp/internal/version"
)

func TestResponseFormatter_FormatAsText(t *testing.T) {
	comparator := version.NewSemanticVersionComparator()
	formatter := NewResponseFormatter(comparator)
	
	t.Run("empty response - no changes and no packages", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion:       "1.21",
			Changes:         []domain.Change{},
			PackageInfo:     map[string][]domain.PackageChange{},
			VersionChanges:  map[string][]domain.Change{},
			VersionPackages: map[string]map[string][]domain.PackageChange{},
		}
		
		result := formatter.FormatAsText(response, "1.21", "")
		expected := "No Go features found for your project (Go 1.21)."
		
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
	
	t.Run("single version with language changes only", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.22",
			Summary:   "Go 1.22 introduces new language features",
			Changes: []domain.Change{
				{Category: "language", Description: "for-range over integers", Impact: "new"},
			},
			PackageInfo: map[string][]domain.PackageChange{},
			VersionChanges: map[string][]domain.Change{
				"1.22": {
					{Category: "language", Description: "for-range over integers", Impact: "new"},
				},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.22": {},
			},
		}
		
		result := formatter.FormatAsText(response, "1.22", "")
		
		// Verify header
		if !strings.Contains(result, "Go Features Available in Your Project (Go 1.22)") {
			t.Error("Missing project header")
		}
		
		// Verify summary
		if !strings.Contains(result, "Summary: Go 1.22 introduces new language features") {
			t.Error("Missing summary section")
		}
		
		// Verify version section
		if !strings.Contains(result, "Go 1.22 Features:") {
			t.Error("Missing version features header")
		}
		
		// Verify language changes section
		if !strings.Contains(result, "Language & Runtime Changes:") {
			t.Error("Missing language changes section")
		}
		
		// Verify change formatting
		if !strings.Contains(result, "- language (new): for-range over integers") {
			t.Error("Missing or incorrectly formatted language change")
		}
		
		// Verify note
		if !strings.Contains(result, "Note: These are all the Go features available in your project version") {
			t.Error("Missing closing note")
		}
	})
	
	t.Run("single version with package changes only", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.21",
			Summary:   "Go 1.21 adds new standard library packages",
			Changes:   []domain.Change{},
			PackageInfo: map[string][]domain.PackageChange{
				"slices": {
					{Function: "Sort", Description: "sorts a slice", Impact: "new"},
				},
			},
			VersionChanges: map[string][]domain.Change{
				"1.21": {}, // Need empty slice, not missing key
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.21": {
					"slices": {
						{Function: "Sort", Description: "sorts a slice", Impact: "new"},
					},
				},
			},
		}
		
		result := formatter.FormatAsText(response, "1.21", "")
		
		// Should have standard library section
		if !strings.Contains(result, "Standard Library Updates:") {
			t.Error("Missing standard library updates section")
		}
		
		// Should have package header
		if !strings.Contains(result, "Package slices:") {
			t.Error("Missing package header")
		}
		
		// Should have function change
		if !strings.Contains(result, "- Sort (new): sorts a slice") {
			t.Error("Missing or incorrectly formatted function change")
		}
		
		// Should NOT have language changes section
		if strings.Contains(result, "Language & Runtime Changes:") {
			t.Error("Should not have language changes section when no language changes")
		}
	})
	
	t.Run("multiple versions in chronological order", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.23",
			Summary:   "Features from Go 1.21 to 1.23",
			Changes: []domain.Change{
				{Category: "language", Description: "1.21 feature", Impact: "new"},
				{Category: "language", Description: "1.22 feature", Impact: "new"},
				{Category: "language", Description: "1.23 feature", Impact: "new"},
			},
			PackageInfo: map[string][]domain.PackageChange{},
			VersionChanges: map[string][]domain.Change{
				"1.23": {{Category: "language", Description: "1.23 feature", Impact: "new"}},
				"1.21": {{Category: "language", Description: "1.21 feature", Impact: "new"}},
				"1.22": {{Category: "language", Description: "1.22 feature", Impact: "new"}},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.21": {},
				"1.22": {},
				"1.23": {},
			},
		}
		
		result := formatter.FormatAsText(response, "1.23", "")
		
		// Find positions of version headers
		pos21 := strings.Index(result, "Go 1.21 Features:")
		pos22 := strings.Index(result, "Go 1.22 Features:")
		pos23 := strings.Index(result, "Go 1.23 Features:")
		
		if pos21 == -1 || pos22 == -1 || pos23 == -1 {
			t.Fatal("Missing one or more version headers")
		}
		
		// Verify chronological order (oldest to newest)
		if !(pos21 < pos22 && pos22 < pos23) {
			t.Errorf("Versions not in chronological order: 1.21@%d, 1.22@%d, 1.23@%d", pos21, pos22, pos23)
		}
	})
	
	t.Run("package changes with examples", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.21",
			Summary:   "Package with code examples",
			Changes:   []domain.Change{},
			PackageInfo: map[string][]domain.PackageChange{
				"slices": {
					{Function: "Sort", Description: "sorts a slice", Impact: "new", Example: "slices.Sort([]int{3,1,2})"},
				},
			},
			VersionChanges: map[string][]domain.Change{
				"1.21": {},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.21": {
					"slices": {
						{Function: "Sort", Description: "sorts a slice", Impact: "new", Example: "slices.Sort([]int{3,1,2})"},
					},
				},
			},
		}
		
		result := formatter.FormatAsText(response, "1.21", "")
		
		// Should have example
		if !strings.Contains(result, "Example: slices.Sort([]int{3,1,2})") {
			t.Error("Missing code example")
		}
	})
	
	t.Run("package changes without function name", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.22",
			Summary:   "General package improvements",
			Changes:   []domain.Change{},
			PackageInfo: map[string][]domain.PackageChange{
				"net/http": {
					{Description: "overall performance improvements", Impact: "performance"},
				},
			},
			VersionChanges: map[string][]domain.Change{
				"1.22": {},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.22": {
					"net/http": {
						{Description: "overall performance improvements", Impact: "performance"},
					},
				},
			},
		}
		
		result := formatter.FormatAsText(response, "1.22", "")
		
		// Should format without function name
		if !strings.Contains(result, "- (performance): overall performance improvements") {
			t.Error("Missing or incorrectly formatted general package change")
		}
	})
	
	t.Run("package filtering shows only specified package", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.22",
			Summary:   "Multiple packages available",
			Changes:   []domain.Change{},
			PackageInfo: map[string][]domain.PackageChange{
				"net/http": {
					{Function: "ServeMux", Description: "enhanced routing", Impact: "enhancement"},
				},
				"slices": {
					{Function: "Sort", Description: "sorts a slice", Impact: "new"},
				},
			},
			VersionChanges: map[string][]domain.Change{
				"1.22": {},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.22": {
					"net/http": {
						{Function: "ServeMux", Description: "enhanced routing", Impact: "enhancement"},
					},
					"slices": {
						{Function: "Sort", Description: "sorts a slice", Impact: "new"},
					},
				},
			},
		}
		
		result := formatter.FormatAsText(response, "1.22", "net/http")
		
		// Should include filtered package content
		if !strings.Contains(result, "ServeMux") || !strings.Contains(result, "enhanced routing") {
			t.Error("Missing content for filtered package net/http")
		}
		
		// Should exclude other packages
		if strings.Contains(result, "Sort") || strings.Contains(result, "sorts a slice") {
			t.Error("Found excluded package content when filtering for net/http")
		}
		
		// Should NOT show package headers when filtering for specific package
		if strings.Contains(result, "Package net/http:") {
			t.Error("Package header should not appear when filtering for specific package")
		}
	})
	
	t.Run("mixed language and package changes", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.22",
			Summary:   "Language and library improvements",
			Changes: []domain.Change{
				{Category: "language", Description: "for-range over integers", Impact: "new"},
			},
			PackageInfo: map[string][]domain.PackageChange{
				"net/http": {
					{Function: "ServeMux", Description: "enhanced routing", Impact: "enhancement"},
				},
			},
			VersionChanges: map[string][]domain.Change{
				"1.22": {
					{Category: "language", Description: "for-range over integers", Impact: "new"},
				},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.22": {
					"net/http": {
						{Function: "ServeMux", Description: "enhanced routing", Impact: "enhancement"},
					},
				},
			},
		}
		
		result := formatter.FormatAsText(response, "1.22", "")
		
		// Should have both sections
		if !strings.Contains(result, "Language & Runtime Changes:") {
			t.Error("Missing language changes section")
		}
		
		if !strings.Contains(result, "Standard Library Updates:") {
			t.Error("Missing standard library updates section")
		}
		
		// Verify both types of changes are present
		if !strings.Contains(result, "- language (new): for-range over integers") {
			t.Error("Missing language change")
		}
		
		if !strings.Contains(result, "- ServeMux (enhancement): enhanced routing") {
			t.Error("Missing package change")
		}
	})
	
	t.Run("empty versions are skipped", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.23",
			Summary:   "Some versions have no changes",
			Changes: []domain.Change{
				{Category: "language", Description: "1.21 feature", Impact: "new"},
				{Category: "language", Description: "1.23 feature", Impact: "new"},
			},
			PackageInfo: map[string][]domain.PackageChange{},
			VersionChanges: map[string][]domain.Change{
				"1.21": {{Category: "language", Description: "1.21 feature", Impact: "new"}},
				"1.22": {}, // Empty - should be skipped
				"1.23": {{Category: "language", Description: "1.23 feature", Impact: "new"}},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.21": {},
				"1.22": {}, // Empty - should be skipped
				"1.23": {},
			},
		}
		
		result := formatter.FormatAsText(response, "1.23", "")
		
		// Should contain non-empty versions
		if !strings.Contains(result, "Go 1.21 Features:") {
			t.Error("Missing 1.21 features section")
		}
		
		if !strings.Contains(result, "Go 1.23 Features:") {
			t.Error("Missing 1.23 features section")
		}
		
		// Should NOT contain empty version
		if strings.Contains(result, "Go 1.22 Features:") {
			t.Error("Empty version 1.22 should be skipped")
		}
	})
	
	t.Run("all impact types are formatted correctly", func(t *testing.T) {
		response := &domain.FeatureResponse{
			ToVersion: "1.22",
			Summary:   "Various impact types",
			Changes: []domain.Change{
				{Category: "language", Description: "new feature", Impact: "new"},
				{Category: "runtime", Description: "enhanced feature", Impact: "enhancement"},
				{Category: "performance", Description: "faster execution", Impact: "performance"},
				{Category: "breaking", Description: "breaking change", Impact: "breaking"},
				{Category: "deprecated", Description: "old feature deprecated", Impact: "deprecation"},
			},
			PackageInfo: map[string][]domain.PackageChange{},
			VersionChanges: map[string][]domain.Change{
				"1.22": {
					{Category: "language", Description: "new feature", Impact: "new"},
					{Category: "runtime", Description: "enhanced feature", Impact: "enhancement"},
					{Category: "performance", Description: "faster execution", Impact: "performance"},
					{Category: "breaking", Description: "breaking change", Impact: "breaking"},
					{Category: "deprecated", Description: "old feature deprecated", Impact: "deprecation"},
				},
			},
			VersionPackages: map[string]map[string][]domain.PackageChange{
				"1.22": {},
			},
		}
		
		result := formatter.FormatAsText(response, "1.22", "")
		
		// Verify all impact types are properly formatted
		expectedFormats := []string{
			"- language (new): new feature",
			"- runtime (enhancement): enhanced feature", 
			"- performance (performance): faster execution",
			"- breaking (breaking): breaking change",
			"- deprecated (deprecation): old feature deprecated",
		}
		
		for _, expected := range expectedFormats {
			if !strings.Contains(result, expected) {
				t.Errorf("Missing expected format: %s", expected)
			}
		}
	})
}