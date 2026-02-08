package main

import (
	"goops-monitor/runner"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var metrics *Metrics
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {
	// Initialize structured logger (JSON format for production)
	jsonFormat := os.Getenv("GOOPS_LOG_JSON") == "true"
	InitLogger(jsonFormat)

	// Load configuration
	configPath := os.Getenv("GOOPS_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		LogError("Failed to load config", "error", err, "path", configPath)
		os.Exit(1)
	}

	// Initialize metrics
	metrics = NewMetrics()

	// Initialize alerter
	alerter := NewAlerter(cfg.WebhookURL, cfg.AlertCooldown)

	LogInfo("Starting monitor",
		"targets", len(cfg.Targets),
		"interval", cfg.Interval.String(),
	)

	// Start API server
	apiPort := os.Getenv("GOOPS_API_PORT")
	if apiPort == "" {
		apiPort = ":8081"
	}
	go StartAPI(apiPort)

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
					checkStatus(t, cfg.OnFailure, cfg.Retries, alerter)
				}(target)
			}
			wg.Wait()
		case <-sigChan:
			LogInfo("Shutting down gracefully")
			stats := metrics.GetStats()
			for target, s := range stats {
				LogInfo("Final stats", "target", target, "ups", s.Ups, "downs", s.Downs)
			}
			return
		}
	}
}

func checkStatus(target string, onFailure string, maxRetries int, alerter *Alerter) {
	var resp *http.Response
	var err error
	
	start := time.Now()
	
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			LogInfo("Retrying check", "target", target, "attempt", i, "max", maxRetries)
			time.Sleep(1 * time.Second)
		}
		
		resp, err = httpClient.Get(target + "/health")
		if err == nil {
			break
		}
	}
	
	duration := time.Since(start)
	
	if err != nil {
		LogError("Target DOWN", "target", target, "error", err.Error(), "latency_ms", duration.Milliseconds())
		metrics.RecordDown(target, duration)
		alerter.SendAlert(target, err.Error())
		executeRemediation(onFailure)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		LogInfo("Target UP", "target", target, "status", resp.StatusCode, "latency_ms", duration.Milliseconds())
		metrics.RecordUp(target, duration)
	} else {
		LogError("Target DOWN", "target", target, "status", resp.StatusCode, "latency_ms", duration.Milliseconds())
		metrics.RecordDown(target, duration)
		alerter.SendAlert(target, "Status code: "+http.StatusText(resp.StatusCode))
		executeRemediation(onFailure)
	}
}

func executeRemediation(script string) {
	if script == "" {
		return
	}
	
	LogInfo("Executing remediation", "script", script)
	r := runner.New("sh", "-c", script)
	if err := r.Execute(); err != nil {
		LogError("Remediation failed", "script", script, "error", err.Error())
	}
}
