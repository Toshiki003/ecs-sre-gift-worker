package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// 環境変数をクリアしてデフォルト値を検証
	for _, key := range []string{"PORT", "DATABASE_URL", "SQS_QUEUE_URL", "AWS_REGION", "FAILURE_RATE", "LOG_LEVEL", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME"} {
		os.Unsetenv(key)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.DatabaseURL != "" {
		t.Errorf("DatabaseURL = %q, want empty", cfg.DatabaseURL)
	}
	if cfg.SQSQueueURL != "" {
		t.Errorf("SQSQueueURL = %q, want empty", cfg.SQSQueueURL)
	}
	if cfg.AWSRegion != "ap-northeast-1" {
		t.Errorf("AWSRegion = %q, want %q", cfg.AWSRegion, "ap-northeast-1")
	}
	if cfg.FailureRate != 20 {
		t.Errorf("FailureRate = %d, want %d", cfg.FailureRate, 20)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.DBMaxOpenConns != 25 {
		t.Errorf("DBMaxOpenConns = %d, want %d", cfg.DBMaxOpenConns, 25)
	}
	if cfg.DBMaxIdleConns != 5 {
		t.Errorf("DBMaxIdleConns = %d, want %d", cfg.DBMaxIdleConns, 5)
	}
	if cfg.DBConnMaxLifetime != 5*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want %v", cfg.DBConnMaxLifetime, 5*time.Minute)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/giftdb")
	os.Setenv("SQS_QUEUE_URL", "https://sqs.ap-northeast-1.amazonaws.com/123456789012/gift-queue")
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("FAILURE_RATE", "50")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DB_MAX_OPEN_CONNS", "10")
	os.Setenv("DB_MAX_IDLE_CONNS", "3")
	os.Setenv("DB_CONN_MAX_LIFETIME", "10m")
	defer func() {
		for _, key := range []string{"PORT", "DATABASE_URL", "SQS_QUEUE_URL", "AWS_REGION", "FAILURE_RATE", "LOG_LEVEL", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME"} {
			os.Unsetenv(key)
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/giftdb" {
		t.Errorf("DatabaseURL = %q, want postgres URL", cfg.DatabaseURL)
	}
	if cfg.AWSRegion != "us-west-2" {
		t.Errorf("AWSRegion = %q, want %q", cfg.AWSRegion, "us-west-2")
	}
	if cfg.FailureRate != 50 {
		t.Errorf("FailureRate = %d, want %d", cfg.FailureRate, 50)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.DBMaxOpenConns != 10 {
		t.Errorf("DBMaxOpenConns = %d, want %d", cfg.DBMaxOpenConns, 10)
	}
	if cfg.DBMaxIdleConns != 3 {
		t.Errorf("DBMaxIdleConns = %d, want %d", cfg.DBMaxIdleConns, 3)
	}
	if cfg.DBConnMaxLifetime != 10*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want %v", cfg.DBConnMaxLifetime, 10*time.Minute)
	}
}

func TestLoad_InvalidFailureRate(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"negative", "-1"},
		{"over_100", "101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("FAILURE_RATE", tt.value)
			defer os.Unsetenv("FAILURE_RATE")

			_, err := Load()
			if err == nil {
				t.Error("expected error for invalid FAILURE_RATE, got nil")
			}
		})
	}
}

func TestLoad_InvalidIntFallsBackToDefault(t *testing.T) {
	os.Setenv("DB_MAX_OPEN_CONNS", "not-a-number")
	defer os.Unsetenv("DB_MAX_OPEN_CONNS")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DBMaxOpenConns != 25 {
		t.Errorf("DBMaxOpenConns = %d, want default %d", cfg.DBMaxOpenConns, 25)
	}
}

func TestLoad_InvalidDurationFallsBackToDefault(t *testing.T) {
	os.Setenv("DB_CONN_MAX_LIFETIME", "not-a-duration")
	defer os.Unsetenv("DB_CONN_MAX_LIFETIME")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DBConnMaxLifetime != 5*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want default %v", cfg.DBConnMaxLifetime, 5*time.Minute)
	}
}
