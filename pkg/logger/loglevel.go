package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

// A Level is a logging priority. Higher levels are more important.
type LogLevel int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	LogLevel_Debug LogLevel = -1
	// InfoLevel logs are the default logging level.
	LogLevel_Info LogLevel = 0
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	LogLevel_Warn LogLevel = 1
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	LogLevel_Error LogLevel = 2
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	LogLevel_DPanic LogLevel = 3
	// PanicLevel logs a message, then panics.
	LogLevel_Panic LogLevel = 4
	// FatalLevel logs a message, then calls os.Exit(1).
	LogLevel_Fatal LogLevel = 5
)

var loglevelToString = map[LogLevel]string{
	LogLevel_Debug:  "debug",
	LogLevel_Info:   "infolevel",
	LogLevel_Warn:   "warnlevel",
	LogLevel_Error:  "errorlevel",
	LogLevel_DPanic: "dpaniclevel",
	LogLevel_Panic:  "paniclevel",
	LogLevel_Fatal:  "fatallevel",
}

var stringToLogLevel = map[string]LogLevel{
	"debug":       LogLevel_Debug,
	"infolevel":   LogLevel_Info,
	"warnlevel":   LogLevel_Warn,
	"errorlevel":  LogLevel_Error,
	"dpaniclevel": LogLevel_DPanic,
	"paniclevel":  LogLevel_Panic,
	"fatallevel":  LogLevel_Fatal,
}

func (e LogLevel) String() string {
	return loglevelToString[e]
}

func (e LogLevel) IsValid() bool {
	_, ok := loglevelToString[e]
	return ok
}

func ParseLogLevelString(s string) (LogLevel, error) {
	if val, ok := stringToLogLevel[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("invalid LogLevel: %s", s)
}

func (l LogLevel) ToZapLevel() zapcore.Level {
	return zapcore.Level(l)
}

func LogLevelOptions() []LogLevel {
	return []LogLevel{
		LogLevel_Debug,
		LogLevel_Info,
		LogLevel_Warn,
		LogLevel_Error,
		LogLevel_DPanic,
		LogLevel_Panic,
		LogLevel_Fatal,
	}
}
