package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVerboseHandler(t *testing.T) {
	handler := NewVerboseHandler(slog.LevelInfo)
	assert.NotNil(t, handler)
	assert.Equal(t, slog.LevelInfo, handler.level)
}

func TestVerboseHandler_Enabled(t *testing.T) {
	handler := NewVerboseHandler(slog.LevelInfo)

	// Test with Info level
	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo))

	// Test with Debug level
	assert.False(t, handler.Enabled(context.Background(), slog.LevelDebug))

	// Test with Error level
	assert.True(t, handler.Enabled(context.Background(), slog.LevelError))
}

func TestVerboseHandler_Handle(t *testing.T) {
	// Redirect stderr to a buffer for testing
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()
	r, w, _ := os.Pipe()
	os.Stderr = w

	handler := NewVerboseHandler(slog.LevelInfo)

	// Create a record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.AddAttrs(
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
	)

	// Handle the record
	err := handler.Handle(context.Background(), record)
	require.NoError(t, err)

	// Close the write end of the pipe
	w.Close()

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check if the output contains the expected message and attributes
	assert.Contains(t, output, "[INFO] test message")
	assert.Contains(t, output, "key1=value1")
	assert.Contains(t, output, "key2=42")
}

func TestVerboseHandler_WithAttrs(t *testing.T) {
	handler := NewVerboseHandler(slog.LevelInfo)

	// Add attributes
	attrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.String("key2", "value2"),
	}

	// Call WithAttrs
	newHandler := handler.WithAttrs(attrs)

	// Check if the returned handler is the same type
	assert.IsType(t, &VerboseHandler{}, newHandler)
}

func TestVerboseHandler_WithGroup(t *testing.T) {
	handler := NewVerboseHandler(slog.LevelInfo)

	// Call WithGroup
	newHandler := handler.WithGroup("testGroup")

	// Check if the returned handler is the same type
	assert.IsType(t, &VerboseHandler{}, newHandler)
}

func TestVerboseHandler_LogLevel(t *testing.T) {
	// Redirect stderr to a buffer for testing
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()
	r, w, _ := os.Pipe()
	os.Stderr = w

	handler := NewVerboseHandler(slog.LevelInfo)
	logger := slog.New(handler)

	// Log messages at different levels
	logger.Info("info message")
	logger.Debug("debug message")
	logger.Error("error message")

	// Close the write end of the pipe
	w.Close()

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check if the output contains the expected messages
	assert.Contains(t, output, "[INFO] info message")
	assert.NotContains(t, output, "[DEBUG] debug message")
	assert.Contains(t, output, "[ERROR] error message")
}
