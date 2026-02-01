package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Alerter struct {
	webhookURL string
	cooldown   time.Duration
	lastAlert  map[string]time.Time
	mu         sync.Mutex
}

func NewAlerter(webhookURL string, cooldown time.Duration) *Alerter {
	return &Alerter{
		webhookURL: webhookURL,
		cooldown:   cooldown,
		lastAlert:  make(map[string]time.Time),
	}
}

func (a *Alerter) SendAlert(target string, message string) {
	if a.webhookURL == "" {
		return
	}

	a.mu.Lock()
	last, ok := a.lastAlert[target]
	if ok && time.Since(last) < a.cooldown {
		a.mu.Unlock()
		return // Cooldown active, skip alert
	}
	a.lastAlert[target] = time.Now()
	a.mu.Unlock()

	go func() {
		payload := map[string]string{
			"text": fmt.Sprintf("🚨 **ALERT** 🚨\nTarget: %s\nError: %s", target, message),
		}
		
		body, _ := json.Marshal(payload)
		resp, err := http.Post(a.webhookURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Failed to send alert: %v", err)
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode >= 400 {
			log.Printf("Webhook failed with status: %d", resp.StatusCode)
		} else {
			log.Printf("Alert sent for %s", target)
		}
	}()
}
