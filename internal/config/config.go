package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration. Loaded once at startup; fail fast if invalid.
type Config struct {
	GRPCPort int
}

// Load reads configuration from environment. Returns error if required values are missing or invalid.
func Load() (*Config, error) {
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		portStr = os.Getenv("HTTP_PORT")
	}
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return nil, fmt.Errorf("GRPC_PORT (or HTTP_PORT) must be a valid port (1-65535), got %q", portStr)
	}
	return &Config{GRPCPort: port}, nil
}
