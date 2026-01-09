package logger

import (
	"log/slog"
	"os"
)

// NewCommandLogger creates a new slog logger for commands.
// It removes the timestamp from the log output to keep it clean for CLI usage.
func NewCommandLogger(level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	return slog.New(slog.NewTextHandler(os.Stderr, opts))
}
