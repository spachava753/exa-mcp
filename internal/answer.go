package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spachava753/exa-mcp/internal/exasdk"
)

// AnswerArgs defines input parameters
type AnswerArgs struct {
	Query       string `json:"query" jsonschema:"required,description=The question or query to answer"`
	IncludeText bool   `json:"includeText,omitempty" jsonschema:"description=If true, include full text content from sources"`
}

// AnswerOutput is the structured output
type AnswerOutput struct {
	Answer    string           `json:"answer"`
	Citations []AnswerCitation `json:"citations"`
}

type AnswerCitation struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	PublishedDate string `json:"publishedDate,omitempty"`
	Author        string `json:"author,omitempty"`
	Text          string `json:"text,omitempty"`
}

var AnswerToolDef = mcp.Tool{
	Name:        "exa_answer",
	Description: "Get an AI-generated answer to a question based on web search results. Returns a concise answer with citations to source documents.",
	Annotations: &mcp.ToolAnnotations{
		Title: "Exa Answer",
	},
}

func AnswerTool(ctx context.Context, req *mcp.CallToolRequest, args AnswerArgs) (*mcp.CallToolResult, AnswerOutput, error) {
	client, err := NewExaClient()
	if err != nil {
		return nil, AnswerOutput{}, fmt.Errorf("create client: %w", err)
	}

	answerReq := &exasdk.AnswerReq{
		Query: args.Query,
	}

	if args.IncludeText {
		answerReq.Text.SetTo(true)
	}

	resp, err := client.Answer(ctx, answerReq)
	if err != nil {
		return nil, AnswerOutput{}, fmt.Errorf("answer: %w", err)
	}

	output := AnswerOutput{
		Answer: resp.Answer.Value,
	}

	for _, c := range resp.Citations {
		citation := AnswerCitation{
			Title: c.Title.Value,
			URL:   c.URL.Value.String(),
		}
		if c.PublishedDate.Set {
			citation.PublishedDate = c.PublishedDate.Value
		}
		if c.Author.Set {
			citation.Author = c.Author.Value
		}
		if c.Text.Set {
			citation.Text = c.Text.Value
		}
		output.Citations = append(output.Citations, citation)
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonBytes)},
		},
	}, output, nil
}
