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

var contextRequestKeys = []string{"http_request_id", "grpc_request_id"}

// DefaultWriterValue is a default writer value.
var DefaultWriterValue = os.Stdout
