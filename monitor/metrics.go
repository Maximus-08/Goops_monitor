package main

import (
	"log"
	"sync"
	"time"
)

type Metrics struct {
	Stats map[string]*TargetStats
	mu    sync.Mutex
}

type TargetStats struct {
	Ups           int `json:"ups"`
	Downs         int `json:"downs"`
	TotalDuration int64 `json:"total_duration_ms"` // in milliseconds
	LastLatency   int64 `json:"last_latency_ms"`   // in milliseconds
}

func NewMetrics() *Metrics {
	return &Metrics{
		Stats: make(map[string]*TargetStats),
	}
}

func (m *Metrics) RecordUp(target string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, ok := m.Stats[target]; !ok {
		m.Stats[target] = &TargetStats{}
	}
	s := m.Stats[target]
	s.Ups++
	s.LastLatency = duration.Milliseconds()
	s.TotalDuration += duration.Milliseconds()
	
	log.Printf("Recorded UP status for %s (Latency: %dms)", target, s.LastLatency)
}

func (m *Metrics) RecordDown(target string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, ok := m.Stats[target]; !ok {
		m.Stats[target] = &TargetStats{}
	}
	s := m.Stats[target]
	s.Downs++
	s.LastLatency = duration.Milliseconds()
	s.TotalDuration += duration.Milliseconds()
	
	log.Printf("Recorded DOWN status for %s (Latency: %dms)", target, s.LastLatency)
}

func (m *Metrics) GetStats() map[string]*TargetStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Return a copy to be safe
	copy := make(map[string]*TargetStats)
	for k, v := range m.Stats {
		copy[k] = &TargetStats{
			Ups:           v.Ups,
			Downs:         v.Downs,
			TotalDuration: v.TotalDuration,
			LastLatency:   v.LastLatency,
		}
	}
	return copy
}
