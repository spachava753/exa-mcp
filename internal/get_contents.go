package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spachava753/exa-mcp/internal/exasdk"
)

// GetContentsArgs defines input parameters
type GetContentsArgs struct {
	URLs             []string `json:"urls" jsonschema:"Array of URLs to fetch content from"`
	Livecrawl        string   `json:"livecrawl,omitempty" jsonschema:"Livecrawl mode: never, fallback (default), always, or preferred"`
	MaxTextChars     int      `json:"maxTextChars,omitempty" jsonschema:"Maximum characters for text content"`
	IncludeSummary   bool     `json:"includeSummary,omitempty" jsonschema:"If true, include AI-generated summary"`
	SummaryQuery     string   `json:"summaryQuery,omitempty" jsonschema:"Custom query for summary generation"`
}

// GetContentsOutput is the structured output
type GetContentsOutput struct {
	Results []ContentResult `json:"results"`
	Context string          `json:"context,omitempty"`
}

type ContentResult struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Text          string   `json:"text,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	Author        string   `json:"author,omitempty"`
	PublishedDate string   `json:"publishedDate,omitempty"`
	Highlights    []string `json:"highlights,omitempty"`
}

var GetContentsToolDef = mcp.Tool{
	Name:        "exa_get_contents",
	Description: "Fetch and extract content from specific URLs. Returns clean text, optional summaries, and metadata. Supports live crawling for fresh content.",
	Annotations: &mcp.ToolAnnotations{
		Title: "Exa Get Contents",
	},
}

func GetContentsTool(ctx context.Context, req *mcp.CallToolRequest, args GetContentsArgs) (*mcp.CallToolResult, GetContentsOutput, error) {
	client, err := NewExaClient()
	if err != nil {
		return nil, GetContentsOutput{}, fmt.Errorf("create client: %w", err)
	}

	contentsReq := &exasdk.GetContentsReq{
		Urls: args.URLs,
	}

	// Set text options
	if args.MaxTextChars > 0 {
		var textOpt exasdk.GetContentsReqText1
		textOpt.MaxCharacters.SetTo(args.MaxTextChars)
		contentsReq.Text.SetTo(exasdk.NewGetContentsReqText1GetContentsReqText(textOpt))
	} else {
		contentsReq.Text.SetTo(exasdk.NewBoolGetContentsReqText(true))
	}

	if args.Livecrawl != "" {
		contentsReq.Livecrawl.SetTo(exasdk.GetContentsReqLivecrawl(args.Livecrawl))
	}

	if args.IncludeSummary {
		var summary exasdk.GetContentsReqSummary
		if args.SummaryQuery != "" {
			summary.Query.SetTo(args.SummaryQuery)
		}
		contentsReq.Summary.SetTo(summary)
	}

	resp, err := client.GetContents(ctx, contentsReq)
	if err != nil {
		return nil, GetContentsOutput{}, fmt.Errorf("get contents: %w", err)
	}

	output := GetContentsOutput{}

	for _, r := range resp.Results {
		result := ContentResult{
			Title: r.Title.Value,
			URL:   r.URL.Value.String(),
		}
		if r.Text.Set {
			result.Text = r.Text.Value
		}
		if r.Summary.Set {
			result.Summary = r.Summary.Value
		}
		if r.Author.Set {
			result.Author = r.Author.Value
		}
		if r.PublishedDate.Set {
			result.PublishedDate = r.PublishedDate.Value
		}
		result.Highlights = r.Highlights
		output.Results = append(output.Results, result)
	}

	if resp.Context.Set {
		output.Context = resp.Context.Value
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonBytes)},
		},
	}, output, nil
}
