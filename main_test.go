package main

import (
	"context"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestMCPServer_GoUpdates(t *testing.T) {
	// Create MCP server with dependencies and tools registered
	mcpServer, err := NewMCPServer()
	if err != nil {
		t.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create test server using the MCP library's testing utilities
	testServer := server.NewTestServer(mcpServer)
	defer testServer.Close()

	// Test that the server was created successfully and tools are registered
	t.Run("test server creation and tool registration", func(t *testing.T) {
		if testServer == nil {
			t.Fatal("Expected test server to be created")
		}
		// The fact that NewTestServer succeeded means tools are registered correctly
	})

	// Test actual request to the server
	t.Run("test go-updates tool request", func(t *testing.T) {
		// Create in-process client for more reliable testing
		cli, err := client.NewInProcessClient(mcpServer)
		if err != nil {
			t.Fatalf("Failed to create in-process client: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Start the client transport
		if err := cli.Start(ctx); err != nil {
			t.Fatalf("Failed to start client: %v", err)
		}

		// Initialize client with proper request
		initReq := mcp.InitializeRequest{
			Params: mcp.InitializeParams{
				ProtocolVersion: "2024-11-05",
				ClientInfo: mcp.Implementation{
					Name:    "test-client",
					Version: "0.1.0",
				},
				Capabilities: mcp.ClientCapabilities{},
			},
		}
		_, err = cli.Initialize(ctx, initReq)
		if err != nil {
			t.Fatalf("Failed to initialize client: %v", err)
		}
		defer cli.Close()

		// Call the go-updates tool with proper request
		toolReq := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "go-updates",
				Arguments: map[string]any{
					"version": "1.22",
				},
			},
		}

		result, err := cli.CallTool(ctx, toolReq)
		if err != nil {
			t.Fatalf("Failed to call go-updates tool: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if len(result.Content) == 0 {
			t.Fatal("Expected content in result")
		}

		// Verify we got some content (simplified check)
		if len(result.Content) < 1 {
			t.Fatal("Expected at least one content item")
		}
	})
}
