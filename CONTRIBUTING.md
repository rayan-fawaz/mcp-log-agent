# Contributing to LogMCP

Thank you for your interest in contributing to LogMCP! This document provides guidelines and instructions for contributing.

## Ways to Contribute

### 1. Log Source Adapters

The most impactful contribution is adding support for new log sources. Here are the ones we'd love to have:

| Source | Status | Priority |
|--------|--------|----------|
| Elasticsearch | 🔲 Not started | High |
| AWS CloudWatch | 🔲 Not started | High |
| Datadog | 🔲 Not started | High |
| Splunk | 🔲 Not started | Medium |
| Grafana Loki | 🔲 Not started | Medium |
| Google Cloud Logging | 🔲 Not started | Medium |
| Azure Monitor | 🔲 Not started | Medium |
| Papertrail | 🔲 Not started | Low |
| Loggly | 🔲 Not started | Low |

#### How to Add a Log Source

1. Create a new adapter in `log-server/adapters/`:

```go
// log-server/adapters/elasticsearch.go
package adapters

import (
    "context"
    "github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchAdapter struct {
    client *elasticsearch.Client
    index  string
}

func NewElasticsearchAdapter(url, index string) (*ElasticsearchAdapter, error) {
    cfg := elasticsearch.Config{
        Addresses: []string{url},
    }
    client, err := elasticsearch.NewClient(cfg)
    if err != nil {
        return nil, err
    }
    return &ElasticsearchAdapter{client: client, index: index}, nil
}

func (a *ElasticsearchAdapter) GetLogs(ctx context.Context, start, end int64, region string) ([]LogEntry, error) {
    // Implement Elasticsearch query
    // Use @timestamp field for time range
    // Use region field for filtering
    return nil, nil
}
```

2. Implement the `LogSource` interface:

```go
type LogSource interface {
    GetLogs(ctx context.Context, start, end int64, region string) ([]LogEntry, error)
    GetStats(ctx context.Context) (map[string]int, error)
    Close() error
}
```

3. Add configuration in `log-server/cmd/main.go`:

```go
var logSource data.LogSource

switch os.Getenv("LOG_SOURCE") {
case "elasticsearch":
    logSource = adapters.NewElasticsearchAdapter(
        os.Getenv("ES_URL"),
        os.Getenv("ES_INDEX"),
    )
case "sumologic":
    logSource = adapters.NewSumologicAdapter(
        os.Getenv("SUMO_ACCESS_ID"),
        os.Getenv("SUMO_ACCESS_KEY"),
    )
default:
    logSource = data.NewSQLiteStore() // Default to demo mode
}
```

### 2. Additional MCP Tools

Add new tools to `mcp-server/cmd/main.go`:

```go
// Example: Add a "search_logs" tool for full-text search
server.AddTool(mcp.NewTool("search_logs",
    mcp.WithDescription("Full-text search across all logs"),
    mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
    mcp.WithNumber("limit", mcp.Description("Max results (default 100)")),
), searchLogsHandler)
```

Ideas for new tools:
- `search_logs` - Full-text search
- `get_log_patterns` - Identify recurring patterns
- `compare_timeranges` - Compare two time periods
- `get_error_summary` - Aggregate errors by type
- `stream_logs` - Real-time log streaming

### 3. Documentation

- Add integration guides for specific log sources
- Write tutorials for common use cases
- Improve API documentation
- Add more example prompts

### 4. Bug Fixes

Found a bug? Please:
1. Check if an issue already exists
2. Create a new issue with reproduction steps
3. Submit a PR with the fix

## Development Setup

### Prerequisites

- Go 1.21+
- Git
- VS Code with Agent Mode (for testing)

### Building

```bash
# Clone the repo
git clone https://github.com/rayan-fawaz/mcp-log-agent.git
cd mcp-log-agent

# Build log server
cd log-server/cmd
go build -o log-server

# Build MCP server
cd ../../mcp-server/cmd
go build -o mcp-server
```

### Running Tests

```bash
cd log-server
go test ./...

cd ../mcp-server
go test ./...
```

### Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/elasticsearch-adapter`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit with a clear message (`git commit -m "Add Elasticsearch adapter"`)
6. Push to your fork (`git push origin feature/elasticsearch-adapter`)
7. Open a Pull Request

### PR Requirements

- [ ] Tests pass
- [ ] Code is formatted (`go fmt`)
- [ ] Documentation updated if needed
- [ ] Commit messages are clear
- [ ] No sensitive data or credentials

## Issue Guidelines

When creating an issue, include:

- **Bug reports**: Steps to reproduce, expected vs actual behavior, environment details
- **Feature requests**: Use case, proposed solution, alternatives considered
- **Questions**: Check existing issues first, be specific

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers get started
- Celebrate contributions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Questions? Open an issue or reach out!
