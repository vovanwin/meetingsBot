package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct {
	l *slog.Logger
}

func NewLoggerTracer(l *slog.Logger) *Logger {
	return &Logger{l: l}
}

// Log метод для трассировки запросов pgx с учетом уровня логирования и дополнительной информации.
func (l *Logger) Log(
	ctx context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]interface{},
) {
	// Добавляем информацию о месте вызова
	_, file, line, _ := runtime.Caller(6)
	location := slog.String("source", fmt.Sprintf("%s:%d", file, line))

	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}
	attrs = append(attrs, location) // Добавляем источник вызова

	var lvl slog.Level
	switch level {
	case tracelog.LogLevelTrace:
		lvl = slog.LevelDebug - 1
	case tracelog.LogLevelDebug:
		lvl = slog.LevelDebug
	case tracelog.LogLevelInfo:
		lvl = slog.LevelInfo
	case tracelog.LogLevelWarn:
		lvl = slog.LevelWarn
	case tracelog.LogLevelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelError
		attrs = append(attrs, slog.Any("INVALID_PGX_LOG_LEVEL", level))
	}

	// Логирование с использованием переданного контекста
	l.l.LogAttrs(ctx, lvl, msg, attrs...)
}
