package main

import (
	"fmt"
	"goops-monitor/runner"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

	fmt.Printf("Starting monitor for %d targets with interval %v\n", len(cfg.Targets), cfg.Interval)

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
			var wg sync.WaitGroup
			for _, target := range cfg.Targets {
				wg.Add(1)
				go func(t string) {
					defer wg.Done()
					checkStatus(t, cfg.OnFailure, cfg.Retries)
				}(target)
			}
			wg.Wait()
		case <-sigChan:
			fmt.Println("\nShutting down gracefully...")
			stats := metrics.GetStats()
			fmt.Printf("Final stats: %+v\n", stats)
			return
		}
	}
}

func checkStatus(target string, onFailure string, maxRetries int) {
	var resp *http.Response
	var err error
	
	start := time.Now()
	
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying check for %s (%d/%d)...", target, i, maxRetries)
			time.Sleep(1 * time.Second)
		}
		
		resp, err = http.Get(target + "/health")
		if err == nil {
			break
		}
	}
	
	duration := time.Since(start)
	
	if err != nil {
		log.Printf("DOWN: %s (%v) - Latency: %v", target, err, duration)
		metrics.RecordDown(target, duration)
		executeRemediation(onFailure)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("UP: %s (Status: %d) - Latency: %v", target, resp.StatusCode, duration)
		metrics.RecordUp(target, duration)
	} else {
		log.Printf("DOWN: %s (Status: %d) - Latency: %v", target, resp.StatusCode, duration)
		metrics.RecordDown(target, duration)
		executeRemediation(onFailure)
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
