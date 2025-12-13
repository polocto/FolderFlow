package logger

import (
	"log/slog"
	"os"
)

// Log is the global logger instance.
var Log *slog.Logger

func init() {
	// Open a file for logging
	file, err := os.OpenFile("folderflow.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Create a new logger that writes to the file
	Log = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// Info logs a message at Info level.
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

// Warn logs a message at Warn level.
func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

// Error logs a message at Error level.
func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

// Debug logs a message at Debug level.
func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}
