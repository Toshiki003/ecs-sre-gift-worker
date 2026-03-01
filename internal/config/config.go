package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション設定を保持する。
// すべての設定は環境変数から読み込む。
type Config struct {
	// Server
	Port string

	// Database
	DatabaseURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	// SQS
	SQSQueueURL string
	AWSRegion   string

	// Worker
	FailureRate int

	// Logging
	LogLevel string
}

// Load は環境変数から設定を読み込む。
func Load() (*Config, error) {
	cfg := &Config{
		Port:              envOrDefault("PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		DBMaxOpenConns:    envIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    envIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: envDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		SQSQueueURL:       os.Getenv("SQS_QUEUE_URL"),
		AWSRegion:         envOrDefault("AWS_REGION", "ap-northeast-1"),
		FailureRate:       envIntOrDefault("FAILURE_RATE", 20),
		LogLevel:          envOrDefault("LOG_LEVEL", "info"),
	}

	if cfg.FailureRate < 0 || cfg.FailureRate > 100 {
		return nil, fmt.Errorf("FAILURE_RATE must be between 0 and 100, got %d", cfg.FailureRate)
	}

	return cfg, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envIntOrDefault(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

func envDurationOrDefault(key string, defaultVal time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultVal
	}
	return d
}
