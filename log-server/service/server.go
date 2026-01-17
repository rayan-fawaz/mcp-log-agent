package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/logmcp/log-server/data"
)

// Server handles HTTP requests for logs.
type Server struct {
	store data.Store
}

// New creates a new Server instance.
func New(store data.Store) *Server {
	return &Server{store: store}
}

// GetLogs handles log query requests.
func (s *Server) GetLogs(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if region == "" || startDate == "" || endDate == "" {
		http.Error(w, "Missing required parameters: region, start_date, end_date", http.StatusBadRequest)
		return
	}

	logs, err := s.store.GetLogs(startDate, endDate, region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch logs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"region": region,
		"start":  startDate,
		"end":    endDate,
		"count":  len(logs),
		"logs":   logs,
	})
}

// Health returns server health status.
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Stats returns log statistics.
func (s *Server) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.GetStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ToEpoch converts a date to epoch milliseconds.
func (s *Server) ToEpoch(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	day := r.URL.Query().Get("day")
	timeStr := r.URL.Query().Get("time")

	if year == "" || month == "" || day == "" || timeStr == "" {
		http.Error(w, "Missing parameters: year, month, day, time", http.StatusBadRequest)
		return
	}

	dateStr := fmt.Sprintf("%s-%s-%sT%s:00Z", year, month, day, timeStr)
	t, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid date format: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"epoch_ms": t.UnixMilli(),
		"date":     dateStr,
	})
}

// ToReadable converts epoch milliseconds to a readable date.
func (s *Server) ToReadable(w http.ResponseWriter, r *http.Request) {
	epochMs := r.URL.Query().Get("epoch_ms")
	if epochMs == "" {
		http.Error(w, "Missing parameter: epoch_ms", http.StatusBadRequest)
		return
	}

	ms, err := strconv.ParseInt(epochMs, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid epoch_ms: %v", err), http.StatusBadRequest)
		return
	}

	t := time.UnixMilli(ms)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"epoch_ms": ms,
		"readable": t.Format(time.RFC3339),
		"utc":      t.UTC().String(),
	})
}
