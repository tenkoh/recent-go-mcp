# Recent Go MCP Server

An MCP (Model Context Protocol) server that provides comprehensive Go language updates and best practices to LLM coding agents. This helps agents avoid using outdated Go patterns and leverage modern language features across 12 Go versions.

**NOTICE**
This repository is experimental and for personal use. (Of course, anyone can use this repository.) LLM coding agent is fully utilized to write code, test and data about Go's features in each release. Some mistakes would be contained. I hope any issues or PRs to improve this repository.

**NOTICE**
The output of this tool might be oversize with some MCP hosts like Claude Code.
I will update this tool with pagenation feature soon.

## Features

- 🔄 **Comprehensive Version Coverage**: Supports Go 1.13 through 1.24 (12 versions)
- 📦 **Package-Specific Filtering**: Get updates for specific standard library packages (net/http, slices, maps, log/slog, etc.)
- 📚 **Rich Information**: Includes examples, impact assessment, and upgrade recommendations
- 🚀 **Single Binary**: All release data embedded using go:embed for easy deployment

## Integration

### With Claude Desktop

Add to your MCP configuration:

```json
{
  "mcpServers": {
    "recent-go-mcp": {
      "command": "go",
      "args":[
        "run",
        "github.com/tenkoh/recent-go-mcp@latest"
      ]
    }
  }
}
```

You can install the command using `go install`, then register the executable binary's path to MCP configuration.

### With Other MCP Clients

The server uses stdio transport and follows the MCP specification. It can be integrated with any MCP-compatible LLM client.

## Usage

The server implements the Model Context Protocol and can be used with any MCP-compatible client.

### Tool: `go-updates`

Get information about Go language updates and best practices.

**Parameters:**
- `version` (required): Go version to check updates from (supported: "1.13" through "1.24")
- `package` (optional): Specific standard library package to filter updates (e.g., "net/http", "slices", "maps", "log/slog")

### Examples

#### Get all updates from Go 1.21 to Go 1.24
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

#### Get slices package specific updates from Go 1.20
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "go-updates",
    "arguments": {
      "version": "1.20",
      "package": "slices"
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


## Data Coverage

**Comprehensive Go Version Support:**
- Go 1.13 through Go 1.24 (12 versions)
- Complete coverage from legacy versions to the latest features

**The embedded data covers:**
- **Language Features**: New syntax, constructs, and language improvements
- **Standard Library Updates**: Package additions, enhancements, and new APIs
- **Modern Packages**: slices, maps, log/slog, go/version, and other Go 1.21+ additions
- **Runtime Improvements**: Performance optimizations and memory management
- **Toolchain Enhancements**: Build system, module system, and developer tooling updates
- **Best Practice Recommendations**: Modern patterns and upgrade guidance

## Contribution
Contributions are really welcomed. Please make an issue or a pull request casually.

## Development

```bash
# Run all tests (includes MCP protocol tests)
go test ./...

# Run MCP server integration tests specifically
go test -v -run TestMCPServer

# Build the server
go build -o recent-go-mcp

# Test manually with JSON-RPC
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./recent-go-mcp
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "go-updates", "arguments": {"version": "1.24"}}}' | ./recent-go-mcp
```

## License

MIT

## Author

tenkoh (using Claude Sonnet 4)