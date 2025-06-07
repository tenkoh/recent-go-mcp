package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer("recent-go-mcp", "1.0.0", 
		server.WithToolCapabilities(false))

	// Define the go-updates tool
	goUpdatesTool := mcp.NewTool("go-updates",
		mcp.WithDescription("Get information about Go language updates and best practices for a specific version. Helps LLM coding agents avoid using outdated Go patterns and leverage new features."),
		mcp.WithString("version", 
			mcp.Required(),
			mcp.Description("Go version to check updates from (e.g., '1.21', '1.22', '1.23')")),
		mcp.WithString("package", 
			mcp.Description("Optional: specific standard library package to get updates for (e.g., 'net/http', 'context')")))

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
	
	// Get updates
	updates, err := getUpdatesForVersion(version, packageName)
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
		return fmt.Sprintf("✅ Go %s is up to date! No newer releases available.", fromVersion)
	}
	
	text := fmt.Sprintf("🔄 **Go Updates from %s to %s**\n\n", fromVersion, updates.ToVersion)
	text += fmt.Sprintf("📋 **Summary**: %s\n\n", updates.Summary)
	
	if len(updates.Changes) > 0 {
		text += "## 🚀 **General Changes**\n\n"
		for _, change := range updates.Changes {
			icon := getImpactIcon(change.Impact)
			text += fmt.Sprintf("- %s **%s**: %s\n", icon, change.Category, change.Description)
		}
		text += "\n"
	}
	
	if len(updates.PackageInfo) > 0 {
		if packageName != "" {
			text += fmt.Sprintf("## 📦 **Package Updates: %s**\n\n", packageName)
		} else {
			text += "## 📦 **Package Updates**\n\n"
		}
		
		for pkg, changes := range updates.PackageInfo {
			if packageName == "" {
				text += fmt.Sprintf("### `%s`\n", pkg)
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
	
	text += "---\n"
	text += "💡 **Tip**: Use these updates to modernize your code and leverage the latest Go features!\n"
	
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