package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerOnce     sync.Once
	zapLogger      *zap.Logger
	isLevelEnabled = make(map[zapcore.Level]bool)
)

// New initializes the global zap logger instance using the provided configuration.
// This function is safe to call multiple times but only initializes the logger once (singleton).
//
// Parameters:
//   - cfg: Config structure that defines logging mode, level, and file output settings.
//
// Returns:
//   - *zap.Logger: the initialized zap logger instance.
//   - error: any error encountered during logger setup (currently always nil).
func New(cfg Config) (*zap.Logger, error) {
	var err error
	loggerOnce.Do(func() {
		writeSyncer := getWriterSyncer(cfg.LogFilePath)

		var encoder zapcore.Encoder
		var defaultLevel zapcore.Level

		switch cfg.Mode {
		case Mode_Production.String():
			encoder = getProductionEncoder()
			defaultLevel = zapcore.InfoLevel
		default:
			encoder = getDevelopmentEncoder()
			defaultLevel = zapcore.DebugLevel
		}

		logLevel := defaultLevel
		if cfg.LogLevel != nil {
			logLevel = cfg.LogLevel.ToZapLevel()
		}

		core := zapcore.NewCore(encoder, writeSyncer, logLevel)
		zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
		zap.ReplaceGlobals(zapLogger)

		for _, lvl := range []zapcore.Level{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			zapcore.WarnLevel,
			zapcore.ErrorLevel,
			zapcore.DPanicLevel,
			zapcore.PanicLevel,
			zapcore.FatalLevel,
		} {
			isLevelEnabled[lvl] = zapLogger.Core().Enabled(lvl)
		}
	})
	return zapLogger, err
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API.
// Sugaring a Logger is quite inexpensive, so it's reasonable for a single
// application to use both Loggers and SugaredLoggers, converting between them
// on the boundaries of performance-sensitive code.
func Sugar() *zap.SugaredLogger {
	return zapLogger.Sugar()
}

func IsInitialized() bool {
	return zapLogger != nil
}
