package internal

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetServer creates and configures the MCP server with Exa AI tools.
func GetServer(version string) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "exa-mcp-server",
		Title:   "Exa MCP Server",
		Version: version,
	}, nil)

	// Register Exa tools
	mcp.AddTool(server, &SearchToolDef, SearchTool)
	mcp.AddTool(server, &FindSimilarToolDef, FindSimilarTool)
	mcp.AddTool(server, &GetContentsToolDef, GetContentsTool)
	mcp.AddTool(server, &AnswerToolDef, AnswerTool)

	return server
}
