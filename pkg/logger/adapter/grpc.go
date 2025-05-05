package adapter

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func NewInterceptorLogger(zl *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]
			switch v := value.(type) {
			case string:
				f = append(f, zap.String(fmt.Sprint(key), v))
			case int:
				f = append(f, zap.Int(fmt.Sprint(key), v))
			case bool:
				f = append(f, zap.Bool(fmt.Sprint(key), v))
			default:
				f = append(f, zap.Any(fmt.Sprint(key), v))
			}
		}
		logger := zl.WithOptions(zap.AddCallerSkip(1)).With(f...)
		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
