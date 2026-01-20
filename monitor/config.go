package main

import (
	"encoding/json"
	"os"
	"time"
)

// LoadConfig reads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	if cfg.Interval == 0 {
		cfg.Interval = 10 * time.Second
	}

	return &cfg, nil
}
