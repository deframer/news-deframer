package logger

import (
	"context"

	"goa.design/clue/log"
)

// NewLoggerContext initializes a clue/log context for command-line tools.
func NewLoggerContext(ctx context.Context, debug bool) context.Context {
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx = log.Context(ctx, log.WithFormat(format))
	if debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}
	return ctx
}
