package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartAPI starts the HTTP API server
func StartAPI(port string) {
	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/status", handleStatus)
	
	LogInfo("Starting API server", "port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		LogError("API server failed", "error", err.Error())
	}
}

// handleMetrics now uses Prometheus handler
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stats := metrics.GetStats()
	allHealthy := true
	for _, s := range stats {
		if s.Downs > 0 {
			allHealthy = false
			break
		}
	}
	
	status := map[string]interface{}{
		"healthy": allHealthy,
		"metrics": stats,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
