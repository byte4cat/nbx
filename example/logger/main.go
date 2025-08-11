package main

import (
	"github.com/byte4cat/nbx/pkg/logger"
	"github.com/byte4cat/nbx/pkg/logger/adapter"
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
	// logger.Fatalf("This is a fatal message: %s", "Critical error!")
	// logger.DPanicf("This is a DPanic message: %s", "Panic in development mode!")
	// logger.Panicf("This is a panic message: %s", "System failure!")

	// Example of using the adapter
	interceptor := adapter.NewInterceptorLogger(log)
	_ = interceptor // do something useful
}
