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
	http.HandleFunc("/ready", handleReady)
	http.HandleFunc("/live", handleLive)
	
	LogInfo("Starting API server", "port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		LogError("API server failed", "error", err.Error())
	}
}

// handleMetrics uses Prometheus handler
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// handleReady returns 200 if the monitor has collected at least one check
func handleReady(w http.ResponseWriter, r *http.Request) {
	stats := metrics.GetStats()
	if len(stats) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

// handleLive always returns 200 (process is alive)
func handleLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alive"))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stats := metrics.GetStats()
	allHealthy := true
	
	targetStats := make(map[string]interface{})
	for target, s := range stats {
		total := s.Ups + s.Downs
		uptimePct := 0.0
		if total > 0 {
			uptimePct = float64(s.Ups) / float64(total) * 100
		}
		
		if s.Downs > 0 {
			allHealthy = false
		}
		
		targetStats[target] = map[string]interface{}{
			"ups":        s.Ups,
			"downs":      s.Downs,
			"uptime_pct": uptimePct,
			"latency_ms": s.LastLatency,
		}
	}
	
	status := map[string]interface{}{
		"healthy": allHealthy,
		"targets": targetStats,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
