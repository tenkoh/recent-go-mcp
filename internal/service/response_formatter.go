package service

import (
	"fmt"
	"sort"
	
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

// FormatAsText formats a FeatureResponse as human-readable text
func (f *DefaultResponseFormatter) FormatAsText(response *domain.FeatureResponse, version string, packageName string) string {
	if len(response.Changes) == 0 && len(response.PackageInfo) == 0 {
		return fmt.Sprintf("✅ No Go features found for your project (Go %s).", response.ToVersion)
	}
	
	text := fmt.Sprintf("🔄 **Go Features Available in Your Project (Go %s)**\n\n", response.ToVersion)
	text += fmt.Sprintf("📋 **Summary**: %s\n\n", response.Summary)
	
	// Get sorted versions for chronological display
	var versions []string
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
		
		text += fmt.Sprintf("## 📦 **Go %s Features**\n\n", version)
		
		// Show general changes for this version
		if len(versionChanges) > 0 {
			text += "### 🚀 **Language & Runtime Changes**\n"
			for _, change := range versionChanges {
				icon := f.getImpactIcon(change.Impact)
				text += fmt.Sprintf("- %s **%s**: %s\n", icon, change.Category, change.Description)
			}
			text += "\n"
		}
		
		// Show package changes for this version
		if len(versionPackages) > 0 {
			text += "### 📚 **Standard Library Updates**\n"
			
			for pkg, changes := range versionPackages {
				if packageName != "" && pkg != packageName {
					continue // Skip if filtering for specific package
				}
				
				if packageName == "" {
					text += fmt.Sprintf("#### `%s`\n", pkg)
				}
				
				for _, change := range changes {
					icon := f.getImpactIcon(change.Impact)
					if change.Function != "" {
						text += fmt.Sprintf("- %s **%s**: %s\n", icon, change.Function, change.Description)
					} else {
						text += fmt.Sprintf("- %s %s\n", icon, change.Description)
					}
					
					if change.Example != "" {
						text += fmt.Sprintf("  ```go\n  %s\n  ```\n", change.Example)
					}
				}
				text += "\n"
			}
		}
		
		text += "\n"
	}
	
	text += "---\n"
	text += "💡 **Tip**: These are all the Go features available in your project version. Use them to write modern, efficient Go code!\n"
	
	return text
}

// sortVersions sorts versions using the version comparator
func (f *DefaultResponseFormatter) sortVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		return f.comparator.Compare(versions[i], versions[j]) < 0
	})
}

// getImpactIcon returns an appropriate emoji icon for the impact type
func (f *DefaultResponseFormatter) getImpactIcon(impact string) string {
	switch impact {
	case "new":
		return "🆕"
	case "enhancement":
		return "✨"
	case "performance":
		return "⚡"
	case "breaking":
		return "⚠️"
	case "deprecation":
		return "🗑️"
	default:
		return "📝"
	}
}