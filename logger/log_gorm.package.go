package logger

import (
	"context"
	"errors"
	"strings"
	"time"

	util "snap_chat_server/utils"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogInterface = gormlogger.Interface

type gormLogger struct {
	SlowThreshold             time.Duration
	LogLevel                  gormlogger.LogLevel
	SourceField               string
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	Colorful                  bool
	SkipErrRecordNotFound     bool
	Logger                    zerolog.Logger
}

func NewGormPackageLogger() *gormLogger {
	return &gormLogger{
		LogLevel:                  gormlogger.Info,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
		Logger:                    NewConsole().Instance,
	}
}

func (l *gormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Info().Msgf(s, args...)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Warn().Msgf(s, args...)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Error().Msgf(s, args...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	fields := map[string]interface{}{
		"sql":      util.WordLimiter(strings.ReplaceAll(strings.ReplaceAll(sql, "\n", ""), "\t", ""), 100),
		"duration": elapsed,
		"rows":     rowsAffected,
	}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		l.Logger.Error().Err(err).Fields(fields).Msg("[GORM] query error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.Warn().Fields(fields).Msgf("[GORM] slow query")
		return
	}

	l.Logger.Debug().Fields(fields).Msgf("[GORM] query")
}
