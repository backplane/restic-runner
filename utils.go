package main

import (
	"log/slog"
	"strings"
)

// setLogLevel sets the log level
func setLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "INFO":
		logLevel.Set(slog.LevelInfo)
	case "WARN":
		logLevel.Set(slog.LevelWarn)
	case "ERROR":
		logLevel.Set(slog.LevelError)
	// case "DEBUG":
	default:
		logLevel.Set(slog.LevelDebug)
	}
}
