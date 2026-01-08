package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type slogGormLogger struct {
	logger                    *slog.Logger
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func newSlogGormLogger(l *slog.Logger, lvl logger.LogLevel) logger.Interface {
	return &slogGormLogger{
		logger:                    l.With("component", "database"),
		LogLevel:                  lvl,
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             200 * time.Millisecond,
	}
}

func (l *slogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *slogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.InfoContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *slogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *slogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *slogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)) {
		if l.LogLevel >= logger.Error {
			l.logger.ErrorContext(ctx, "query error", "err", err, "elapsed", elapsed, "rows", rows, "sql", sql)
		}
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn {
		l.logger.WarnContext(ctx, "slow query", "elapsed", elapsed, "rows", rows, "sql", sql)
		return
	}

	if l.LogLevel >= logger.Info {
		l.logger.InfoContext(ctx, "query", "elapsed", elapsed, "rows", rows, "sql", sql)
	}
}
