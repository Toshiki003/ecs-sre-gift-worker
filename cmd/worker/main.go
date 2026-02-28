package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Toshiki003/ecs-sre-gift-worker/internal/config"
	"github.com/Toshiki003/ecs-sre-gift-worker/internal/logger"
)

func main() {
	_, err := config.Load()
	if err != nil {
		l := logger.New("worker")
		l.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	l := logger.New("worker")
	l.Info("starting worker")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: T-4.1 で SQS ポーリングループを実装
	_ = ctx

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("shutting down worker")
	cancel()
}
