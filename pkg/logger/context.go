package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// ContextKey типы для хранения значений в контексте
type ContextKey string

const (
	RequestIDKey   ContextKey = "request_id"
	UserIDKey      ContextKey = "user_id"
	OrganizationID ContextKey = "org_id"
	TraceIDKey     ContextKey = "trace_id"
)

// WithContext добавляет значение в контекст
func WithContext(ctx context.Context, key ContextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// FromContext извлекает значение из контекста
func FromContext(ctx context.Context, key ContextKey) (interface{}, bool) {
	return ctx.Value(key), ctx.Value(key) != nil
}

// Logger обертка над logrus.Logger с дополнительным функционалом
type Logger struct {
	logger *logrus.Logger
}

// New создает новый Logger
func New(log *logrus.Logger) *Logger {
	return &Logger{logger: log}
}

// withCaller добавляет информацию о файле и строке вызова
func (l *Logger) withCaller() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return l.logger.WithField("caller", "unknown")
	}

	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		parts := strings.Split(fn.Name(), ".")
		funcName = parts[len(parts)-1]
	}

	// Сокращаем путь файла до последних двух сегментов
	fileParts := strings.Split(file, "/")
	if len(fileParts) > 2 {
		file = strings.Join(fileParts[len(fileParts)-2:], "/")
	}

	return l.logger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
		"func": funcName,
	})
}

// Info логирует информационное сообщение с контекстом
func (l *Logger) Info(ctx context.Context, message string, fields ...logrus.Fields) {
	entry := l.withCaller()
	entry = l.addContextFields(entry, ctx)

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Info(message)
}

// Debug логирует отладочное сообщение с контекстом
func (l *Logger) Debug(ctx context.Context, message string, fields ...logrus.Fields) {
	entry := l.withCaller()
	entry = l.addContextFields(entry, ctx)

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Debug(message)
}

// Warn логирует предупреждение с контекстом
func (l *Logger) Warn(ctx context.Context, message string, fields ...logrus.Fields) {
	entry := l.withCaller()
	entry = l.addContextFields(entry, ctx)

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Warn(message)
}

// Error логирует ошибку с контекстом и stack trace
func (l *Logger) Error(ctx context.Context, message string, err error, fields ...logrus.Fields) {
	entry := l.withCaller()
	entry = l.addContextFields(entry, ctx)

	if err != nil {
		entry = entry.WithError(err)
		// Добавляем stack trace для критических ошибок
		entry = entry.WithField("stack", l.getStackTrace())
	}

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Error(message)
}

// Fatal логирует критическую ошибку и выходит
func (l *Logger) Fatal(ctx context.Context, message string, err error, fields ...logrus.Fields) {
	entry := l.withCaller()
	entry = l.addContextFields(entry, ctx)

	if err != nil {
		entry = entry.WithError(err)
		entry = entry.WithField("stack", l.getStackTrace())
	}

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Fatal(message)
}

// addContextFields добавляет поля из контекста в entry
func (l *Logger) addContextFields(entry *logrus.Entry, ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if requestID, ok := FromContext(ctx, RequestIDKey); ok {
		fields["request_id"] = requestID
	}

	if userID, ok := FromContext(ctx, UserIDKey); ok {
		fields["user_id"] = userID
	}

	if orgID, ok := FromContext(ctx, OrganizationID); ok {
		fields["org_id"] = orgID
	}

	if traceID, ok := FromContext(ctx, TraceIDKey); ok {
		fields["trace_id"] = traceID
	}

	return entry.WithFields(fields)
}

// getStackTrace получает stack trace для логирования
func (l *Logger) getStackTrace() string {
	const maxDepth = 5
	var buf strings.Builder

	for i := 2; i < maxDepth+2; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		funcName := "unknown"
		if fn != nil {
			parts := strings.Split(fn.Name(), ".")
			funcName = parts[len(parts)-1]
		}

		// Сокращаем путь
		fileParts := strings.Split(file, "/")
		if len(fileParts) > 1 {
			file = strings.Join(fileParts[len(fileParts)-1:], "/")
		}

		buf.WriteString(fmt.Sprintf("%s:%d in %s\n", file, line, funcName))
	}

	return buf.String()
}

// WithField добавляет поле к логгеру (для цепочки вызовов)
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.withCaller().WithField(key, value)
}

// WithFields добавляет поля к логгеру (для цепочки вызовов)
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.withCaller().WithFields(fields)
}

// Raw возвращает базовый logrus.Logger
func (l *Logger) Raw() *logrus.Logger {
	return l.logger
}
