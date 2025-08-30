package logkit

import (
	"log/slog"
	"os"
)

const (
	// DefaultTimeTemplate is a default time format string.
	DefaultTimeTemplate = "02.01.2006 15:04:05.000"
	// DefaultLogType is a default log type.
	DefaultLogType = "json"
	// DefaultLevel is a default log level.
	DefaultLevel = "error"
	// DefaultLevelValue is a default log level value.
	DefaultLevelValue = slog.LevelError
	// DefaultWriter is a default writer to use for logging.
	DefaultWriter = "stdout"
)

// DefaultWriterValue is a default writer value.
var DefaultWriterValue = os.Stdout

// Additional log levels.
const (
	// LevelTrace is a trace log level - the lowest possible level.
	LevelTrace = slog.Level(-8)
	// LevelDebug is an alias for slog.LevelDebug.
	LevelDebug = slog.LevelDebug
	// LevelVerbose is a verbose log level - the middleground between Debug and Info levels.
	LevelVerbose = slog.Level(-2)
	// LevelInfo is an alias for slog.LevelInfo.
	LevelInfo = slog.LevelInfo
	// LevelWarn is an alias for slog.LevelWarn.
	LevelWarn = slog.LevelWarn
	// LevelError is an alias for slog.LevelError.
	LevelError = slog.LevelError
	// LevelFatal is a fatal log level - the highest possible level.
	LevelFatal = slog.Level(16)
)

// A helper to map log levels to their names.
var levelNames = map[slog.Level]string{
	LevelTrace:      "TRACE",
	slog.LevelDebug: "DEBUG",
	LevelVerbose:    "VERBOSE",
	slog.LevelInfo:  "INFO",
	slog.LevelWarn:  "WARN",
	slog.LevelError: "ERROR",
	LevelFatal:      "FATAL",
}

// A helper to map log levels received in configuration to their values.
var levelValues = map[string]slog.Level{
	"trace":   LevelTrace,
	"debug":   LevelDebug,
	"verbose": LevelVerbose,
	"info":    LevelInfo,
	"warn":    LevelWarn,
	"error":   LevelError,
	"fatal":   LevelFatal,
}
