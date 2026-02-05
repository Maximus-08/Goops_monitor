package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Interval  time.Duration `json:"interval"`
	Targets       []string      `json:"targets"`
	Retries       int           `json:"retries"`
	WebhookURL    string        `json:"webhook_url"`
	AlertCooldown time.Duration `json:"alert_cooldown"`
	OnFailure     string        `json:"on_failure"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Config
func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		Interval      string `json:"interval"`
		AlertCooldown string `json:"alert_cooldown"`
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

	if aux.AlertCooldown != "" {
		cooldown, err := time.ParseDuration(aux.AlertCooldown)
		if err != nil {
			return err
		}
		c.AlertCooldown = cooldown
	}
	
	return nil
}

// LoadConfig reads configuration from a JSON file and applies env overrides
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	applyEnvOverrides(&cfg)

	if cfg.Interval == 0 {
		cfg.Interval = 10 * time.Second
	}

	return &cfg, nil
}

// applyEnvOverrides overrides config values with GOOPS_* environment variables
func applyEnvOverrides(cfg *Config) {
	if val := os.Getenv("GOOPS_INTERVAL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.Interval = d
		}
	}

	if val := os.Getenv("GOOPS_TARGETS"); val != "" {
		cfg.Targets = strings.Split(val, ",")
	}

	if val := os.Getenv("GOOPS_RETRIES"); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			cfg.Retries = n
		}
	}

	if val := os.Getenv("GOOPS_WEBHOOK_URL"); val != "" {
		cfg.WebhookURL = val
	}

	if val := os.Getenv("GOOPS_ALERT_COOLDOWN"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.AlertCooldown = d
		}
	}

	if val := os.Getenv("GOOPS_ON_FAILURE"); val != "" {
		cfg.OnFailure = val
	}
}
