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

func TestMultiHandler_Enabled(t *testing.T) {
	// Create a MultiHandler with two handlers
	handler1 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler2 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{handler1, handler2},
	}

	// Test with Info level
	ctx := context.Background()
	assert.True(t, multiHandler.Enabled(ctx, slog.LevelInfo))

	// Test with Debug level
	assert.True(t, multiHandler.Enabled(ctx, slog.LevelDebug))

	// Test with Error level
	assert.True(t, multiHandler.Enabled(ctx, slog.LevelError))
}

func TestMultiHandler_Handle(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a handler that writes to the buffer
	handler1 := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler2 := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{handler1, handler2},
	}

	// Create a record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := multiHandler.Handle(context.Background(), record)
	require.NoError(t, err)

	// Check if the message was written to the buffer
	assert.Contains(t, buf.String(), "test message")
}

func TestMultiHandler_WithAttrs(t *testing.T) {
	// Create a MultiHandler with two handlers
	handler1 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler2 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{handler1, handler2},
	}

	// Add attributes
	attrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.String("key2", "value2"),
	}

	// Call WithAttrs
	newHandler := multiHandler.WithAttrs(attrs)

	// Check if the returned handler is a MultiHandler
	assert.IsType(t, &MultiHandler{}, newHandler)
}

func TestMultiHandler_WithGroup(t *testing.T) {
	// Create a MultiHandler with two handlers
	handler1 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler2 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{handler1, handler2},
	}

	// Call WithGroup
	newHandler := multiHandler.WithGroup("testGroup")

	// Check if the returned handler is a MultiHandler
	assert.IsType(t, &MultiHandler{}, newHandler)
}

func TestInit(t *testing.T) {
	// Create a temporary directory for logs
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	defer require.NoError(t, os.RemoveAll(logDir))

	// Test verbose mode
	err := Init(true, false)
	require.NoError(t, err)

	// Check if the log file was created
	_, err = os.Stat("logs/folderflow.log")
	assert.NoError(t, err)

	// Test non-verbose mode
	err = Init(false, false)
	require.NoError(t, err)

	// Check if the log file was created
	_, err = os.Stat("logs/folderflow.log")
	assert.NoError(t, err)
}

func TestLoggerOutput(t *testing.T) {
	// Create a temporary directory for logs
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	defer require.NoError(t, os.RemoveAll(logDir))

	// Initialize logger
	err := Init(true, false)
	require.NoError(t, err)

	// Log a message
	slog.Info("test info message")
	slog.Error("test error message")

	// Check if the log file contains the messages
	fileContent, err := os.ReadFile("logs/folderflow.log")
	require.NoError(t, err)

	assert.Contains(t, string(fileContent), "test info message")
	assert.Contains(t, string(fileContent), "test error message")
}
