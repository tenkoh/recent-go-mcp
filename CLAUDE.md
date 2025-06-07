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

This project follows clean architecture principles with dependency injection for improved testability and maintainability.

### Core Components

- **main.go**: MCP server implementation with dependency injection and embedded data
- **internal/domain/**: Domain models and interfaces defining contracts
- **internal/service/**: Business logic layer with feature retrieval and formatting
- **internal/storage/**: Data access layer with embedded repository implementation
- **internal/version/**: Version comparison utilities with comprehensive testing
- **data/releases/**: Individual JSON files for each Go version (embedded via go:embed)

### Key Design Patterns

- **Clean Architecture**: Clear separation of concerns with domain, service, and infrastructure layers
- **Dependency Injection**: All dependencies injected through interfaces for testability
- **Repository Pattern**: Abstract data access through ReleaseRepository interface
- **Embedded Data**: Uses `//go:embed` in main.go to include JSON data in the binary
- **MCP Protocol**: Implements Model Context Protocol for LLM integration
- **SOLID Principles**: Single responsibility, dependency inversion, and open/closed principles

### Adding New Go Versions

1. Create a new JSON file in `data/releases/` following the pattern `go{version}.json`
2. The embedded repository automatically discovers and loads all `.json` files in the releases directory
3. No code changes required - the system dynamically loads new versions
4. Follow the existing JSON structure with version, changes, and package updates
5. Version sorting and comparison handled automatically by the version comparator

### Testing

- **Unit Tests**: Comprehensive test coverage for business logic using mocks
- **Version Comparator**: Fully tested semantic version comparison
- **Service Layer**: Testable business logic with dependency injection
- **Repository Pattern**: Mockable data access for isolated testing

### Data Structure

- Releases are automatically sorted by semantic version (newest first)
- Each release contains general changes and package-specific updates
- Changes are categorized by impact: "new", "enhancement", "performance", "breaking", "deprecation"
- Version comparison uses proper semantic versioning logic

### MCP Integration

The server provides a single tool `go-updates` that:
- Takes a Go version (required) and optional package name
- Returns formatted text + JSON response
- Includes examples, impact indicators, and best practices