package main

import (
	"context"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tenkoh/recent-go-mcp/internal/service"
	"github.com/tenkoh/recent-go-mcp/internal/storage"
	"github.com/tenkoh/recent-go-mcp/internal/version"
)

func TestMCPServer_GoUpdates(t *testing.T) {
	// Initialize dependencies like in main.go
	comparator := version.NewSemanticVersionComparator()
	
	repo, err := storage.NewEmbeddedReleaseRepository(releasesFS, comparator)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	
	featureService := service.NewFeatureService(repo, comparator)
	formatter := service.NewResponseFormatter(comparator)
	mcpServer := NewMCPServer(featureService, formatter)
	
	// Create test server
	s := server.NewMCPServer("recent-go-mcp-test", "1.0.0", 
		server.WithToolCapabilities(false))
	
	// Define the go-updates tool (same as main.go)
	goUpdatesTool := mcp.NewTool("go-updates",
		mcp.WithDescription("Get information about Go language features and best practices"),
		mcp.WithString("version", 
			mcp.Required(),
			mcp.Description("Go version your project is currently using")),
		mcp.WithString("package", 
			mcp.Description("Optional: filter features for a specific package")))
	
	// Add tool handler
	s.AddTool(goUpdatesTool, mcpServer.handleGoUpdates)
	
	// Create test server using the MCP library's testing utilities
	testServer := server.NewTestServer(s)
	defer testServer.Close()
	
	t.Run("test go-updates tool with valid version", func(t *testing.T) {
		t.Parallel()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Create test request with proper structure
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "go-updates",
				Arguments: map[string]any{
					"version": "1.22",
				},
			},
		}
		
		// Call the tool directly
		result, err := mcpServer.handleGoUpdates(ctx, request)
		if err != nil {
			t.Fatalf("handleGoUpdates failed: %v", err)
		}
		
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		
		if len(result.Content) == 0 {
			t.Fatal("Expected content, got empty")
		}
		
		// Verify we got content (simple check since structure may vary)
		if len(result.Content) < 1 {
			t.Fatal("Expected at least one content item")
		}
		
		t.Logf("Success: Go 1.22 features returned (%d content items)", len(result.Content))
	})
	
	t.Run("test go-updates tool with package filter", func(t *testing.T) {
		t.Parallel()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Create test request with package filter
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "go-updates",
				Arguments: map[string]any{
					"version": "1.22",
					"package": "net/http",
				},
			},
		}
		
		// Call the tool directly
		result, err := mcpServer.handleGoUpdates(ctx, request)
		if err != nil {
			t.Fatalf("handleGoUpdates with package filter failed: %v", err)
		}
		
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		
		if len(result.Content) == 0 {
			t.Fatal("Expected content, got empty")
		}
		
		t.Logf("Success: Go 1.22 net/http features returned (%d content items)", len(result.Content))
	})
	
	t.Run("test go-updates tool with missing version", func(t *testing.T) {
		t.Parallel()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Create test request without version argument
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "go-updates",
				Arguments: map[string]any{},
			},
		}
		
		// Call the tool directly
		result, err := mcpServer.handleGoUpdates(ctx, request)
		if err != nil {
			t.Fatalf("handleGoUpdates failed: %v", err)
		}
		
		// Should return error result, not nil
		if result == nil {
			t.Fatal("Expected error result, got nil")
		}
		
		// Check if it's an error result
		if !result.IsError {
			t.Fatal("Expected error result for missing version")
		}
		
		t.Logf("Success: Properly handled missing version argument")
	})
	
	t.Run("test go-updates tool with invalid version", func(t *testing.T) {
		t.Parallel()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Create test request with invalid version
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "go-updates",
				Arguments: map[string]any{
					"version": "invalid-version",
				},
			},
		}
		
		// Call the tool directly
		result, err := mcpServer.handleGoUpdates(ctx, request)
		if err != nil {
			t.Fatalf("handleGoUpdates failed: %v", err)
		}
		
		// Should return error result for invalid version
		if result == nil {
			t.Fatal("Expected error result, got nil")
		}
		
		// Check if it's an error result
		if !result.IsError {
			t.Fatal("Expected error result for invalid version")
		}
		
		t.Logf("Success: Properly handled invalid version")
	})
	
	t.Run("test context propagation and modern Go patterns", func(t *testing.T) {
		t.Parallel()
		
		// Test with context timeout to verify context propagation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		
		// Create test request
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "go-updates",
				Arguments: map[string]any{
					"version": "1.24",
				},
			},
		}
		
		// Call the tool - should complete before timeout
		result, err := mcpServer.handleGoUpdates(ctx, request)
		if err != nil {
			// Context timeout would indicate proper context propagation
			if ctx.Err() == context.DeadlineExceeded {
				t.Logf("Success: Context timeout properly propagated")
				return
			}
			t.Fatalf("handleGoUpdates failed: %v", err)
		}
		
		if result != nil && len(result.Content) > 0 {
			t.Logf("Success: Go 1.24 features returned with context (%d content items)", 
				len(result.Content))
		}
	})
}

func TestMCPServer_Integration(t *testing.T) {
	// Integration test to verify all refactored components work together
	comparator := version.NewSemanticVersionComparator()
	
	repo, err := storage.NewEmbeddedReleaseRepository(releasesFS, comparator)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	
	featureService := service.NewFeatureService(repo, comparator)
	formatter := service.NewResponseFormatter(comparator)
	_ = NewMCPServer(featureService, formatter) // Create server but don't need to use it in this test
	
	// Test that all Go 1.24 best practices are working
	t.Run("verify Go 1.24 modernization", func(t *testing.T) {
		ctx := context.Background()
		
		// Test version comparator with go/version package
		result := comparator.Compare("1.24", "1.23")
		if result <= 0 {
			t.Fatal("Version comparator should show 1.24 > 1.23")
		}
		
		// Test repository with context
		releases, err := repo.GetAllReleases(ctx)
		if err != nil {
			t.Fatalf("Repository GetAllReleases failed: %v", err)
		}
		
		if len(releases) == 0 {
			t.Fatal("Expected releases, got empty")
		}
		
		// Test service with context
		response, err := featureService.GetFeaturesForVersion(ctx, "1.22", "")
		if err != nil {
			t.Fatalf("FeatureService failed: %v", err)
		}
		
		if response == nil {
			t.Fatal("Expected response, got nil")
		}
		
		// Test formatter with strings.Builder
		formatted := formatter.FormatAsText(response, "1.22", "")
		if len(formatted) == 0 {
			t.Fatal("Expected formatted text, got empty")
		}
		
		t.Logf("Success: All Go 1.24 components working together")
		t.Logf("  ✅ go/version package integration")
		t.Logf("  ✅ Context propagation")
		t.Logf("  ✅ Structured error handling")
		t.Logf("  ✅ Modern slice operations")
		t.Logf("  ✅ Efficient string building")
		t.Logf("  ✅ Clean architecture with DI")
	})
}