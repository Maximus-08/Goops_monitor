package main

import (
	"sync"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Stats         map[string]*TargetStats
	upGauge       *prometheus.GaugeVec
	latencyHist   *prometheus.HistogramVec
	mu            sync.Mutex
}

type TargetStats struct {
	Ups           int `json:"ups"`
	Downs         int `json:"downs"`
	TotalDuration int64 `json:"total_duration_ms"` // in milliseconds
	LastLatency   int64 `json:"last_latency_ms"`   // in milliseconds
}

func NewMetrics() *Metrics {
	m := &Metrics{
		Stats: make(map[string]*TargetStats),
		upGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goops_target_up",
				Help: "Current status of the target (1 = UP, 0 = DOWN)",
			},
			[]string{"target"},
		),
		latencyHist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "goops_target_latency_seconds",
				Help:    "Response latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"target"},
		),
	}
	
	prometheus.MustRegister(m.upGauge)
	prometheus.MustRegister(m.latencyHist)
	return m
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
	
	// Prometheus updates
	m.upGauge.WithLabelValues(target).Set(1)
	m.latencyHist.WithLabelValues(target).Observe(duration.Seconds())
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
	
	// Prometheus updates
	m.upGauge.WithLabelValues(target).Set(0)
	m.latencyHist.WithLabelValues(target).Observe(duration.Seconds())
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
