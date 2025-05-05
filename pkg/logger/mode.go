package logger

import (
	"fmt"
	"strings"
)

// Mode can be either "development" or "production"
type Mode string

const (
	// Development mode is used for development and testing purposes.
	Mode_Development Mode = "development"
	// Production mode is used for production environments.
	Mode_Production Mode = "production"
)

var modeToString = map[Mode]string{
	Mode_Development: "development",
	Mode_Production:  "production",
}

var stringToMode = map[string]Mode{
	"development": Mode_Development,
	"production":  Mode_Production,
}

func (e Mode) String() string {
	return modeToString[e]
}

func (e Mode) IsValid() bool {
	_, ok := modeToString[e]
	return ok
}

func ParseModeString(s string) (Mode, error) {
	if val, ok := stringToMode[strings.ToLower(s)]; ok {
		return val, nil
	}
	return "", fmt.Errorf("invalid Mode: %s", s)
}

func ModeOptions() []Mode {
	return []Mode{
		Mode_Development,
		Mode_Production,
	}
}
