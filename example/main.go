package main

import (
	"log/slog"
	"os"

	"github.com/kechako/seezlog"
)

func main() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	h := seezlog.NewHandler(os.Stdout, true, opts)
	logger := slog.New(h)

	logger.Debug("Debug message")
	logger.Info("Info message", "user_id", 12345)
	logger.Warn("Warning message", slog.String("ip_address", "192.168.1.100"))
	logger.Error("Error message", slog.Any("error", os.ErrNotExist))

	// With attributes and groups
	loggerWithAttrs := logger.With("service", "auth", "version", "v1.2.3")
	loggerWithAttrs.Info("User authenticated")

	groupLogger := logger.WithGroup("request")
	groupLogger.Info("Request received", "method", "GET", "path", "/api/users")
}
