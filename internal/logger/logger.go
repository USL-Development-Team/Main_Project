package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// SetupLogger configures structured logging for the application
func SetupLogger(level slog.Level, logFilePath string) *slog.Logger {
	var writers []io.Writer

	// Always write to stdout
	writers = append(writers, os.Stdout)

	// Also write to file if path provided
	if logFilePath != "" {
		// Create logs directory if it doesn't exist
		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// Fallback to stdout only if we can't create log directory
			writers = []io.Writer{os.Stdout}
		} else {
			// Open log file (create if doesn't exist, append if it does)
			logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				// Fallback to stdout only if we can't open log file
				writers = []io.Writer{os.Stdout}
			} else {
				writers = append(writers, logFile)
			}
		}
	}

	// Create multi-writer to write to both stdout and file
	output := io.MultiWriter(writers...)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // Include source file/line in logs
	}

	handler := slog.NewJSONHandler(output, opts)
	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return logger
}

// NewLogger creates a new logger with additional context
func NewLogger(name string) *slog.Logger {
	return slog.Default().With("component", name)
}
