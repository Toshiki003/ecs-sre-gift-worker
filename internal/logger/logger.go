package logger

import (
	"log/slog"
	"os"
)

// Logger はJSON構造化ログを提供する。
// 要件8.3に基づき、service フィールドを常時付与する。
// gift_id, request_id 等の付与は T-2.3 で詳細実装する。
type Logger struct {
	*slog.Logger
}

func New(service string) *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	l := slog.New(handler).With("service", service)
	return &Logger{Logger: l}
}
