package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tenkoh/recent-go-mcp/internal/domain"
	"github.com/tenkoh/recent-go-mcp/internal/service"
	"github.com/tenkoh/recent-go-mcp/internal/storage"
	"github.com/tenkoh/recent-go-mcp/internal/version"
)

// Embed all Go release data files
//go:embed data/releases/*.json
var releasesFS embed.FS

// MCPServer wraps the dependencies for the MCP server
type MCPServer struct {
	featureService domain.FeatureService
	formatter      domain.ResponseFormatter
}

// NewMCPServer creates a new MCP server with injected dependencies
func NewMCPServer(featureService domain.FeatureService, formatter domain.ResponseFormatter) *MCPServer {
	return &MCPServer{
		featureService: featureService,
		formatter:      formatter,
	}
}

func main() {
	// Initialize dependencies
	comparator := version.NewSemanticVersionComparator()
	
	repo, err := storage.NewEmbeddedReleaseRepository(releasesFS, comparator)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	
	featureService := service.NewFeatureService(repo, comparator)
	formatter := service.NewResponseFormatter(comparator)
	
	mcpServer := NewMCPServer(featureService, formatter)
	
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
	s.AddTool(goUpdatesTool, mcpServer.handleGoUpdates)

	// Start server
	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

func (m *MCPServer) handleGoUpdates(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	
	// Get features using the service
	response, err := m.featureService.GetFeaturesForVersion(version, packageName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting features: %v", err)), nil
	}
	
	// Format response as JSON
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error formatting response: %v", err)), nil
	}
	
	// Create detailed text response using formatter
	textResponse := m.formatter.FormatAsText(response, version, packageName)
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(textResponse),
			mcp.NewTextContent(fmt.Sprintf("\n\n--- JSON Response ---\n%s", string(responseJSON))),
		},
	}, nil
}

func init() {
	// Log initialization information
	fmt.Fprintf(os.Stderr, "Initializing recent-go-mcp server with dependency injection architecture\n")
}