package config

import (
	"fmt"
	"os"
	"time"
	"strconv"
)

type Config struct {
	Server ServerConfig
	Stats StatsConfig
	LogLevel string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type StatsConfig struct {
	WindowSeconds int
}

const (
	defaultPort         = "8080"
	defaultStatsWindowSeconds = 60
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 15 * time.Second
	defaultLogLevel = "info"
)

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnvString("PORT", defaultPort),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", defaultReadTimeout),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", defaultWriteTimeout),
			IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", defaultIdleTimeout),
		},
		Stats: StatsConfig{
			WindowSeconds: getEnvInt("STATS_WINDOW_SECONDS", defaultStatsWindowSeconds),
		},
		LogLevel: getEnvString("LOG_LEVEL", defaultLogLevel),
	}

	// Validação das configurações
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("erro na validação das configurações: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Stats.WindowSeconds <= 0 {
		return fmt.Errorf("STATS_WINDOW_SECONDS deve ser maior que zero")
	}

	if c.Server.Port == "" {
		return fmt.Errorf("PORT não pode ser vazio")
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("READ_TIMEOUT deve ser maior que zero")
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("WRITE_TIMEOUT deve ser maior que zero")
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("IDLE_TIMEOUT deve ser maior que zero")
	}

	return nil
}

func getEnvString(key, defaultValue string) string {
	
	if value := os.Getenv(key); value == "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
