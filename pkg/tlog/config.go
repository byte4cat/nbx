package tlog

import "log/slog"

type Config struct {
	StderrLevel slog.Leveler // default: Info
	FileLevel   slog.Leveler // default: Error
	LogFilePath string       // file path, not writing to file if its empty string
	ForceJSON   bool
	NoColor     bool
	TimeFormat  string
}
