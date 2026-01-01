package internal

//go:generate curl -sSL -o openapi-spec.yaml https://raw.githubusercontent.com/exa-labs/openapi-spec/refs/heads/master/exa-openapi-spec.yaml
//go:generate go run github.com/spachava753/exa-mcp/internal/cmd/fixspec
//go:generate ogen --target exasdk --clean --package exasdk openapi-spec.yaml
