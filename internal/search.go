package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spachava753/exa-mcp/internal/exasdk"
)

// SearchArgs defines input parameters for the search tool
type SearchArgs struct {
	Query          string   `json:"query" jsonschema:"The query string for the search"`
	Type           string   `json:"type,omitempty" jsonschema:"Search type: neural (embeddings-based), fast (streamlined), auto (default - intelligently combines methods), deep (comprehensive with query expansion)"`
	Category       string   `json:"category,omitempty" jsonschema:"A data category to focus on"`
	NumResults     int      `json:"numResults,omitempty" jsonschema:"Number of results to return (default 10)"`
	IncludeDomains []string `json:"includeDomains,omitempty" jsonschema:"List of domains to include in the search"`
	ExcludeDomains []string `json:"excludeDomains,omitempty" jsonschema:"List of domains to exclude from search results"`
	IncludeText    []string `json:"includeText,omitempty" jsonschema:"Strings that must be present in webpage text (max 1 string, up to 5 words)"`
	ExcludeText    []string `json:"excludeText,omitempty" jsonschema:"Strings that must not be present in webpage text"`
	GetContents    bool     `json:"getContents,omitempty" jsonschema:"If true, return page contents along with search results"`
}

// SearchOutput is the structured output
type SearchOutput struct {
	Results []SearchResult `json:"results"`
	Context string         `json:"context,omitempty"`
}

type SearchResult struct {
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	PublishedDate string  `json:"publishedDate,omitempty"`
	Author        string  `json:"author,omitempty"`
	Score         float64 `json:"score,omitempty"`
	Text          string  `json:"text,omitempty"`
	Summary       string  `json:"summary,omitempty"`
}

var SearchToolDef = mcp.Tool{
	Name:        "exa_search",
	Description: "Perform a web search using Exa's AI-powered search engine. Returns relevant results with optional page contents. Supports neural (semantic), fast, auto, and deep search types.",
	Annotations: &mcp.ToolAnnotations{
		Title: "Exa Search",
	},
}

func SearchTool(ctx context.Context, req *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, SearchOutput, error) {
	client, err := NewExaClient()
	if err != nil {
		return nil, SearchOutput{}, fmt.Errorf("create client: %w", err)
	}

	searchReq := &exasdk.SearchReq{
		Query:          args.Query,
		IncludeDomains: args.IncludeDomains,
		ExcludeDomains: args.ExcludeDomains,
		IncludeText:    args.IncludeText,
		ExcludeText:    args.ExcludeText,
	}

	if args.NumResults > 0 {
		searchReq.NumResults.SetTo(args.NumResults)
	}

	if args.Type != "" {
		searchReq.Type.SetTo(exasdk.SearchReqType(args.Type))
	}

	if args.Category != "" {
		searchReq.Category.SetTo(exasdk.SearchReqCategory(args.Category))
	}

	if args.GetContents {
		var contentsReq exasdk.ContentsRequest
		contentsReq.Text.SetTo(exasdk.NewBoolContentsRequestText(true))
		searchReq.Contents.SetTo(contentsReq)
	}

	resp, err := client.Search(ctx, searchReq)
	if err != nil {
		return nil, SearchOutput{}, fmt.Errorf("search: %w", err)
	}

	output := SearchOutput{}
	
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
