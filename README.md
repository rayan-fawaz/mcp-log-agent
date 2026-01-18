# LogMCP

> **An MCP server that gives AI agents the power to query, analyze, and act on your application logs.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![MCP](https://img.shields.io/badge/MCP-Compatible-blue?style=flat)](https://modelcontextprotocol.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

LogMCP connects AI assistants (GitHub Copilot, Claude, etc.) to your log data through the [Model Context Protocol](https://modelcontextprotocol.io). Ask questions in natural language, get instant answers, and automate incident response workflows.

![Demo](https://via.placeholder.com/800x400?text=VS+Code+Agent+Mode+Demo)

---

## ✨ What Can It Do?

| Capability | Example |
|------------|---------|
| **Query logs** | *"Show me all errors from the EU region in the last hour"* |
| **Analyze patterns** | *"What's causing the spike in 401 authentication failures?"* |
| **Debug issues** | *"Find slow requests over 500ms and identify the root cause"* |
| **Correlate events** | *"Were there any deployments before this error started?"* |
| **Automate actions** | *"Summarize these errors and create a Jira ticket"* |
| **Support teams** | *"Explain this 503 error in simple terms for the customer"* |

---

## 🎯 Why LogMCP?

**The Problem:**
- Engineers spend hours searching through millions of log lines
- Distributed microservices make correlation difficult
- Support teams can't answer customer questions without engineering help
- Context switching between tools (logs, Jira, GitHub) slows incident response

**The Solution:**
- AI agents query logs using natural language
- Automatic pattern detection and summarization
- One interface to logs + Jira + GitHub (via official MCP servers)
- Support teams get instant answers without escalating

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    VS Code Agent Mode                        │
│              (GitHub Copilot, Claude, GPT, etc.)            │
└─────────────────────────────┬───────────────────────────────┘
                              │ MCP Protocol
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       MCP Servers                            │
├─────────────────┬──────────────────┬────────────────────────┤
│    LogMCP       │   GitHub MCP     │    Atlassian MCP       │
│  (this repo)    │   (official)     │     (official)         │
├─────────────────┼──────────────────┼────────────────────────┤
│ • get_logs      │ • create_issue   │ • create_jira_issue    │
│ • to_epoch      │ • create_pr      │ • search_issues        │
│ • to_readable   │ • search_code    │ • update_issue         │
│ • send_teams    │ • list_commits   │ • add_comment          │
└────────┬────────┴──────────────────┴────────────────────────┘
         │
         │ HTTP/REST
         ▼
┌─────────────────┐      ┌─────────────────────────────────────┐
│   Log Server    │◄─────│  Your Log Source                    │
│    (SQLite)     │      │  • Demo logs (included)             │
└─────────────────┘      │  • Sumologic, Datadog, Splunk, etc. │
                         │  • CloudWatch, Elasticsearch        │
                         │  • Any JSON log source              │
                         └─────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites
- [Go 1.21+](https://golang.org/dl/)
- [VS Code](https://code.visualstudio.com/) with Agent Mode enabled

### 1. Clone the repository

```bash
git clone https://github.com/rayan-fawaz/mcp-log-agent.git
cd mcp-log-agent
```

### 2. Start the Log Server

```bash
cd log-server/cmd
go run main.go
```

You'll see:
```
Loading demo logs from ../demo_logs
Loaded 10 logs for NA region
Loaded 10 logs for EU region
Loaded 10 logs for AP region
Log server starting on port :8081
```

### 3. Start the MCP Server

Open a new terminal:

```bash
cd mcp-server/cmd
go run main.go
```

You'll see:
```
LogMCP server running on port 8080
```

### 4. Configure VS Code

Add to your VS Code settings (`Ctrl+Shift+P` → "Open User Settings JSON"):

```json
{
  "mcp": {
    "servers": {
      "logmcp": {
        "type": "http",
        "url": "http://localhost:8080/mcp"
      }
    }
  }
}
```

### 5. Try It!

Open VS Code Agent Mode (Copilot Chat) and ask:

> Get logs from the NA region between 1768473996370 and 1768473996400

> What errors are happening in the EU region? Summarize the patterns.

> Convert January 18, 2026 at 10:46:36 to epoch milliseconds

---

## 🔌 Connecting Your Own Logs

LogMCP ships with demo logs, but you can easily connect your own log sources.

### Option 1: Replace Demo Logs (Simplest)

Export your logs as JSON and place them in `demo_logs/`:

```json
[
  {
    "raw": {
      "time": 1768473996379,
      "log": "lvl=error reqID=abc-123 status=500 msg=\"Database timeout\""
    }
  }
]
```

### Option 2: Connect to Sumologic

Modify `log-server/data/sqlite.go` to fetch from Sumologic API:

```go
// In loadDemoData() or create a new function:
func fetchFromSumologic(query string, start, end int64) ([]LogEntry, error) {
    // Sumologic Search Job API
    client := &http.Client{}
    req, _ := http.NewRequest("POST", 
        "https://api.sumologic.com/api/v1/search/jobs",
        strings.NewReader(fmt.Sprintf(`{"query":"%s","from":%d,"to":%d}`, query, start, end)))
    req.SetBasicAuth(os.Getenv("SUMOLOGIC_ACCESS_ID"), os.Getenv("SUMOLOGIC_ACCESS_KEY"))
    // ... handle response
}
```

### Option 3: Connect to Other Sources

The log server is a simple REST API. You can:

| Source | Integration Approach |
|--------|---------------------|
| **Elasticsearch** | Replace SQLite queries with ES client |
| **CloudWatch** | Use AWS SDK to fetch log groups |
| **Datadog** | Use Datadog Logs API |
| **Splunk** | Use Splunk REST API |
| **File-based** | Read from log files directly |

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed integration guides.

### Option 4: Build a Custom Adapter

Create a new log source adapter:

```go
// log-server/adapters/your_source.go
type YourSourceAdapter struct {
    apiKey string
}

func (a *YourSourceAdapter) GetLogs(start, end int64, region string) ([]LogEntry, error) {
    // Implement your log fetching logic
}
```

---

## 🛠️ MCP Tools Reference

| Tool | Description | Parameters |
|------|-------------|------------|
| `get_logs` | Query logs by region and time range | `region` (NA/EU/AP), `start` (epoch ms), `end` (epoch ms) |
| `to_epoch` | Convert human date to epoch milliseconds | `year`, `month`, `day`, `time` (HH:MM:SS) |
| `to_readable` | Convert epoch to readable date | `epoch_ms` |
| `send_teams` | Send message to Microsoft Teams | `message` |

### Example Tool Responses

**get_logs:**
```json
{
  "region": "NA",
  "count": 5,
  "logs": [
    {"time": 1768473996379, "message": "lvl=error status=500 msg=\"DB timeout\""},
    {"time": 1768473996380, "message": "lvl=info status=200 msg=\"Request completed\""}
  ]
}
```

---

## 🔗 Adding GitHub & Jira Integration

LogMCP works alongside official MCP servers. Update your VS Code settings:

```json
{
  "mcp": {
    "servers": {
      "logmcp": {
        "type": "http",
        "url": "http://localhost:8080/mcp"
      },
      "github": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-github"],
        "env": {
          "GITHUB_PERSONAL_ACCESS_TOKEN": "<your-token>"
        }
      },
      "atlassian": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-atlassian"],
        "env": {
          "ATLASSIAN_API_TOKEN": "<your-token>",
          "ATLASSIAN_EMAIL": "you@company.com",
          "ATLASSIAN_DOMAIN": "yourcompany.atlassian.net"
        }
      }
    }
  }
}
```

Now you can ask:
> Find all 500 errors in the last hour and create a GitHub issue with the summary

---

## 📁 Project Structure

```
mcp-log-agent/
├── log-server/              # REST API for log storage & queries
│   ├── cmd/main.go          # Entry point
│   ├── data/sqlite.go       # Database layer (swap for your source)
│   ├── service/server.go    # HTTP handlers
│   └── utils/loader.go      # Log parsing utilities
├── mcp-server/              # MCP protocol server
│   ├── cmd/main.go          # MCP tools & handlers
│   ├── mcp.json             # VS Code config
│   └── mcp.example.json     # Config with GitHub/Jira
├── demo_logs/               # Sample log data
│   ├── sample_logs_na.json  # North America region
│   ├── sample_logs_eu.json  # Europe region
│   └── sample_logs_ap.json  # Asia Pacific region
└── demo/
    └── prompts.md           # Example prompts to try
```

---

## ⚙️ Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `LOG_SERVER_PORT` | 8081 | Port for the log REST API |
| `MCP_SERVER_PORT` | 8080 | Port for the MCP server |
| `LOG_SERVER_URL` | http://localhost:8081 | URL for MCP to reach log server |
| `DEMO_LOGS_PATH` | ../demo_logs | Path to demo log files |
| `TEAMS_WEBHOOK_URL` | - | Microsoft Teams webhook URL |

Copy `.env.example` to `.env` and configure as needed.

---

## 🧪 API Reference

The log server exposes these REST endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/logs?region=NA&start_date=X&end_date=Y` | Query logs |
| GET | `/health` | Health check |
| GET | `/stats` | Log counts by region |
| GET | `/time/epoch?year=2026&month=1&day=18&time=10:30:00` | Convert to epoch |
| GET | `/time/readable?epoch_ms=1768473996379` | Convert to readable |

---

## 🤝 Contributing

Contributions are welcome! Areas where help is needed:

- [ ] **Log source adapters** - Elasticsearch, CloudWatch, Datadog, Splunk
- [ ] **Additional MCP tools** - Log streaming, alerts, metrics
- [ ] **UI dashboard** - Web interface for log visualization
- [ ] **Documentation** - More examples and tutorials

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 📜 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io) by Anthropic
- [mcp-go](https://github.com/mark3labs/mcp-go) - Go SDK for MCP
- Built for the VS Code Agent Mode ecosystem

---

**⭐ Star this repo if you find it useful!**
