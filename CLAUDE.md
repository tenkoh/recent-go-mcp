# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server that provides Go language updates and best practices to LLM coding agents. The server helps agents avoid using outdated Go patterns and leverage new language features.

## Common Commands

- `go build -o recent-go-mcp` - Build the MCP server binary
- `go test ./...` - Run all tests
- `go mod tidy` - Clean up dependencies
- `go fmt ./...` - Format code
- `go vet ./...` - Run static analysis

### Testing the MCP Server

```bash
# Test tools list
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./recent-go-mcp

# Test go-updates tool
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.21"}}}' | ./recent-go-mcp

# Test with package filter
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.22", "package": "net/http"}}}' | ./recent-go-mcp
```

## Architecture

### Core Components

- **main.go**: MCP server implementation with tool handlers
- **types.go**: Data structures for Go releases and updates
- **data.go**: Data loading, version comparison, and filtering logic
- **releases.json**: Embedded Go release data (embedded via go:embed)

### Key Design Patterns

- **Embedded Data**: Uses `//go:embed` to include JSON data in the binary
- **MCP Protocol**: Implements Model Context Protocol for LLM integration
- **Version Comparison**: Smart semantic version comparison for Go releases
- **Filtering**: Package-specific filtering for targeted updates

### Adding New Go Versions

1. Update `releases.json` with new release data
2. Follow the existing JSON structure with version, changes, and package updates
3. The `init()` function automatically loads and sorts releases

### Data Structure

- Releases are sorted by version (newest first)
- Each release contains general changes and package-specific updates
- Changes are categorized by impact: "new", "enhancement", "performance", "breaking", "deprecation"

### MCP Integration

The server provides a single tool `go-updates` that:
- Takes a Go version (required) and optional package name
- Returns formatted text + JSON response
- Includes examples, impact indicators, and best practices