package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port    string
	GinMode string

	// Database configuration
	DatabaseURL string

	// Logging configuration
	LogLevel string

	// Environment
	Environment string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:        getEnvWithDefault("PORT", "8080"),
		GinMode:     getEnvWithDefault("GIN_MODE", "debug"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
		Environment: getEnvWithDefault("ENV", "development"),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate ensures all required configuration is present
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// Validate GinMode
	if c.GinMode != "debug" && c.GinMode != "release" {
		return fmt.Errorf("GIN_MODE must be 'debug' or 'release', got '%s'", c.GinMode)
	}

	// Validate LogLevel
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error, got '%s'", c.LogLevel)
	}

	return nil
}

// getEnvWithDefault returns the value of an environment variable or a default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
