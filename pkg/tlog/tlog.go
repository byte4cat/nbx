package tlog

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

var DefaultConfig = &Config{
	StderrLevel: slog.LevelInfo,
	FileLevel:   slog.LevelError,
	LogFilePath: "service.log",
	NoColor:     true,
	TimeFormat:  time.RFC3339,
	ForceText:   false,
	ForceJSON:   true,
}

func New(cfg *Config) *slog.Logger {
	var handlers []slog.Handler

	if cfg == nil {
		cfg = DefaultConfig
	}

	stderrLevel := cfg.StderrLevel
	if stderrLevel == nil {
		stderrLevel = slog.LevelInfo
	}

	fileLevel := cfg.FileLevel
	if fileLevel == nil {
		fileLevel = slog.LevelError
	}

	timeFormat := cfg.TimeFormat
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey && a.Value.Kind() == slog.KindAny {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				// 保留最後兩層 (例如 core/http.go)
				dir := filepath.Base(filepath.Dir(source.File))
				source.File = filepath.Join(dir, filepath.Base(source.File))
			}
		}
		return a
	}

	isTTY := isatty.IsTerminal(os.Stderr.Fd()) ||
		isatty.IsCygwinTerminal(os.Stderr.Fd())

	var stderrHandler slog.Handler
	if (isTTY || cfg.ForceText) && !cfg.NoColor && !cfg.ForceJSON {
		stderrHandler = tint.NewHandler(os.Stderr, &tint.Options{
			Level:       stderrLevel,
			TimeFormat:  timeFormat,
			AddSource:   true,
			ReplaceAttr: replaceAttr,
		})
	} else {
		stderrHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:       stderrLevel,
			AddSource:   true,
			ReplaceAttr: replaceAttr,
		})
	}

	handlers = append(handlers, stderrHandler)

	if cfg.LogFilePath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.LogFilePath,
			MaxSize:    100,  // max size 100MB
			MaxBackups: 5,    // keep 5 old files
			MaxAge:     30,   // keep it 30 days
			Compress:   true, // compress old files (.gz)
		}

		fileHandler := slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{
			Level:       fileLevel,
			AddSource:   true,
			ReplaceAttr: replaceAttr,
		})
		handlers = append(handlers, fileHandler)
	}

	return slog.New(&MultiHandler{handlers: handlers})
}
