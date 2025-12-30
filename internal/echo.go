package internal

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// EchoArgs defines the input parameters for the echo tool.
// Use jsonschema tags to provide descriptions for MCP clients.
type EchoArgs struct {
	Message string `json:"message" jsonschema:"The message to echo back"`
}

// EchoOutput defines the structured output of the echo tool.
// This enables type-safe responses and better client integration.
type EchoOutput struct {
	Echo string `json:"echo"`
}

// EchoToolDef defines the tool metadata exposed to MCP clients.
var EchoToolDef = mcp.Tool{
	Name:        "echo",
	Description: "Echoes back the provided message. A simple example tool demonstrating the MCP tool pattern.",
	Annotations: &mcp.ToolAnnotations{
		Title: "Echo",
	},
}

// EchoTool implements the echo tool logic.
// The function signature follows the MCP SDK pattern:
//   - ctx: context for cancellation
//   - req: the raw MCP request (rarely needed)
//   - args: typed arguments parsed from the request
// Returns: (CallToolResult, structured output, error)
func EchoTool(ctx context.Context, req *mcp.CallToolRequest, args EchoArgs) (*mcp.CallToolResult, EchoOutput, error) {
	output := EchoOutput{Echo: args.Message}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output.Echo},
		},
	}, output, nil
}
