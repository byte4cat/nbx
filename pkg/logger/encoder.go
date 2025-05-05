package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getProductionEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(cfg)
}

func getDevelopmentEncoder() zapcore.Encoder {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(cfg)
}
