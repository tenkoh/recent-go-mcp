package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"reflect"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tenkoh/recent-go-mcp/internal/domain"
	"github.com/tenkoh/recent-go-mcp/internal/service"
	"github.com/tenkoh/recent-go-mcp/internal/storage"
	"github.com/tenkoh/recent-go-mcp/internal/version"
)

// Version of the MCP server
const Version = "0.2.0"

// Embed all Go release data files
//
//go:embed data/releases/*.json
var releasesFS embed.FS

// MCPServer wraps the dependencies for the MCP server
type MCPServer struct {
	featureService domain.FeatureService
	formatter      domain.ResponseFormatter
}

// NewMCPServer creates a new MCP server with dependencies initialized and tools registered
func NewMCPServer() (*server.MCPServer, error) {
	// Initialize dependencies
	comparator := version.NewSemanticVersionComparator()

	repo, err := storage.NewEmbeddedReleaseRepository(releasesFS, comparator)
	if err != nil {
		return nil, err
	}

	featureService := service.NewFeatureService(repo, comparator)
	formatter := service.NewResponseFormatter(comparator)

	// Create the wrapper for dependency injection
	mcpWrapper := &MCPServer{
		featureService: featureService,
		formatter:      formatter,
	}

	// Create MCP server
	s := server.NewMCPServer("recent-go-mcp", Version,
		server.WithToolCapabilities(false))

	// Define the go-updates tool
	goUpdatesTool := mcp.NewTool("go-updates",
		mcp.WithDescription("Get comprehensive Go language features and best practices for your project version in structured Markdown format. Supports Go 1.13-1.25, displaying all available features chronologically to help LLM coding agents use modern Go patterns and standard library functions efficiently."),
		mcp.WithString("version",
			mcp.Required(),
			mcp.Description("Go version your project is currently using (supported: '1.13' through '1.25', e.g., '1.21', '1.22', '1.23', '1.25')")),
		mcp.WithString("package",
			mcp.Description("Optional: filter features for a specific standard library package (e.g., 'net/http', 'context', 'slices', 'maps')")))

	// Add tool handler
	s.AddTool(goUpdatesTool, mcpWrapper.handleGoUpdates)

	return s, nil
}

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Initializing recent-go-mcp server",
		"component", "recent-go-mcp",
		"version", Version,
		"supportedGoVersions", "1.13-1.25",
		"architecture", "clean-architecture-with-DI")

	// Create MCP server with dependencies and tools
	mcpServer, err := NewMCPServer()
	if err != nil {
		logger.Error("Failed to create MCP server", "error", err)
		os.Exit(1)
	}
	logger.Info("MCP server created with dependencies and tools registered")

	// Start server
	logger.Info("Starting MCP server")
	if err := server.ServeStdio(mcpServer); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func (m *MCPServer) handleGoUpdates(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger := slog.Default()
	args := request.GetArguments()

	logger.Debug("Processing go-updates request", "args", args)

	// Extract version argument (required)
	versionArg, exists := args["version"]
	if !exists {
		logger.Warn("Missing required version argument")
		return mcp.NewToolResultError("version argument is required"), nil
	}

	version, ok := versionArg.(string)
	if !ok {
		logger.Warn("Invalid version argument type", "type", typeof(versionArg))
		return mcp.NewToolResultError("version must be a string"), nil
	}

	// Extract package argument (optional)
	var packageName string
	if pkgArg, exists := args["package"]; exists {
		if pkg, ok := pkgArg.(string); ok {
			packageName = pkg
		}
	}

	logger.Info("Processing feature request",
		"version", version,
		"package", packageName,
		"hasPackageFilter", packageName != "")

	// Get features using the service with context
	response, err := m.featureService.GetFeaturesForVersion(ctx, version, packageName)
	if err != nil {
		logger.Error("Failed to get features",
			"error", err,
			"version", version,
			"package", packageName)

		// Check if it's a structured error and provide better error messages
		if domain.IsNotFoundError(err) {
			return mcp.NewToolResultError("Version not found: " + version), nil
		}
		if domain.IsValidationError(err) {
			return mcp.NewToolResultError("Invalid input: " + err.Error()), nil
		}

		return mcp.NewToolResultError("Error getting features: " + err.Error()), nil
	}

	logger.Debug("Features retrieved successfully",
		"changesCount", len(response.Changes),
		"packagesCount", len(response.PackageInfo))

	// Create detailed markdown response using formatter
	markdownResponse := m.formatter.FormatAsText(response, version, packageName)

	logger.Info("Request processed successfully",
		"version", version,
		"package", packageName,
		"responseLength", len(markdownResponse))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(markdownResponse),
		},
	}, nil
}

// typeof returns the type name of a value for logging
func typeof(v any) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).String()
}
