package logger

import (
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

// Init инициализирует глобальный логгер
func Init(logLevel string) {
	lvl := slog.LevelInfo
	switch logLevel {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	level.Store(lvl)

	// Обработчик для логирования в файл
	fileHandler := slog.NewJSONHandler(&lumberjack.Logger{
		Filename:   logsFilePath,
		MaxSize:    fileMegabytesMaxSize,
		MaxBackups: fileMaxBackupsCount,
		MaxAge:     fileMaxAge,
	}, &slog.HandlerOptions{Level: lvl})

	// Обработчик для логирования в консоль
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})

	// Используем кастомный multiHandler для логирования в оба источника
	globalLogger = slog.New(&multiHandler{handlers: []slog.Handler{fileHandler, consoleHandler}})
}

// Debug обертка над методом
func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

// Info обертка над методом
func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

// Warn обертка над методом
func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

// Error обертка над методом
func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}

// With обертка над методом
func With(args ...any) *slog.Logger {
	return globalLogger.With(args...)
}

// An Attr is a key-value pair.
type Attr struct {
	Key   string
	Value string
}

// String обертка над методом
func String(key, value string) Attr {
	return Attr{key, value}
}
