package utils

import (
	"os/user"

	"github.com/byte4cat/nbx/pkg/clog/v2"
)

// GetUsernameOrPanic gets the current system username
func GetUsernameOrPanic() string {
	u, err := user.Current()
	if err != nil {
		clog.Panic("failed to get current user: %v", err)
	}

	return u.Username
}
