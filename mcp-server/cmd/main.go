package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"LogMCP",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Register tools
	s.AddTool(logsTool(), handleGetLogs)
	s.AddTool(toEpochTool(), handleToEpoch)
	s.AddTool(toReadableTool(), handleToReadable)
	s.AddTool(teamsTool(), handleTeams)

	// Start server
	port := getEnv("MCP_SERVER_PORT", "8080")
	fmt.Printf("LogMCP server running on port %s\n", port)

	srv := server.NewStreamableHTTPServer(s, []server.StreamableHTTPOption{}...)
	if err := srv.Start(":" + port); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

// Tool definitions

func logsTool() mcp.Tool {
	return mcp.NewTool("get_logs",
		mcp.WithDescription("Query application logs by region and time range"),
		mcp.WithString("region", mcp.Required(), mcp.Description("Region: NA, EU, or AP")),
		mcp.WithString("start", mcp.Required(), mcp.Description("Start time (epoch milliseconds)")),
		mcp.WithString("end", mcp.Required(), mcp.Description("End time (epoch milliseconds)")),
	)
}

func toEpochTool() mcp.Tool {
	return mcp.NewTool("to_epoch",
		mcp.WithDescription("Convert date to epoch milliseconds"),
		mcp.WithString("year", mcp.Required(), mcp.Description("Year (YYYY)")),
		mcp.WithString("month", mcp.Required(), mcp.Description("Month (1-12)")),
		mcp.WithString("day", mcp.Required(), mcp.Description("Day (1-31)")),
		mcp.WithString("time", mcp.Required(), mcp.Description("Time (HH:MM:SS)")),
	)
}

func toReadableTool() mcp.Tool {
	return mcp.NewTool("to_readable",
		mcp.WithDescription("Convert epoch milliseconds to readable date"),
		mcp.WithString("epoch_ms", mcp.Required(), mcp.Description("Epoch milliseconds")),
	)
}

func teamsTool() mcp.Tool {
	return mcp.NewTool("send_teams",
		mcp.WithDescription("Send message to Microsoft Teams"),
		mcp.WithString("message", mcp.Required(), mcp.Description("Message to send")),
	)
}

// Tool handlers

func handleGetLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	region, _ := req.RequireString("region")
	start, _ := req.RequireString("start")
	end, _ := req.RequireString("end")

	url := fmt.Sprintf("%s/logs?region=%s&start_date=%s&end_date=%s",
		getEnv("LOG_SERVER_URL", "http://localhost:8081"), region, start, end)

	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Log server unavailable: %v", err)), nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return mcp.NewToolResultText(string(body)), nil
}

func handleToEpoch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	year, _ := req.RequireString("year")
	month, _ := req.RequireString("month")
	day, _ := req.RequireString("day")
	timeStr, _ := req.RequireString("time")

	url := fmt.Sprintf("%s/time/epoch?year=%s&month=%s&day=%s&time=%s",
		getEnv("LOG_SERVER_URL", "http://localhost:8081"), year, month, day, timeStr)

	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return mcp.NewToolResultText(string(body)), nil
}

func handleToReadable(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	epochMs, _ := req.RequireString("epoch_ms")

	url := fmt.Sprintf("%s/time/readable?epoch_ms=%s",
		getEnv("LOG_SERVER_URL", "http://localhost:8081"), epochMs)

	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return mcp.NewToolResultText(string(body)), nil
}

func handleTeams(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	webhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
	if webhookURL == "" {
		return mcp.NewToolResultError("Teams not configured (set TEAMS_WEBHOOK_URL)"), nil
	}

	message, _ := req.RequireString("message")
	payload, _ := json.Marshal(map[string]string{"text": message})

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send: %v", err)), nil
	}
	defer resp.Body.Close()

	return mcp.NewToolResultText("Message sent to Teams"), nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
