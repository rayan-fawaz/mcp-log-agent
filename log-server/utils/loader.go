package utils

import (
	"encoding/json"
	"os"
)

// Log represents a parsed log entry.
type Log struct {
	Raw RawMessage `json:"raw"`
}

// RawMessage contains the raw log data.
type RawMessage struct {
	Timestamp int64  `json:"time"`
	Log       string `json:"log"`
}

// LoadFile reads a file from disk.
func LoadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Parse parses JSON log data.
func Parse(data []byte) ([]Log, error) {
	var logs []Log
	if err := json.Unmarshal(data, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
