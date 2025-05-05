package logger

import (
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getWriterSyncer(logPath string) zapcore.WriteSyncer {
	if logPath == "" {
		// if logPath is empty, write to stdout
		return zapcore.AddSync(os.Stdout)
	}

	// if logPath is provided, use lumberjack for file rotation
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    100, // MB
			MaxAge:     30,  // days
			MaxBackups: 3,
			LocalTime:  false,
			Compress:   false,
		}),
		zapcore.AddSync(os.Stdout),
	)
}
