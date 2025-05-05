package logger

// Config defines the configuration for the logger.
//
// It controls the logger’s operating mode, output destinations, and log verbosity.
// When no file path is specified, logs are printed to stdout only.
type Config struct {
	// Mode sets the environment mode for the logger.
	// Accepted values are "development" or "production".
	// Defaults to "development" if unspecified.
	Mode string

	// LogFilePath is the path to the log file used for file logging with rotation.
	//
	// If provided, logs will be written both to this file and stdout.
	// If empty, only stdout will be used.
	LogFilePath string

	// LogLevel sets the verbosity level of the logger.
	//
	// Valid values are: "debug", "info", "warn", "error", "fatal" "panic",
	// and "dpanic".
	// Defaults to "info" if not set.
	LogLevel *LogLevel
}

func DefaultConfig() Config {
	return Config{
		Mode:        "development",
		LogFilePath: "",
		LogLevel:    nil,
	}
}
