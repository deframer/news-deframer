package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"goa.design/clue/log"
	"gorm.io/gorm/logger"
	gorm_util "gorm.io/gorm/utils"
)

// idea: https://901abc26123bee.github.io/how-to-implement-custom-gorm-logger/

const UserUIDKey = "user_uid"
const slowThreshold = 200 * time.Millisecond

/*
// set UserUIDKey for custom gorm logging
ctx = context.WithValue(ctx, orm.UserUIDKey, uid)
*/

// CustomGormLogger defines custom gorm logger
type CustomGormLogger struct {
	ctx context.Context
}

// NewCustomGormLogger new a gorm custom logger instance
func NewCustomGormLogger(ctx context.Context) *CustomGormLogger {
	return &CustomGormLogger{
		ctx: ctx,
	}
}

// LogMode log mode
func (l *CustomGormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Printf(l.ctx, msg, data...)
}

// Warn print warn messages
func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warnf(ctx, msg, data...)
}

// Error print error messages
func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Error(ctx, fmt.Errorf(msg, data...))
}

// Trace print custom sql message with uid
func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	userUid := ctx.Value(UserUIDKey)
	elapsed := float64(time.Since(begin)) / 1e6

	userStr := ""
	if userUid != nil {
		userStr = fmt.Sprintf("[user_uid %s] ", userUid)
	}

	sql, rows := fc()
	formatSql := strings.ReplaceAll(sql, "\"", "")

	log.Printf(l.ctx, "[file]=%s", gorm_util.FileWithLineNum())
	log.Printf(l.ctx, "[%.3fms][rows %d]%s%s", elapsed, rows, userStr, formatSql)

	if err != nil {
		log.Warnf(l.ctx, "[%.3fms] SQL ERROR, %v", elapsed, err)
	}
	if elapsed > float64(slowThreshold.Nanoseconds()) {
		log.Printf(l.ctx, "[%.3fms] SLOW SQL", elapsed)
	}
}
