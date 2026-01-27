package main

import (
	"fmt"
	"goops-monitor/runner"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var metrics *Metrics

func main() {
	// Load configuration
	cfg, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize metrics
	metrics = NewMetrics()

	fmt.Printf("Starting monitor for %s with interval %v\n", cfg.Target, cfg.Interval)

	// Start API server
	go StartAPI(":8081")

	// Start a local test server for demonstration
	go startServer(":8080")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkStatus(cfg)
		case <-sigChan:
			fmt.Println("\nShutting down gracefully...")
			stats := metrics.GetStats()
			fmt.Printf("Final stats - Ups: %d, Downs: %d\n", stats["ups"], stats["downs"])
			return
		}
	}
}

func checkStatus(cfg *Config) {
	resp, err := http.Get(cfg.Target + "/health")
	if err != nil {
		log.Printf("DOWN: %s (%v)", cfg.Target, err)
		metrics.RecordDown()
		executeRemediation(cfg.OnFailure)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("UP: %s (Status: %d)", cfg.Target, resp.StatusCode)
		metrics.RecordUp()
	} else {
		log.Printf("DOWN: %s (Status: %d)", cfg.Target, resp.StatusCode)
		metrics.RecordDown()
		executeRemediation(cfg.OnFailure)
	}
}

func executeRemediation(script string) {
	if script == "" {
		return
	}
	
	log.Printf("Executing remediation: %s", script)
	r := runner.New("sh", "-c", script)
	if err := r.Execute(); err != nil {
		log.Printf("Remediation failed: %v", err)
	}
}
