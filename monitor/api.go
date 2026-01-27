package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// StartAPI starts the HTTP API server
func StartAPI(port string) {
	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/status", handleStatus)
	
	log.Printf("Starting API server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("API server failed: %v", err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	stats := metrics.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stats := metrics.GetStats()
	status := map[string]interface{}{
		"healthy": stats["downs"] == 0,
		"metrics": stats,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
