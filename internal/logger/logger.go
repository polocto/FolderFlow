package logger

import (
	"context"
	"log/slog"
	"os"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		hs[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: hs}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		hs[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: hs}
}

func Init(verbose, debug bool) (func() error, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		"logs/folderflow.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}
	var leveler slog.Leveler = slog.LevelInfo
	if debug {
		leveler = slog.LevelDebug
	}

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: leveler,
	})

	var handler slog.Handler

	if verbose {
		consoleHandler := NewVerboseHandler(leveler.Level())

		// Combine
		handler = &MultiHandler{
			handlers: []slog.Handler{
				fileHandler,
				consoleHandler,
			},
		}
	} else {
		handler = fileHandler
	}

	slog.SetDefault(slog.New(handler))
	return file.Close, nil
}
