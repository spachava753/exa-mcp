# My MCP Server

A Model Context Protocol (MCP) server template for Go.

## Quick Start

1. **Clone and customize:**
```bash
# Update module path
go mod edit -module github.com/yourusername/your-mcp-server

# Update imports in main.go and internal/
# Update server name in internal/server.go
```

2. **Build and run:**
```bash
go build -o my-mcp-server .
./my-mcp-server
```

3. **Configure with Claude Desktop** (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "my-mcp-server": {
      "command": "/path/to/my-mcp-server"
    }
  }
}
```

## Adding Tools

See [AGENTS.md](AGENTS.md) for detailed instructions on adding new tools.

## Example Tool

The template includes an `echo` tool that demonstrates the basic pattern:

```go
// Define args with jsonschema tags
type EchoArgs struct {
    Message string `json:"message" jsonschema:"The message to echo back"`
}

// Define tool metadata
var EchoToolDef = mcp.Tool{
    Name:        "echo",
    Description: "Echoes back the provided message",
}

// Implement handler
func EchoTool(ctx context.Context, req *mcp.CallToolRequest, args EchoArgs) (*mcp.CallToolResult, EchoOutput, error) {
    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: args.Message}},
    }, EchoOutput{Echo: args.Message}, nil
}
```

## License

MIT
