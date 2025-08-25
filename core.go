package logkit

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Logger is a wrapper structure for an underlying logger.
type Logger struct {
	l              *slog.Logger
	extraCtxFields []any
}

// addContextData extracts values from the context using keys defined via WithExtraContextFields.
// For each key:
//   - If the value in the context is of type slog.Attr, it is added to the log as is.
//   - Otherwise, the value is wrapped as an attribute with a string key derived from the context key:
//   - If the key is a string, it is used directly.
//   - If the key implements fmt.Stringer, its String() method is used as the attribute key.
//   - All other key types are ignored.
//
// This logic ensures compatibility with slog's requirement that attribute keys be strings,
// while allowing context keys to be any comparable type (e.g. custom structs) as long as
// they provide a string representation via fmt.Stringer or are plain strings.
func (logg Logger) addContextData(ctx context.Context, args ...any) []any {
	for _, k := range logg.extraCtxFields {
		v := ctx.Value(k)
		if v == nil {
			continue
		}
		// Logging slog.Attr as is.
		if attr, ok := v.(slog.Attr); ok {
			args = append(args, attr)
			continue
		}
		// Key should be slog-compatible: either string or fmt.Stringer.
		switch k := k.(type) {
		case string:
			args = append(args, slog.Any(k, v))
		case fmt.Stringer:
			args = append(args, slog.Any(k.String(), v))
		}
	}

	return args
}

// Trace logs a message with level Trace on the standard logger.
func (logg Logger) Trace(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelTrace, msg, logg.addContextData(ctx, args...)...)
}

// Debug logs a message with level Debug on the standard logger.
func (logg Logger) Debug(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelDebug, msg, logg.addContextData(ctx, args...)...)
}

// Verbose logs a message with level Verbose on the standard logger.
func (logg Logger) Verbose(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelVerbose, msg, logg.addContextData(ctx, args...)...)
}

// Info logs a message with level Info on the standard logger.
func (logg Logger) Info(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelInfo, msg, logg.addContextData(ctx, args...)...)
}

// Warn logs a message with level Warn on the standard logger.
func (logg Logger) Warn(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelWarn, msg, logg.addContextData(ctx, args...)...)
}

// Error logs a message with level Error on the standard logger.
func (logg Logger) Error(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelError, msg, logg.addContextData(ctx, args...)...)
}

// Fatal logs a message with level Error on the standard logger and then calls os.Exit(1).
func (logg Logger) Fatal(ctx context.Context, msg string, args ...any) {
	logg.l.Log(ctx, LevelFatal, msg, logg.addContextData(ctx, args...)...)
	os.Exit(1)
}

// With returns a new Logger that adds the given key-value pairs to the logger's context.
func (logg Logger) With(args ...any) *Logger {
	return &Logger{logg.l.With(args...), logg.extraCtxFields}
}
