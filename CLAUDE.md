# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server that provides Go language updates and best practices to LLM coding agents. The server helps agents avoid using outdated Go patterns and leverage new language features.

**Current Status**: v0.1.0 (release ready)
**Supported Go Versions**: 1.13 through 1.24 (12 versions)
**Architecture**: Clean architecture with Go 1.24 best practices

## Common Commands

- `go build -o recent-go-mcp` - Build the MCP server binary
- `go test ./...` - Run all tests (includes MCP protocol tests)
- `go test -v -run TestMCPServer` - Run MCP server integration tests
- `go mod tidy` - Clean up dependencies
- `go fmt ./...` - Format code
- `go vet ./...` - Run static analysis

### Testing the MCP Server

The project includes comprehensive automated tests using the MCP Go library:

```bash
# Run MCP server tests (recommended)
go test -v -run TestMCPServer

# Manual JSON-RPC testing (if needed)
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./recent-go-mcp
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.24"}}}' | ./recent-go-mcp
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.22", "package": "slices"}}}' | ./recent-go-mcp
```

## Architecture

This project follows clean architecture principles with dependency injection and Go 1.24 best practices for improved testability and maintainability.

### Core Components

- **main.go**: MCP server implementation with `NewMCPServer()` factory function that handles full initialization
- **internal/domain/**: Domain models, interfaces, and structured error types with proper error wrapping
- **internal/service/**: Business logic layer with context propagation and modern Go patterns
- **internal/storage/**: Data access layer using embedded filesystem and modern slice operations
- **internal/version/**: Version comparison using Go 1.22+ `go/version` package (simplified from manual parsing)
- **data/releases/**: Individual JSON files for Go 1.13-1.24 (embedded via go:embed)

### Key Design Patterns

- **Clean Architecture**: Clear separation of concerns with domain, service, and infrastructure layers
- **Dependency Injection**: All dependencies injected through interfaces for testability
- **Repository Pattern**: Abstract data access through ReleaseRepository interface
- **Embedded Data**: Uses `//go:embed` in main.go to include JSON data in the binary
- **MCP Protocol**: Implements Model Context Protocol for LLM integration
- **SOLID Principles**: Single responsibility, dependency inversion, and open/closed principles
- **Go 1.24 Best Practices**: Modern error handling, context propagation, structured logging, and efficient operations

### Adding New Go Versions

1. Create a new JSON file in `data/releases/` following the pattern `go{version}.json`
2. The embedded repository automatically discovers and loads all `.json` files in the releases directory
3. No code changes required - the system dynamically loads new versions
4. Follow the existing JSON structure with version, changes, and package updates
5. Version sorting and comparison handled automatically by the version comparator

### Testing

- **MCP Protocol Tests**: Real MCP client-server communication using `NewInProcessClient`
- **Unit Tests**: Comprehensive test coverage for business logic using mocks
- **Version Comparator**: Fully tested semantic version comparison using `go/version` package
- **Service Layer**: Testable business logic with dependency injection and context propagation
- **Repository Pattern**: Mockable data access for isolated testing
- **Automated Testing**: `go test ./...` runs all tests including MCP integration tests

### Data Structure

- Releases are automatically sorted by semantic version (newest first)
- Each release contains general changes and package-specific updates
- Changes are categorized by impact: "new", "enhancement", "performance", "breaking", "deprecation"
- Version comparison uses proper semantic versioning logic

### MCP Integration

The server provides a single tool `go-updates` that:
- **Supports Go 1.13 through 1.24**: Comprehensive version coverage (12 Go versions)
- **Version Parameter**: Required Go version (e.g., "1.21", "1.22", "1.23", "1.24")
- **Package Filtering**: Optional package name for focused results (e.g., "net/http", "slices", "maps")
- **Dual Response Format**: Returns both formatted text and structured JSON
- **Modern Go Features**: Highlights `slices`, `maps`, `log/slog`, `go/version`, and other Go 1.21+ features
- **Best Practices**: Includes examples, impact indicators, and upgrade recommendations

## Go 1.24 Modernization Features

This codebase demonstrates Go 1.24 best practices:

- **Structured Error Handling**: Custom error types with proper wrapping and context
- **Context Propagation**: `context.Context` throughout the service layer
- **Modern Slice Operations**: `slices.SortFunc`, `slices.Clone`, `slices.IndexFunc`, etc.
- **Efficient String Building**: `strings.Builder` for performance
- **Structured Logging**: `log/slog` with JSON output and contextual fields
- **Official Version Comparison**: `go/version.Compare` instead of manual parsing
- **Clean Architecture**: Dependency injection with interfaces for testability