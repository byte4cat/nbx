package tlog

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg Config) *slog.Logger {
	var handlers []slog.Handler

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

	isTTY := isatty.IsTerminal(os.Stderr.Fd()) ||
		isatty.IsCygwinTerminal(os.Stderr.Fd())

	var stderrHandler slog.Handler
	if isTTY && !cfg.NoColor && !cfg.ForceJSON {
		stderrHandler = tint.NewHandler(os.Stderr, &tint.Options{
			Level:      stderrLevel,
			TimeFormat: timeFormat,
			AddSource:  true,
		})
	} else {
		stderrHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     stderrLevel,
			AddSource: true,
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
			Level:     fileLevel,
			AddSource: true,
		})
		handlers = append(handlers, fileHandler)
	}

	return slog.New(&MultiHandler{handlers: handlers})
}
