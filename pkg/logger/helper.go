package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Infof logs a message at Info level using fmt.Sprintf-style formatting.
func Infof(msg string, args ...any) {
	logf(zapcore.InfoLevel, msg, args...)
}

// Info logs a message at Info level with structured fields.
func Info(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

// Debugf logs a message at Debug level using fmt.Sprintf-style formatting.
func Debugf(msg string, args ...any) {
	logf(zapcore.DebugLevel, msg, args...)
}

// Debug logs a message at Debug level with structured fields.
func Debug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

// Warnf logs a message at Warn level using fmt.Sprintf-style formatting.
func Warnf(msg string, args ...any) {
	logf(zapcore.WarnLevel, msg, args...)
}

// Warn logs a message at Warn level with structured fields.
func Warn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

// Errorf logs a message at Error level using fmt.Sprintf-style formatting.
func Errorf(msg string, args ...any) {
	logf(zapcore.ErrorLevel, msg, args...)
}

// Error logs a message at Error level with structured fields.
func Error(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}

// Fatalf logs a message at Fatal level using fmt.Sprintf-style formatting.
// The application will terminate immediately.
func Fatalf(msg string, args ...any) {
	logf(zapcore.FatalLevel, msg, args...)
}

// Fatal logs a message at Fatal level with structured fields.
// The application will terminate immediately.
func Fatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

// Panicf logs a message at Panic level using fmt.Sprintf-style formatting.
// It then panics.
func Panicf(msg string, args ...any) {
	logf(zapcore.PanicLevel, msg, args...)
}

// Panic logs a message at Panic level with structured fields.
// It then panics.
func Panic(msg string, fields ...zap.Field) {
	zapLogger.Panic(msg, fields...)
}

// DPanicLevel logs are particularly important errors.
// In development the logger panics after writing the message.
func DPanicf(msg string, args ...any) {
	logf(zapcore.DPanicLevel, msg, args...)
}

var levelLoggers = map[zapcore.Level]func(string){
	zapcore.DebugLevel:  func(msg string) { zapLogger.Debug(msg) },
	zapcore.InfoLevel:   func(msg string) { zapLogger.Info(msg) },
	zapcore.WarnLevel:   func(msg string) { zapLogger.Warn(msg) },
	zapcore.ErrorLevel:  func(msg string) { zapLogger.Error(msg) },
	zapcore.DPanicLevel: func(msg string) { zapLogger.DPanic(msg) },
	zapcore.PanicLevel:  func(msg string) { zapLogger.Panic(msg) },
	zapcore.FatalLevel:  func(msg string) { zapLogger.Fatal(msg) },
}

func logf(level zapcore.Level, msg string, args ...any) {
	if zapLogger == nil {
		panic("logger not initialized, call logger.New(config) first")
	}

	if logFunc, ok := levelLoggers[level]; ok {
		logFunc(fmt.Sprintf(msg, args...))
	}
}
