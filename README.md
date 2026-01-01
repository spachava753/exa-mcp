# Exa MCP Server

An MCP (Model Context Protocol) server that exposes Exa AI's search capabilities as tools.

## Tools

This server provides the following tools:

### exa_search
Perform web searches using Exa's AI-powered search engine. Supports:
- **neural**: Embeddings-based semantic search
- **fast**: Streamlined neural search
- **auto**: Intelligently combines search methods (default)
- **deep**: Comprehensive search with query expansion

### exa_find_similar
Find web pages similar to a given URL. Useful for discovering related content, competitor analysis, or finding more resources on a topic.

### exa_get_contents
Fetch and extract content from specific URLs. Returns clean text, optional summaries, and metadata. Supports live crawling for fresh content.

### exa_answer
Get an AI-generated answer to a question based on web search results. Returns a concise answer with citations to source documents.

## Configuration

Set the `EXA_API_KEY` environment variable with your Exa API key.

```bash
export EXA_API_KEY="your-api-key-here"
```

## Usage

### As an MCP Server

Run the server:

```bash
./exa-mcp-server
```

Or with go run:

```bash
go run . 
```

### Version

```bash
./exa-mcp-server version
```

## Development

### Regenerating the SDK

The SDK is generated from the OpenAPI spec using [ogen](https://github.com/ogen-go/ogen):

```bash
ogen --target exasdk --clean exa-openapi-spec.yaml
```

### Building

```bash
go build .
```

## License

See LICENSE file.
