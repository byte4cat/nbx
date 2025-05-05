package main

import (
	"github.com/yimincai/nbx/pkg/logger"
	"github.com/yimincai/nbx/pkg/logger/adapter"
)

func main() {
	cfg := logger.Config{
		Mode:        "development",
		LogFilePath: "./logs/app.log",
	}

	log, _ := logger.New(cfg)

	logger.Debugf("This is a debug message: %s", "Debugging info")
	logger.Infof("This is an info message: %s", "Hello, World!")
	logger.Warnf("This is a warning message: %s", "Be careful!")
	logger.Errorf("This is an error message: %s", "Something went wrong!")
	// logger.Fatalf("This is a fatal message: %s", "Critical error!")
	// logger.DPanicf("This is a DPanic message: %s", "Panic in development mode!")
	// logger.Panicf("This is a panic message: %s", "System failure!")

	// Example of using the adapter
	interceptor := adapter.NewInterceptorLogger(log)
	_ = interceptor // do something useful
}
