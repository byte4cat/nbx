package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Infof logs a message at Info level using fmt.Sprintf-style formatting.
func Infof(msg string, args ...any) {
	logger.Info(msg, toZapFields(args)...)
}

// Info logs a message at Info level with structured fields.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Debugf logs a message at Debug level using fmt.Sprintf-style formatting.
func Debugf(msg string, args ...any) {
	logger.Debug(msg, toZapFields(args)...)
}

// Debug logs a message at Debug level with structured fields.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warnf logs a message at Warn level using fmt.Sprintf-style formatting.
func Warnf(msg string, args ...any) {
	logger.Warn(msg, toZapFields(args)...)
}

// Warn logs a message at Warn level with structured fields.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Errorf logs a message at Error level using fmt.Sprintf-style formatting.
func Errorf(msg string, args ...any) {
	logger.Error(msg, toZapFields(args)...)
}

// Error logs a message at Error level with structured fields.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatalf logs a message at Fatal level using fmt.Sprintf-style formatting.
// The application will terminate immediately.
func Fatalf(msg string, args ...any) {
	logger.Fatal(msg, toZapFields(args)...)
}

// Fatal logs a message at Fatal level with structured fields.
// The application will terminate immediately.
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Panicf logs a message at Panic level using fmt.Sprintf-style formatting.
// It then panics.
func Panicf(msg string, args ...any) {
	logger.Panic(msg, toZapFields(args)...)
}

// Panic logs a message at Panic level with structured fields.
// It then panics.
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// DPanicLevel logs are particularly important errors.
// In development the logger panics after writing the message.
func DPanicf(msg string, args ...any) {
	logger.DPanic(msg, toZapFields(args)...)
}

func toZapFields(args []any) []zap.Field {
	fields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("invalid_key_%d", i)
		}
		fields = append(fields, zap.Any(key, args[i+1]))
	}
	if len(args)%2 != 0 {
		fields = append(fields, zap.Any("invalid_last_key", args[len(args)-1]))
	}
	return fields
}
