package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spachava753/exa-mcp/internal/exasdk"
)

// FindSimilarArgs defines input parameters
type FindSimilarArgs struct {
	URL            string   `json:"url" jsonschema:"The URL for which to find similar links"`
	NumResults     int      `json:"numResults,omitempty" jsonschema:"Number of results to return (default 10)"`
	IncludeDomains []string `json:"includeDomains,omitempty" jsonschema:"List of domains to include"`
	ExcludeDomains []string `json:"excludeDomains,omitempty" jsonschema:"List of domains to exclude"`
	GetContents    bool     `json:"getContents,omitempty" jsonschema:"If true, return page contents"`
}

// FindSimilarOutput is the structured output
type FindSimilarOutput struct {
	Results []SearchResult `json:"results"`
	Context string         `json:"context,omitempty"`
}

var FindSimilarToolDef = mcp.Tool{
	Name:        "exa_find_similar",
	Description: "Find web pages similar to a given URL. Useful for discovering related content, competitor analysis, or finding more resources on a topic.",
	Annotations: &mcp.ToolAnnotations{
		Title: "Exa Find Similar",
	},
}

func FindSimilarTool(ctx context.Context, req *mcp.CallToolRequest, args FindSimilarArgs) (*mcp.CallToolResult, FindSimilarOutput, error) {
	client, err := NewExaClient()
	if err != nil {
		return nil, FindSimilarOutput{}, fmt.Errorf("create client: %w", err)
	}

	findReq := &exasdk.FindSimilarReq{
		URL:            args.URL,
		IncludeDomains: args.IncludeDomains,
		ExcludeDomains: args.ExcludeDomains,
	}

	if args.NumResults > 0 {
		findReq.NumResults.SetTo(args.NumResults)
	}

	if args.GetContents {
		var contentsReq exasdk.ContentsRequest
		contentsReq.Text.SetTo(exasdk.NewBoolContentsRequestText(true))
		findReq.Contents.SetTo(contentsReq)
	}

	resp, err := client.FindSimilar(ctx, findReq)
	if err != nil {
		return nil, FindSimilarOutput{}, fmt.Errorf("find similar: %w", err)
	}

	output := FindSimilarOutput{}

	for _, r := range resp.Results {
		result := SearchResult{
			Title: r.Title.Value,
			URL:   r.URL.Value.String(),
		}
		if r.PublishedDate.Set {
			result.PublishedDate = r.PublishedDate.Value
		}
		if r.Author.Set {
			result.Author = r.Author.Value
		}
		if r.Score.Set {
			result.Score = r.Score.Value
		}
		if r.Text.Set {
			result.Text = r.Text.Value
		}
		if r.Summary.Set {
			result.Summary = r.Summary.Value
		}
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
