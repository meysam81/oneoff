package config

import (
	"fmt"
	"runtime"

	"github.com/caarlos0/env/v11"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Port int    `env:"PORT" envDefault:"8080"`
	Host string `env:"HOST" envDefault:"localhost"`

	// Database configuration
	DBPath string `env:"DB_PATH" envDefault:"./oneoff.db"`

	// Worker configuration
	WorkersCount int `env:"WORKERS_COUNT" envDefault:"0"` // 0 = auto (N/2 cores)

	// Logging configuration
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// Timezone configuration
	DefaultTimezone string `env:"DEFAULT_TIMEZONE" envDefault:"UTC"`

	// Retention configuration
	LogRetentionDays int `env:"LOG_RETENTION_DAYS" envDefault:"90"`

	// Default priority
	DefaultPriority int `env:"DEFAULT_PRIORITY" envDefault:"5"`

	// Environment
	Environment string `env:"ENVIRONMENT" envDefault:"production"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Auto-detect worker count if set to 0
	if cfg.WorkersCount == 0 {
		cfg.WorkersCount = runtime.NumCPU() / 2
		if cfg.WorkersCount < 1 {
			cfg.WorkersCount = 1
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", c.Port)
	}

	if c.DBPath == "" {
		return fmt.Errorf("DB_PATH cannot be empty")
	}

	if c.WorkersCount < 1 {
		return fmt.Errorf("WORKERS_COUNT must be at least 1")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	if c.DefaultPriority < 1 || c.DefaultPriority > 10 {
		return fmt.Errorf("DEFAULT_PRIORITY must be between 1 and 10")
	}

	if c.LogRetentionDays < 1 {
		return fmt.Errorf("LOG_RETENTION_DAYS must be at least 1")
	}

	return nil
}

// Address returns the full server address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
