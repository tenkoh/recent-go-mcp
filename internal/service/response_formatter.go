package service

import (
	"slices"
	"strings"

	"github.com/tenkoh/recent-go-mcp/internal/domain"
)

// DefaultResponseFormatter implements ResponseFormatter
type DefaultResponseFormatter struct {
	comparator domain.VersionComparator
}

// NewResponseFormatter creates a new response formatter
func NewResponseFormatter(comparator domain.VersionComparator) domain.ResponseFormatter {
	return &DefaultResponseFormatter{
		comparator: comparator,
	}
}

// FormatAsText formats a FeatureResponse as LLM-readable text
func (f *DefaultResponseFormatter) FormatAsText(response *domain.FeatureResponse, version string, packageName string) string {
	if len(response.Changes) == 0 && len(response.PackageInfo) == 0 {
		return "No Go features found for your project (Go " + response.ToVersion + ")."
	}

	// Use strings.Builder for efficient string construction
	var builder strings.Builder
	builder.Grow(2048) // Pre-allocate reasonable buffer size

	// Write header
	builder.WriteString("Go Features Available in Your Project (Go ")
	builder.WriteString(response.ToVersion)
	builder.WriteString(")\n\n")

	// Write summary
	builder.WriteString("Summary: ")
	builder.WriteString(response.Summary)
	builder.WriteString("\n\n")

	// Get sorted versions for chronological display using slices
	versions := make([]string, 0, len(response.VersionChanges))
	for version := range response.VersionChanges {
		versions = append(versions, version)
	}
	f.sortVersions(versions)

	// Display features by version (oldest to newest)
	for _, version := range versions {
		versionChanges := response.VersionChanges[version]
		versionPackages := response.VersionPackages[version]

		// Skip if no relevant changes for this version
		if len(versionChanges) == 0 && len(versionPackages) == 0 {
			continue
		}

		// Write version header
		builder.WriteString("Go ")
		builder.WriteString(version)
		builder.WriteString(" Features:\n\n")

		// Show general changes for this version
		if len(versionChanges) > 0 {
			builder.WriteString("Language & Runtime Changes:\n")
			for _, change := range versionChanges {
				builder.WriteString("- ")
				builder.WriteString(change.Category)
				builder.WriteString(" (")
				builder.WriteString(change.Impact)
				builder.WriteString("): ")
				builder.WriteString(change.Description)
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}

		// Show package changes for this version
		if len(versionPackages) > 0 {
			builder.WriteString("Standard Library Updates:\n")

			for pkg, changes := range versionPackages {
				if packageName != "" && pkg != packageName {
					continue // Skip if filtering for specific package
				}

				if packageName == "" {
					builder.WriteString("Package ")
					builder.WriteString(pkg)
					builder.WriteString(":\n")
				}

				for _, change := range changes {
					builder.WriteString("- ")
					if change.Function != "" {
						builder.WriteString(change.Function)
						builder.WriteString(" (")
						builder.WriteString(change.Impact)
						builder.WriteString("): ")
					} else {
						builder.WriteString("(")
						builder.WriteString(change.Impact)
						builder.WriteString("): ")
					}
					builder.WriteString(change.Description)
					builder.WriteString("\n")

					if change.Example != "" {
						builder.WriteString("  Example: ")
						builder.WriteString(change.Example)
						builder.WriteString("\n")
					}
				}
				builder.WriteString("\n")
			}
		}

		builder.WriteString("\n")
	}

	builder.WriteString("Note: These are all the Go features available in your project version. Use them to write modern, efficient Go code.\n")

	return builder.String()
}

// sortVersions sorts versions using the version comparator with modern slices
func (f *DefaultResponseFormatter) sortVersions(versions []string) {
	slices.SortFunc(versions, func(a, b string) int {
		return f.comparator.Compare(a, b)
	})
}
