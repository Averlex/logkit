// Package logkit provides a structured logging wrapper over slog.Logger
// with enhanced context handling, custom levels, and flexible configuration.
package logkit

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Option defines a function that allows to configure underlying logger on construction.
type Option func(c *Config) error

// Config defines an inner logger configuration.
type Config struct {
	handlerOpts    *slog.HandlerOptions
	logType        string
	handler        slog.Handler
	writer         io.Writer
	timeTemplate   string
	level          slog.Level
	setupLevel     bool
	extraCtxFields []any
}

// WithConfig allows to apply custom configuration.
// Expected following config structure:
//
//	{
//			format        string, // "text" or "json"
//			level         string, // "debug", "info", "warn", "error"
//			time_template string, // any valid time format
//			log_stream:   string, // "stdout", "stderr"
//	}
func WithConfig(cfg map[string]any) Option {
	return func(c *Config) error {
		optionalFields := map[string]any{
			"format":        "",
			"level":         "",
			"time_template": "",
			"log_stream":    "",
		}

		ve := &validationError{}
		ve.invalidTypes = validateTypes(cfg, optionalFields)

		validateLogLevel(cfg, ve)
		validateTimeFormat(cfg, ve)
		validateWriter(cfg, ve)
		validateLogType(cfg, ve)

		if ve.hasErrors() {
			return fmt.Errorf("config data is invalid: %s", ve.Error())
		}

		if level, ok := cfg["level"]; ok {
			levelStr := strings.ToLower(level.(string))
			if level, ok := levelValues[levelStr]; ok {
				c.level = level
			} else {
				c.setupLevel = true
			}
		}

		if timeTmpl, ok := cfg["time_template"]; ok {
			c.timeTemplate = timeTmpl.(string)
		}

		if writer, ok := cfg["log_stream"]; ok {
			switch strings.ToLower(writer.(string)) {
			case "stdout":
				c.writer = os.Stdout
			case "stderr":
				c.writer = os.Stderr
			}
		}

		if logType, ok := cfg["format"]; ok {
			c.logType = logType.(string)
		}

		c.checkDefaults()
		c.handler = buildHandler(c)

		return nil
	}
}

// WithWriter allows to apply custom configuration.
func WithWriter(w io.Writer) Option {
	return func(c *Config) error {
		if w == nil {
			return fmt.Errorf("expected io.Writer, got nil")
		}

		c.writer = w
		c.handler = buildHandler(c)

		return nil
	}
}

// WithDefaults applies default configuration to the logger.
// May be overwritten by WithConfig and/or WithWriter options.
func WithDefaults() Option {
	return WithConfig(map[string]any{
		"format":        DefaultLogType,
		"level":         DefaultLevel,
		"time_template": DefaultTimeTemplate,
		"log_stream":    DefaultWriter,
	})
}

// WithExtraContextFields configures the logger to automatically include
// values from the context in every log record, using the provided keys.
//
// The keys must be comparable (as required by context.Context) and are used
// directly in ctx.Value(key) to retrieve the corresponding values.
// Compatible key types are:
//   - string: used directly as the attribute key in logs.
//   - Any type implementing fmt.Stringer: its String() method is used as the attribute key.
//
// This includes named string types (e.g., type MyString string) that implement fmt.Stringer.
// The logger will use these keys to extract values from the context and add them as log attributes:
//   - If a value is of type slog.Attr, it is added to the log as is - using internal key in slog.Attr pair.
//   - Otherwise, it is added as a regular attribute with the key as described above.
//
// If no key is found in the context, it is safely ignored by the logger.
//
// Example:
//
//	var RequestIDKey struct{}
//	func (RequestIDKey) String() string { return "request_id" }
//
//	logger := NewLogger(WithExtraContextFields(RequestIDKey))
//	ctx := context.WithValue(ctx, RequestIDKey, "abc-123")
//	logger.Info(ctx, "event") // → request_id=abc-123
//
// If no fields are provided, the option does nothing and returns nil.
// If any key types are invalid, a single error listing all erroneous types is returned.
func WithExtraContextFields(fields ...any) Option {
	return func(c *Config) error {
		if len(fields) == 0 {
			return nil
		}

		if err := validateLoggableContextKeys(fields...); err != nil {
			return err
		}

		c.extraCtxFields = append(c.extraCtxFields, fields...)
		return nil
	}
}

// NewLogger returns a new Logger with the given log type and level.
// If no opts are provided, it returns a default logger.
//
// The log type can be "text" or "json". The log level can be "debug", "info", "warn" or "error".
//
// timeTemplate is a time format string. Any format which is valid for time.Time format is acceptable.
//
// Empty log level corresponds to "error", as well as empty log type corresponds to "json".
// Empty time format is equal to the default value which is "02.01.2006 15:04:05.000".
// Empty writer option equals to using os.Stdout. Custom writer might be set using WithWriter option.
//
// If the log type or level is unknown, it returns an error.
func NewLogger(opts ...Option) (*Logger, error) {
	cfg := &Config{}

	if err := WithDefaults()(cfg); err != nil {
		return nil, errors.New("default logger initialization failed")
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("logger initialization failed: %w", err)
		}
	}

	return &Logger{slog.New(cfg.handler), cfg.extraCtxFields}, nil
}
