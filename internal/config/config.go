package config

import (
	"errors"
	"os"
)

// Config holds all configuration values for the server.
type Config struct {
	Token     string // STATUSCAST_TOKEN (required)
	Domain    string // STATUSCAST_DOMAIN e.g. "myco.statuscast.com" (required)
	Transport string // TRANSPORT: "stdio" (default) | "http"
	Port      string // PORT: default "8080" (only used when Transport == "http")
}

// Load reads configuration from environment variables and returns an error
// if any required variables are missing.
func Load() (*Config, error) {
	cfg := &Config{
		Token:     os.Getenv("STATUSCAST_TOKEN"),
		Domain:    os.Getenv("STATUSCAST_DOMAIN"),
		Transport: os.Getenv("TRANSPORT"),
		Port:      os.Getenv("PORT"),
	}

	if cfg.Token == "" {
		return nil, errors.New("STATUSCAST_TOKEN environment variable is required")
	}
	if cfg.Domain == "" {
		return nil, errors.New("STATUSCAST_DOMAIN environment variable is required")
	}

	if cfg.Transport == "" {
		cfg.Transport = "stdio"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
