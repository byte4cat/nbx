package logger

import (
	"sync"

	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerOnce     sync.Once
	logger         *zap.Logger
	isLevelEnabled = make(map[zapcore.Level]bool)
	baseLogger     *zap.Logger
)

// New initializes a zap logger instance with the provided configuration and caller skip.
// This function is safe to call multiple times but the core logger setup happens only once.
// It returns the configured logger instance.
//
// Parameters:
//    - cfg: Config structure.
//    - callerSkip: The number of stack frames to skip to find the original caller.
//
// Returns:
//    - *zap.Logger: the configured zap logger instance.
//    - error: any error encountered during logger setup.
//
// Note: This function does NOT automatically set this logger as the global zap logger.
// You must call zap.ReplaceGlobals() externally if needed.
func New(cfg Config, callerSkip int) (*zap.Logger, error) {
	var core zapcore.Core

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

		core = zapcore.NewCore(encoder, writeSyncer, logLevel)

		baseLogger = zap.New(core)
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip))
	})

	return logger, nil
}

func GetLoggerWithoutCaller(cfg Config) (*zap.Logger, error) {
	if baseLogger == nil {
		return nil, errors.New("base logger not initialized")
	}
	return baseLogger, nil
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API.
// Sugaring a Logger is quite inexpensive, so it's reasonable for a single
// application to use both Loggers and SugaredLoggers, converting between them
// on the boundaries of performance-sensitive code.
func Sugar() *zap.SugaredLogger {
	if logger == nil {
		// Handle case where logger wasn't initialized, maybe return a no-op logger's sugar
		return zap.NewNop().Sugar()
	}
	return logger.Sugar()
}

// IsLevelEnabled checks if a given level is enabled for the main application logger.
// Ensure mainAppLogger is initialized before calling.
func IsLevelEnabled(level zapcore.Level) bool {
	if logger == nil {
		// Handle case where logger wasn't initialized
		return false
	}

	// return false // Level not in map? Should not happen for standard levels
	return logger.Core().Enabled(level)
}

// InitGlobalLogger initializes the global zap logger with the provided
// configuration and caller skip.
func InitGlobalLogger(cfg Config, callerSkip int) (*zap.Logger, error) {
	l, err := New(cfg, callerSkip)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(l)

	return l, nil
}
