package internal

import (
	"context"
	"os"

	"github.com/spachava753/exa-mcp/internal/exasdk"
)

// apiKeySource implements exasdk.SecuritySource
type apiKeySource struct {
	apiKey string
}

func (s *apiKeySource) Apikey(ctx context.Context, operationName exasdk.OperationName) (exasdk.Apikey, error) {
	return exasdk.Apikey{
		APIKey: s.apiKey,
	}, nil
}

// NewExaClient creates a new Exa API client
func NewExaClient() (*exasdk.Client, error) {
	apiKey := os.Getenv("EXA_API_KEY")
	
	return exasdk.NewClient(
		"https://api.exa.ai",
		&apiKeySource{apiKey: apiKey},
	)
}
