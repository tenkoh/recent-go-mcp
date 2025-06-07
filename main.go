package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer("recent-go-mcp", "1.0.0", 
		server.WithToolCapabilities(false))

	// Define the go-updates tool
	goUpdatesTool := mcp.NewTool("go-updates",
		mcp.WithDescription("Get information about Go language features and best practices available in your project's Go version. Shows all features from the oldest supported version up to your current version, helping LLM coding agents use appropriate Go patterns."),
		mcp.WithString("version", 
			mcp.Required(),
			mcp.Description("Go version your project is currently using (e.g., '1.21', '1.22', '1.23')")),
		mcp.WithString("package", 
			mcp.Description("Optional: filter features for a specific standard library package (e.g., 'net/http', 'context')")))

	// Add tool handler
	s.AddTool(goUpdatesTool, handleGoUpdates)

	// Start server
	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

func handleGoUpdates(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	
	// Extract version argument (required)
	versionArg, exists := args["version"]
	if !exists {
		return mcp.NewToolResultError("version argument is required"), nil
	}
	
	version, ok := versionArg.(string)
	if !ok {
		return mcp.NewToolResultError("version must be a string"), nil
	}
	
	// Extract package argument (optional)
	var packageName string
	if pkgArg, exists := args["package"]; exists {
		if pkg, ok := pkgArg.(string); ok {
			packageName = pkg
		}
	}
	
	// Get features
	updates, err := getFeaturesForVersion(version, packageName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting updates: %v", err)), nil
	}
	
	// Format response as JSON
	responseJSON, err := json.MarshalIndent(updates, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error formatting response: %v", err)), nil
	}
	
	// Create detailed text response
	textResponse := formatUpdatesAsText(updates, version, packageName)
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(textResponse),
			mcp.NewTextContent(fmt.Sprintf("\n\n--- JSON Response ---\n%s", string(responseJSON))),
		},
	}, nil
}

func formatUpdatesAsText(updates *UpdateResponse, fromVersion, packageName string) string {
	if len(updates.Changes) == 0 && len(updates.PackageInfo) == 0 {
		return fmt.Sprintf("✅ No Go features found for your project (Go %s).", updates.ToVersion)
	}
	
	text := fmt.Sprintf("🔄 **Go Features Available in Your Project (Go %s)**\n\n", updates.ToVersion)
	text += fmt.Sprintf("📋 **Summary**: %s\n\n", updates.Summary)
	
	// Get sorted versions for chronological display
	var versions []string
	for version := range updates.VersionChanges {
		versions = append(versions, version)
	}
	sort.Strings(versions) // Simple string sort works for our version format
	
	// Display features by version (oldest to newest)
	for _, version := range versions {
		versionChanges := updates.VersionChanges[version]
		versionPackages := updates.VersionPackages[version]
		
		// Skip if no relevant changes for this version
		if len(versionChanges) == 0 && len(versionPackages) == 0 {
			continue
		}
		
		text += fmt.Sprintf("## 📦 **Go %s Features**\n\n", version)
		
		// Show general changes for this version
		if len(versionChanges) > 0 {
			text += "### 🚀 **Language & Runtime Changes**\n"
			for _, change := range versionChanges {
				icon := getImpactIcon(change.Impact)
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
					icon := getImpactIcon(change.Impact)
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

func getImpactIcon(impact string) string {
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

func init() {
	// Ensure we have release data loaded
	if len(goReleases) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: No Go release data loaded\n")
	} else {
		fmt.Fprintf(os.Stderr, "Loaded %d Go releases\n", len(goReleases))
	}
}