package main

import (
	"errors"

	"github.com/byte4cat/nbx/v2/pkg/logger"
	"github.com/byte4cat/nbx/v2/pkg/logger/adapter"
)

func main() {
	cfg := logger.Config{
		Mode:        "development",
		LogFilePath: "./logs/app.log",
	}

	log, _ := logger.New(cfg, 3)

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	logger.Debug("Debug with fields", logger.String("key1", "value1"))
	logger.Info("Info with fields", logger.Int("key2", 42))
	logger.Warn("Warn with fields", logger.Bool("key3", true))
	logger.Error("Error with fields", logger.Err(errors.New("sample error")))

	logger.Debugf("Formatted debug", "debug", "details")
	logger.Infof("Formatted info: %s = %d", "value", 100)
	logger.Warnf("Formatted warn: %s", "be careful")
	logger.Errorf("Formatted error: %v", errors.New("formatted error"))

	// Example of using the adapter
	interceptor := adapter.NewInterceptorLogger(log)
	_ = interceptor // do something useful
}
