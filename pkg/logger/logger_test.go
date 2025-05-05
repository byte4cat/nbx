package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewLogger_SingletonBehavior(t *testing.T) {
	t.Run("should initialize logger only once", func(t *testing.T) {
		logger1, err1 := New(Config{
			Mode:        Mode_Development.String(),
			LogFilePath: "./test1.log",
		})
		require.NoError(t, err1)
		assert.NotNil(t, logger1)

		logger2, err2 := New(Config{
			Mode:        Mode_Development.String(),
			LogFilePath: "./test2.log",
		})
		require.NoError(t, err2)
		assert.NotNil(t, logger2)

		// Sould be the same instance
		assert.Same(t, logger1, logger2, "logger should only be initialized once")

		// Check if the global logger is the same as the initialized logger
		// zap.ReplaceGlobals(logger1) used in New function
		assert.Same(t, logger1, zap.L(), "global logger should match singleton")
	})

	t.Run("log file path should be empty", func(t *testing.T) {
		logger, err := New(Config{
			Mode:        Mode_Development.String(),
			LogFilePath: "",
		})
		require.NoError(t, err)
		assert.NotNil(t, logger)

		// Check if the global logger is the same as the initialized logger
		assert.Same(t, logger, zap.L(), "global logger should match singleton")
	})
}
