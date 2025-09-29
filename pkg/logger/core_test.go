package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewLogger_SingletonBehavior(t *testing.T) {
	t.Run("log file path should be empty", func(t *testing.T) {
		logger, err := New(Config{
			Mode:        Mode_Development.String(),
			LogFilePath: "",
		}, 0)
		require.NoError(t, err)
		assert.NotNil(t, logger)

		logger.Debug("This is a debug message")

		// Check if the global logger is the same as the initialized logger
		assert.NotSame(t, logger, zap.L(), "global logger should not match singleton")
	})
}
