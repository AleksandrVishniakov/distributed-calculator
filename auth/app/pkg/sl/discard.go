package sl

import (
	"context"
	"log/slog"
)

// NewDiscardLogger returns a slog.Logger implementation with no functionality
func NewDiscardLogger() *slog.Logger {
	return slog.New(DiscardHandler{})
}

type DiscardHandler struct {
}

func (d DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (d DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (d DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return d
}

func (d DiscardHandler) WithGroup(_ string) slog.Handler {
	return d
}
