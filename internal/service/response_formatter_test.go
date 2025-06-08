package service

import (
	"testing"

	"github.com/tenkoh/recent-go-mcp/internal/domain"
	"github.com/tenkoh/recent-go-mcp/internal/version"
)

func TestResponseFormatter_FullStringComparison(t *testing.T) {
	comparator := version.NewSemanticVersionComparator()
	formatter := NewResponseFormatter(comparator)

	t.Run("empty response", func(t *testing.T) {
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
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
		}
	})

	t.Run("language changes only", func(t *testing.T) {
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

		expected := `Go Features Available in Your Project (Go 1.22)

Summary: Go 1.22 introduces new language features

Go 1.22 Features:

Language & Runtime Changes:
- language (new): for-range over integers


Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.
`

		if result != expected {
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
		}
	})

	t.Run("package changes only", func(t *testing.T) {
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
				"1.21": {},
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

		expected := `Go Features Available in Your Project (Go 1.21)

Summary: Go 1.21 adds new standard library packages

Go 1.21 Features:

Standard Library Updates:
Package slices:
- Sort (new): sorts a slice


Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.
`

		if result != expected {
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
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

		expected := `Go Features Available in Your Project (Go 1.22)

Summary: Language and library improvements

Go 1.22 Features:

Language & Runtime Changes:
- language (new): for-range over integers

Standard Library Updates:
Package net/http:
- ServeMux (enhancement): enhanced routing


Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.
`

		if result != expected {
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
		}
	})

	t.Run("package filtering", func(t *testing.T) {
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

		expected := `Go Features Available in Your Project (Go 1.22)

Summary: Multiple packages available

Go 1.22 Features:

Standard Library Updates:
- ServeMux (enhancement): enhanced routing


Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.
`

		if result != expected {
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
		}
	})

	t.Run("code examples", func(t *testing.T) {
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

		expected := `Go Features Available in Your Project (Go 1.21)

Summary: Package with code examples

Go 1.21 Features:

Standard Library Updates:
Package slices:
- Sort (new): sorts a slice
  Example: slices.Sort([]int{3,1,2})


Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.
`

		if result != expected {
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
		}
	})

	t.Run("multiple versions chronological order", func(t *testing.T) {
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

		expected := `Go Features Available in Your Project (Go 1.23)

Summary: Features from Go 1.21 to 1.23

Go 1.21 Features:

Language & Runtime Changes:
- language (new): 1.21 feature


Go 1.22 Features:

Language & Runtime Changes:
- language (new): 1.22 feature


Go 1.23 Features:

Language & Runtime Changes:
- language (new): 1.23 feature


Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.
`

		if result != expected {
			t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
		}
	})
}
