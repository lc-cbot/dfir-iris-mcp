package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BaseURL string
	APIKey  string
}

func Load() (*Config, error) {
	u := os.Getenv("DFIR_IRIS_URL")
	if u == "" {
		return nil, fmt.Errorf("DFIR_IRIS_URL environment variable is required")
	}
	k := os.Getenv("DFIR_IRIS_API_KEY")
	if k == "" {
		return nil, fmt.Errorf("DFIR_IRIS_API_KEY environment variable is required")
	}
	return &Config{
		BaseURL: strings.TrimRight(u, "/"),
		APIKey:  k,
	}, nil
}
