# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server that provides Go language updates and best practices to LLM coding agents in structured Markdown format. The server helps agents avoid using outdated Go patterns and leverage new language features efficiently.

**Current Status**: v0.2.0 (latest release)
**Supported Go Versions**: 1.13 through 1.25 (13 versions)
**Architecture**: Clean architecture with Go 1.25 best practices

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
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.25"}}}' | ./recent-go-mcp
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.22", "package": "slices"}}}' | ./recent-go-mcp
```

## Architecture

This project follows clean architecture principles with dependency injection and Go 1.25 best practices for improved testability and maintainability.

### Core Components

- **main.go**: MCP server implementation with `NewMCPServer()` factory function and version constant management
- **internal/domain/**: Domain models, interfaces, and structured error types with proper error wrapping
- **internal/service/**: Business logic layer with context propagation and modern Go patterns
- **internal/storage/**: Data access layer using embedded filesystem and modern slice operations
- **internal/version/**: Version comparison using Go 1.22+ `go/version` package (simplified from manual parsing)
- **data/releases/**: Individual JSON files for Go 1.13-1.25 (embedded via go:embed)

### Key Design Patterns

- **Clean Architecture**: Clear separation of concerns with domain, service, and infrastructure layers
- **Dependency Injection**: All dependencies injected through interfaces for testability
- **Repository Pattern**: Abstract data access through ReleaseRepository interface
- **Embedded Data**: Uses `//go:embed` in main.go to include JSON data in the binary
- **MCP Protocol**: Implements Model Context Protocol for LLM integration
- **SOLID Principles**: Single responsibility, dependency inversion, and open/closed principles
- **Go 1.25 Best Practices**: Modern error handling, context propagation, structured logging, and efficient operations

### Adding or Updating Go Version Changelogs

Follow this checklist whenever you add a new release JSON (for example, `go1.25.json`) or revise an existing one:

1. **Source data first**: Review the official Go release notes, language spec updates, and relevant package documentation before writing. Avoid speculation—if a feature or API name is not in the upstream material, leave it out.
2. **Use the house prompt**: `make_changelog_prompt.md` documents the required structure, tone, and quality bar. Keep the JSON schema identical (ordering of top-level keys, field casing, impact taxonomy).
3. **Write trustworthy examples**: Every code sample should compile with the stated API. Prefer minimal but realistic snippets showing imports/aliases when needed.
4. **Cross-check coverage**: Ensure all major language, runtime, toolchain, and platform callouts from the release notes are represented. Include breaking changes explicitly with `impact: "breaking"`.
5. **Update metadata**: After adding a new version file, synchronize version ranges in `README.md`, `main.go`, and any other human-facing docs so they advertise the correct highest Go version.
6. **Validation pass**: Run `go test ./...` to make sure the embedded data still parses and the MCP surface works. If you add or rename files, run `go fmt ./...` to keep formatting consistent.
7. **No drive-by edits**: Keep unrelated refactors or copy edits out of changelog-focused PRs unless explicitly requested.

Document any nuances discovered while following this flow back into this section so future tasks get faster and more accurate.

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
- **Supports Go 1.13 through 1.25**: Comprehensive version coverage (13 Go versions)
- **Version Parameter**: Required Go version (e.g., "1.22", "1.23", "1.24", "1.25")
- **Package Filtering**: Optional package name for focused results (e.g., "net/http", "slices", "maps")
- **Markdown Response Format**: Returns structured Markdown for optimal LLM consumption
- **Modern Go Features**: Highlights `slices`, `maps`, `log/slog`, `go/version`, and other Go 1.21+ features
- **Best Practices**: Includes examples, impact indicators, and upgrade recommendations
- **Efficient Output**: ~70% size reduction compared to previous dual-format approach

## Go 1.25 Modernization Features

This codebase demonstrates Go 1.25 best practices:

- **Structured Error Handling**: Custom error types with proper wrapping and context
- **Context Propagation**: `context.Context` throughout the service layer
- **Modern Slice Operations**: `slices.SortFunc`, `slices.Clone`, `slices.IndexFunc`, etc.
- **Efficient String Building**: `strings.Builder` for performance
- **Structured Logging**: `log/slog` with JSON output and contextual fields
- **Official Version Comparison**: `go/version.Compare` instead of manual parsing
- **Clean Architecture**: Dependency injection with interfaces for testability
