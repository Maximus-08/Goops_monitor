package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Interval time.Duration
	Target   string
}

func main() {
	cfg := Config{
		Interval: 5 * time.Second,
		Target:   "http://localhost:8080",
	}

	fmt.Printf("Starting monitor for %s with interval %v\n", cfg.Target, cfg.Interval)

	// Start a local test server for demonstration
	go startServer(":8080")

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		checkStatus(cfg.Target)
	}
}

func checkStatus(target string) {
	resp, err := http.Get(target + "/health")
	if err != nil {
		log.Printf("DOWN: %s (%v)", target, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("UP: %s (Status: %d)", target, resp.StatusCode)
	} else {
		log.Printf("DOWN: %s (Status: %d)", target, resp.StatusCode)
	}
}
