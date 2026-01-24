package main

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Interval  time.Duration `json:"interval"`
	Target    string        `json:"target"`
	OnFailure string        `json:"on_failure"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Config
func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		Interval string `json:"interval"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	duration, err := time.ParseDuration(aux.Interval)
	if err != nil {
		return err
	}
	c.Interval = duration
	
	return nil
}

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
