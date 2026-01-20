package dbu

import (
	"log/slog"
	"os"
	"testing"

	"github.com/byte4cat/nbx/v2/pkg/tlog"
)

func TestMain(m *testing.M) {
	logger := tlog.New(&tlog.Config{
		StderrLevel: slog.LevelDebug,
		FileLevel:   slog.LevelError,
		LogFilePath: "./test.log",
	})

	slog.SetDefault(logger)

	os.Exit(m.Run())
}
