# Exa MCP Server

MCP server exposing Exa AI's search capabilities as tools.

## Project Structure

```
.
├── main.go                    # Entry point, server startup
├── internal/
│   ├── server.go              # Server configuration, tool registration
│   ├── exa_client.go          # Exa API client initialization
│   ├── search.go              # exa_search tool
│   ├── answer.go              # exa_answer tool
│   ├── find_similar.go        # exa_find_similar tool
│   ├── get_contents.go        # exa_get_contents tool
│   ├── generate.go            # go:generate directive for SDK
│   ├── openapi-spec.yaml      # Exa OpenAPI spec
│   ├── exasdk/                # Generated Exa SDK (via ogen)
│   └── cmd/fixspec/           # Tool to fix OpenAPI spec issues
├── go.mod
└── README.md
```

## Adding a New Tool

1. Create a new file in `internal/` (e.g., `my_tool.go`)

2. Define the tool structure:
```go
type MyToolArgs struct {
    Param    string   `json:"param" jsonschema:"Description of the parameter"`
    Optional string   `json:"optional,omitempty" jsonschema:"Optional field description"`
}

type MyToolOutput struct {
    Result string `json:"result"`
}

var MyToolDef = mcp.Tool{
    Name:        "my_tool",
    Description: "What this tool does",
}

func MyTool(ctx context.Context, req *mcp.CallToolRequest, args MyToolArgs) (*mcp.CallToolResult, MyToolOutput, error) {
    // Implementation
}
```

3. Register in `server.go`:
```go
mcp.AddTool(server, &MyToolDef, MyTool)
```

## jsonschema Tags

The MCP SDK uses `github.com/google/jsonschema-go` for schema inference. The `jsonschema` struct tag has specific rules:

**The tag value is ONLY used as the description.** It must NOT start with `WORD=` pattern.

```go
// CORRECT - plain description text
Query string `json:"query" jsonschema:"The search query"`

// WRONG - will panic at startup
Query string `json:"query" jsonschema:"required,description=The search query"`
Query string `json:"query" jsonschema:"description=The search query"`
Query string `json:"query" jsonschema:"enum=a,enum=b,description=The query"`
```

**Required fields** are inferred from the `json` tag:
- No `omitempty` → required
- Has `omitempty` → optional

**Enums** are not supported via struct tags. Describe valid values in the description text instead.
