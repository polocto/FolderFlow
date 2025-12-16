package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type VerboseHandler struct {
	level slog.Level
}

func NewVerboseHandler(level slog.Level) *VerboseHandler {
	return &VerboseHandler{level: level}
}

func (h *VerboseHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *VerboseHandler) Handle(ctx context.Context, r slog.Record) error {
	// Print level and message
	fmt.Fprintf(os.Stderr, "[%s] %s", r.Level, r.Message)

	// Iterate attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != slog.TimeKey { // skip timestamp
			fmt.Fprintf(os.Stderr, " %s=%v", a.Key, a.Value)
		}
		return true
	})

	fmt.Fprintln(os.Stderr) // newline
	return nil
}

func (h *VerboseHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *VerboseHandler) WithGroup(name string) slog.Handler       { return h }
