package data

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/logmcp/log-server/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTableSQL = `CREATE TABLE IF NOT EXISTS logs (
		region  TEXT NOT NULL,
		time    INTEGER NOT NULL,
		message TEXT NOT NULL
	)`
	countLogsSQL  = `SELECT COUNT(*) FROM logs`
	insertLogSQL  = `INSERT INTO logs (region, time, message) VALUES (?, ?, ?)`
	selectLogsSQL = `SELECT region, time, message FROM logs WHERE region = ? AND time BETWEEN ? AND ?`
	statsSQL      = `SELECT region, COUNT(*) as count FROM logs GROUP BY region`
)

// Store defines the interface for log storage operations.
type Store interface {
	GetLogs(start, end, region string) ([]LogEntry, error)
	GetStats() (map[string]int, error)
}

// LogEntry represents a single log record.
type LogEntry struct {
	Region  string `json:"region"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

// Database implements Store using SQLite.
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new SQLite database connection and loads demo data if needed.
func NewDatabase() (*Database, error) {
	conn, err := sql.Open("sqlite3", "./logs.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := conn.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	var count int
	if err := conn.QueryRow(countLogsSQL).Scan(&count); err != nil {
		return nil, fmt.Errorf("failed to count logs: %w", err)
	}

	if count == 0 {
		if err := loadDemoData(conn); err != nil {
			return nil, fmt.Errorf("failed to load demo data: %w", err)
		}
	}

	return &Database{conn: conn}, nil
}

// GetLogs retrieves logs for a region within a time range.
func (db *Database) GetLogs(start, end, region string) ([]LogEntry, error) {
	rows, err := db.conn.Query(selectLogsSQL, region, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var entry LogEntry
		if err := rows.Scan(&entry.Region, &entry.Time, &entry.Message); err != nil {
			continue
		}
		logs = append(logs, entry)
	}
	return logs, nil
}

// GetStats returns log counts per region.
func (db *Database) GetStats() (map[string]int, error) {
	rows, err := db.conn.Query(statsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var region string
		var count int
		if err := rows.Scan(&region, &count); err != nil {
			continue
		}
		stats[region] = count
	}
	return stats, nil
}

// loadDemoData populates the database with sample logs.
func loadDemoData(conn *sql.DB) error {
	demoPath := os.Getenv("DEMO_LOGS_PATH")
	if demoPath == "" {
		demoPath = "../demo_logs"
	}

	regions := map[string]string{
		"NA": demoPath + "/sample_logs_na.json",
		"EU": demoPath + "/sample_logs_eu.json",
		"AP": demoPath + "/sample_logs_ap.json",
	}

	fmt.Printf("Loading demo logs from %s\n", demoPath)

	for region, path := range regions {
		data, err := utils.LoadFile(path)
		if err != nil {
			fmt.Printf("Warning: could not load %s logs: %v\n", region, err)
			continue
		}

		logs, err := utils.Parse(data)
		if err != nil {
			fmt.Printf("Warning: could not parse %s logs: %v\n", region, err)
			continue
		}

		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(insertLogSQL)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, log := range logs {
			stmt.Exec(region, log.Raw.Timestamp, log.Raw.Log)
		}
		stmt.Close()

		if err := tx.Commit(); err != nil {
			return err
		}

		fmt.Printf("Loaded %d logs for %s region\n", len(logs), region)
	}

	return nil
}
