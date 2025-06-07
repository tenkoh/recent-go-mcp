# Recent Go MCP Server

An MCP (Model Context Protocol) server that provides Go language updates and best practices to LLM coding agents. This helps agents avoid using outdated Go patterns and leverage the latest language features.

## Features

- 🔄 **Version Updates**: Get comprehensive updates from any Go version to the latest
- 📦 **Package-Specific**: Filter updates for specific standard library packages
- 🎯 **Context-Aware**: Categorized changes (language, runtime, toolchain, performance)
- 📚 **Rich Information**: Includes examples, impact assessment, and best practices
- 🚀 **Single Binary**: All release data embedded for easy deployment

## Installation

```bash
# Build the binary
go build -o recent-go-mcp

# Or install directly
go install github.com/tenkoh/recent-go-mcp@latest
```

## Usage

The server implements the Model Context Protocol and can be used with any MCP-compatible client.

### Tool: `go-updates`

Get information about Go language updates and best practices.

**Parameters:**
- `version` (required): Go version to check updates from (e.g., "1.21", "1.22", "1.23")
- `package` (optional): Specific standard library package to filter updates (e.g., "net/http", "context")

### Examples

#### Get all updates from Go 1.21
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "go-updates",
    "arguments": {
      "version": "1.21"
    }
  }
}
```

#### Get net/http specific updates from Go 1.22
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "go-updates",
    "arguments": {
      "version": "1.22",
      "package": "net/http"
    }
  }
}
```

### Response Format

The tool returns both human-readable text and structured JSON data:

- **Summary**: Overview of changes and their impact
- **General Changes**: Language, runtime, and toolchain improvements
- **Package Updates**: Specific updates per standard library package
- **Examples**: Code examples showing new usage patterns
- **Impact Indicators**: Visual icons showing the type of change (🆕 new, ✨ enhancement, ⚡ performance, etc.)

## Integration

### With Claude Desktop

Add to your MCP configuration:

```json
{
  "mcpServers": {
    "recent-go-mcp": {
      "command": "/path/to/recent-go-mcp"
    }
  }
}
```

### With Other MCP Clients

The server uses stdio transport and follows the MCP specification. It can be integrated with any MCP-compatible LLM client.

## Data Coverage

Currently includes:
- Go 1.21 (August 2023)
- Go 1.22 (February 2024) 
- Go 1.23 (August 2024)

The embedded data covers:
- Language feature additions and changes
- Standard library package updates
- Runtime improvements
- Toolchain enhancements
- Performance optimizations
- Best practice recommendations

## Development

```bash
# Run tests
go test ./...

# Build
go build -o recent-go-mcp

# Test manually
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./recent-go-mcp
```

## License

[Add your license here]