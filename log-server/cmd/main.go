package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/logmcp/log-server/data"
	"github.com/logmcp/log-server/service"
)

func main() {
	port := getPort()

	// Initialize the database
	db, err := data.NewDatabase()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Create the HTTP server
	srv := service.New(db)

	// Register routes
	http.HandleFunc("/logs", srv.GetLogs)
	http.HandleFunc("/health", srv.Health)
	http.HandleFunc("/stats", srv.Stats)
	http.HandleFunc("/time/epoch", srv.ToEpoch)
	http.HandleFunc("/time/readable", srv.ToReadable)

	fmt.Printf("Log server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getPort() string {
	if port := os.Getenv("LOG_SERVER_PORT"); port != "" {
		return ":" + port
	}
	return ":8081"
}
