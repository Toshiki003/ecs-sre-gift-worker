package config

import (
	"fmt"
	"os"
)

// Config はアプリケーション設定を保持する。
// 環境変数から読み込む（T-2.2 で詳細実装）。
type Config struct {
	Port        string
	DatabaseURL string
	SQSQueueURL string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	sqsURL := os.Getenv("SQS_QUEUE_URL")

	_ = fmt.Sprintf("loaded config: port=%s", port)

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		SQSQueueURL: sqsURL,
	}, nil
}
