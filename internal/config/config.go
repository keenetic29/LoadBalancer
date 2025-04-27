package config

import (
	"encoding/json"
	"os"
)

type RateLimitConfig struct {
	DefaultCapacity   int  `json:"default_capacity"`
	DefaultRatePerSec int  `json:"default_rate_per_sec"`
	RefillPeriod      int  `json:"refill_period"`
}

type Config struct {
	Port      string         `json:"port"`
	Backends  []string       `json:"backends"`
	RateLimit RateLimitConfig `json:"rate_limit"`
	ClientsDB string         `json:"clients_db"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}