package dbu

import (
	"fmt"
	"os"
	"testing"

	"github.com/byte4cat/nbx/v2/pkg/logger"
)

func TestMain(m *testing.M) {
	// setup
	_, err := logger.New(logger.DefaultConfig(), 3)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}

	os.Exit(m.Run())
}
