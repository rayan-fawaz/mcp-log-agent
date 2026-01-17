# LogMCP

An AI-powered log analysis agent that reads logs across microservices, summarizes issues, answers questions, and proactively alerts engineers before customers are impacted. Built for VS Code Agent Mode using the Model Context Protocol (MCP).

## The Problem

- **Massive log volume** from 30+ microservices makes manual analysis impractical
- **Difficulty identifying customer issues** in distributed, correlated logs
- **Time-consuming debugging** requiring engineers to search across multiple sources
- **Support teams struggle** to answer questions without engineering help
- **Lack of context** when correlating logs, metrics, and code changes

## The Solution

LogMCP connects AI agents to your observability stack, enabling:

- **AI-powered chat interface** for support teams to get instant answers
- **Jira integration** - automatic ticket commenting with log insights and suggested actions
- **GitHub integration** - links code changes, issues, and PRs to observed log events
- **Proactive alerting** - surfaces critical issues before customers are impacted
- **Accelerated onboarding** - contextual summaries help new engineers ramp up fast

## Example Queries

Ask the AI agent questions like:

- "Show me all errors from the EU region in the last hour"
- "What's causing the 401 authentication failures?"
- "Summarize the 500 errors and create a Jira ticket for the ops team"
- "Find slow requests over 500ms and identify the root cause"
- "Correlate this error with recent code deployments"

## Architecture

```
+----------------------------------------------------------+
|                    VS Code Agent Mode                     |
|                  (Copilot, Claude, etc.)                  |
+----------------------------+-----------------------------+
                             | MCP Protocol
                             v
+----------------------------------------------------------+
|                      MCP Servers                          |
+-------------+-----------------+--------------------------+
|   LogMCP    |   GitHub MCP    |    Atlassian MCP         |
| (this repo) |   (official)    |     (official)           |
+-------------+-----------------+--------------------------+
| get_logs    | create_issue    | create_jira_issue        |
| to_epoch    | create_pr       | search_issues            |
| to_readable | search_code     | update_issue             |
| send_teams  |                 |                          |
+------+------+-----------------+--------------------------+
       |
       | HTTP
       v
+--------------+
|  Log Server  |
|   (SQLite)   |
+--------------+
```

## Quick Start

### 1. Clone and build

```bash
git clone https://github.com/YOUR_USERNAME/logmcp.git
cd logmcp
```

### 2. Start the log server

```bash
cd log-server/cmd
go run main.go
```

Server starts with demo logs automatically.

### 3. Start the MCP server

```bash
cd mcp-server/cmd
go run main.go
```

### 4. Configure VS Code

Add to your VS Code settings (Cmd/Ctrl+Shift+P -> Open User Settings JSON):

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

### 5. Try it

Open VS Code Agent Mode and ask:

> Get logs from the NA region between 1768470000000 and 1768475000000

## Tools

| Tool | Description |
|------|-------------|
| `get_logs` | Query logs by region (NA/EU/AP) and time range |
| `to_epoch` | Convert date to epoch milliseconds |
| `to_readable` | Convert epoch to readable date |
| `send_teams` | Send message to Teams (optional) |

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /logs` | Query logs with region, start_date, end_date |
| `GET /health` | Health check |
| `GET /stats` | Log counts per region |
| `GET /time/epoch` | Convert date to epoch |
| `GET /time/readable` | Convert epoch to date |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_SERVER_PORT` | 8081 | Log server port |
| `MCP_SERVER_PORT` | 8080 | MCP server port |
| `LOG_SERVER_URL` | http://localhost:8081 | Log server URL |
| `DEMO_LOGS_PATH` | ../demo_logs | Path to demo logs |
| `TEAMS_WEBHOOK_URL` | - | Teams webhook (optional) |

## Adding GitHub + Jira

See `mcp-server/mcp.example.json` for how to add the official GitHub and Atlassian MCP servers alongside LogMCP.

## Project Structure

```
logmcp/
├── log-server/          # Log storage and API
│   ├── cmd/             # Entry point
│   ├── data/            # Database layer
│   ├── service/         # HTTP handlers
│   └── utils/           # Log parsing
├── mcp-server/          # MCP protocol server
│   └── cmd/             # Entry point
└── demo_logs/           # Sample log data
```

## License

MIT
