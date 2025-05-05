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
	logger.Infof("This is an info message: %s", "Hello, World!")

	// Example of using the adapter
	interceptor := adapter.NewInterceptorLogger(log)
	_ = interceptor // do something useful
}
