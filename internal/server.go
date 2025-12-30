package internal

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetServer creates and configures the MCP server with all tools.
// Add new tools by calling mcp.AddTool() with the tool definition and handler.
func GetServer(version string) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "my-mcp-server",  // TODO: Replace with your server name
		Title:   "My MCP Server",  // TODO: Replace with your server title
		Version: version,
	}, nil)

	// Register tools here. Each tool needs:
	// 1. A Tool definition (name, description, annotations)
	// 2. A handler function matching the signature:
	//    func(ctx, *mcp.CallToolRequest, ArgsType) (*mcp.CallToolResult, OutputType, error)
	mcp.AddTool(server, &EchoToolDef, EchoTool)

	// Example: Add more tools like this:
	// mcp.AddTool(server, &MyToolDef, MyToolHandler)

	return server
}
