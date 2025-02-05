package logger

import (
	"context"
	"log/slog"
	"sync/atomic"
)

// Константы для настроек логирования
const (
	logsFilePath         = "logs/app.log"
	fileMegabytesMaxSize = 10
	fileMaxBackupsCount  = 3
	fileMaxAge           = 7 // дней
)

// Глобальные переменные
var globalLogger *slog.Logger
var level atomic.Value

// multiHandler — кастомный обработчик для логирования в несколько мест
type multiHandler struct {
	handlers []slog.Handler
}

// Handle передает лог-сообщение всем обработчикам
func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// Enabled проверяет, можно ли логировать сообщение данного уровня
func (m *multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

// WithAttrs создает новый `multiHandler` с дополнительными атрибутами
func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

// WithGroup поддерживает группировку логов
func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}
