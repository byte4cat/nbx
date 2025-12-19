package main

import "github.com/byte4cat/nbx/v2/cmd"

func main() {
	cmd.Execute()
	// logger := tlog.New(tlog.Config{
	// 	Level:       slog.LevelDebug,
	// 	LogFilePath: "./test_run.log",
	// })
	//
	// slog.SetDefault(logger)
	//
	// slog.Info("info test", "hello", "world")
	// slog.Warn("warn test", "hello", "world")
	// slog.Error("error test", "hello", "world")
	// slog.Debug("debug test", "hello", "world")
}
