package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Logger is a wrapper around library logger.
type Logger struct {
	log *slog.Logger
}

// New creates a new Logger instance.
// If there is no such level, the default level is info.
func New(level string) *Logger {
	var l slog.Level

	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: l}

	log := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	return &Logger{log: log}
}

// Debug logs a message at the debug level.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log.Debug(msg, args...)
}

// Info logs a message at the info level.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log.Info(msg, args...)
}

// Warn logs a message at the warn level.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log.Warn(msg, args...)
}

// Error logs a message at the error level.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log.Error(msg, args...)
}
