package internal

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestEchoTool(t *testing.T) {
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	server := GetServer("test")
	_, err := server.Connect(t.Context(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(t.Context(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}

	result, err := session.CallTool(t.Context(), &mcp.CallToolParams{
		Name: "echo",
		Arguments: map[string]any{
			"message": "hello world",
		},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}

	if result.IsError {
		t.Fatalf("unexpected error result")
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content, got %d", len(result.Content))
	}

	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}

	if tc.Text != "hello world" {
		t.Errorf("expected 'hello world', got %q", tc.Text)
	}
}
