package main

import (
	"log"
	"sync"
)

type Metrics struct {
	Ups   int
	Downs int
	mu    sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordUp() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Ups++
	log.Println("Recorded UP status")
}

func (m *Metrics) RecordDown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Downs++
	log.Println("Recorded DOWN status")
}

func (m *Metrics) GetStats() map[string]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]int{
		"ups":   m.Ups,
		"downs": m.Downs,
	}
}
