package main

import (
	"log/slog"

	"github.com/byte4cat/nbx/v2/pkg/tlog"
)

func main() {
	logger := tlog.New(tlog.Config{
		StderrLevel: slog.LevelDebug,
		FileLevel:   slog.LevelError,
		LogFilePath: "./test.log",
	})

	slog.SetDefault(logger)

	slog.Info("info test", "hello", "world")
	slog.Warn("warn test", "hello", "world")
	slog.Error("error test", "hello", "world")
	slog.Debug("debug test", "hello", "world")
}
