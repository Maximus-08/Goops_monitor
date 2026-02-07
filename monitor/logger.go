package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger(jsonFormat bool) {
	var handler slog.Handler

	if jsonFormat {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger = slog.New(handler)
}

// LogInfo logs an info message with optional attributes
func LogInfo(msg string, attrs ...any) {
	if logger == nil {
		InitLogger(false)
	}
	logger.Info(msg, attrs...)
}

// LogError logs an error message with optional attributes
func LogError(msg string, attrs ...any) {
	if logger == nil {
		InitLogger(false)
	}
	logger.Error(msg, attrs...)
}

// LogWarn logs a warning message with optional attributes
func LogWarn(msg string, attrs ...any) {
	if logger == nil {
		InitLogger(false)
	}
	logger.Warn(msg, attrs...)
}
