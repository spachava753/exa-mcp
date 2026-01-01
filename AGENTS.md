# MCP Server Template

This is a Go template for building Model Context Protocol (MCP) servers.

## Project Structure

```
.
├── main.go              # Entry point, server startup
├── internal/
│   ├── server.go        # Server configuration, tool registration
│   └── echo.go          # Example tool implementation
├── go.mod
└── README.md
```

## Adding a New Tool

1. Create a new file in `internal/` (e.g., `my_tool.go`)

2. Define the tool structure:
```go
// Args struct with jsonschema tags for parameter descriptions
type MyToolArgs struct {
    Param string `json:"param" jsonschema:"Description of the parameter"`
}

// Output struct for structured responses
type MyToolOutput struct {
    Result string `json:"result"`
}

// Tool definition with metadata
var MyToolDef = mcp.Tool{
    Name:        "my_tool",
    Description: "What this tool does",
}

// Handler function
func MyTool(ctx context.Context, req *mcp.CallToolRequest, args MyToolArgs) (*mcp.CallToolResult, MyToolOutput, error) {
    // Implementation
}
```

3. Register in `server.go`:
```go
mcp.AddTool(server, &MyToolDef, MyTool)
```

## Key Patterns

- **Tool handlers** return `(*mcp.CallToolResult, OutputType, error)`
- **Error responses**: Set `IsError: true` on CallToolResult
- **Structured output**: The second return value enables typed responses
- Use `jsonschema` tags on args for MCP client documentation
