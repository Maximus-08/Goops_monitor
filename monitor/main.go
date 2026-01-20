package main

import (
	"fmt"
	"log"
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

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		checkStatus(cfg.Target)
	}
}

func checkStatus(target string) {
	// TODO: Implement actual HTTP check
	log.Printf("Checking status of %s...", target)
}
