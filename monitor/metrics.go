package main

import (
	"log"
	"sync"
)

type Metrics struct {
	Stats map[string]*TargetStats
	mu    sync.Mutex
}

type TargetStats struct {
	Ups   int `json:"ups"`
	Downs int `json:"downs"`
}

func NewMetrics() *Metrics {
	return &Metrics{
		Stats: make(map[string]*TargetStats),
	}
}

func (m *Metrics) RecordUp(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, ok := m.Stats[target]; !ok {
		m.Stats[target] = &TargetStats{}
	}
	m.Stats[target].Ups++
	log.Printf("Recorded UP status for %s", target)
}

func (m *Metrics) RecordDown(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, ok := m.Stats[target]; !ok {
		m.Stats[target] = &TargetStats{}
	}
	m.Stats[target].Downs++
	log.Printf("Recorded DOWN status for %s", target)
}

func (m *Metrics) GetStats() map[string]*TargetStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Return a copy to be safe
	copy := make(map[string]*TargetStats)
	for k, v := range m.Stats {
		copy[k] = &TargetStats{Ups: v.Ups, Downs: v.Downs}
	}
	return copy
}
