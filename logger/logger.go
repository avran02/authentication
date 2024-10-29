package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/avran02/authentication/internal/config"
)

func Setup(config config.Server) {
	var ll slog.Leveler
	var isDefaultLogLevel bool
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		ll = slog.LevelDebug
	case "info":
		ll = slog.LevelInfo
	case "warn":
		ll = slog.LevelWarn
	case "error":
		ll = slog.LevelError
	default:
		isDefaultLogLevel = true
		ll = slog.LevelInfo
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: ll,
	}))

	slog.SetDefault(log)

	if isDefaultLogLevel {
		slog.Warn("Logger using default value", "log level", ll)
		return
	}
	slog.Info("Logger", "log level", ll)
}
