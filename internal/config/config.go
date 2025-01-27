package config

import (
	"encoding/json"
	"os"
	"time"
)

type CoinConfig struct {
	ID        string  `json:"id"`
	Threshold float64 `json:"threshold"`
}

type Config struct {
	Coins                  []CoinConfig `json:"coins"`
	IntervalSeconds        int          `json:"intervalSeconds"`
	PriceMoveWindowMinutes int          `json:"priceMoveWindowMinutes"`
	PriceMovePercentage    float64      `json:"priceMovePercentage"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Optional: set defaults if zero
	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 10
	}
	if cfg.PriceMoveWindowMinutes <= 0 {
		cfg.PriceMoveWindowMinutes = 1
	}
	if cfg.PriceMovePercentage <= 0 {
		cfg.PriceMovePercentage = 10
	}

	return &cfg, nil
}

func (c *Config) IntervalDuration() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}

func (c *Config) PriceMoveWindowDuration() time.Duration {
	return time.Duration(c.PriceMoveWindowMinutes) * time.Minute
}
