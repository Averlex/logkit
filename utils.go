package logkit

import (
	"log/slog"
	"strings"
	"time"
)

// buildHandler returns a handler based on config.
func buildHandler(c *Config) slog.Handler {
	c.handlerOpts = &slog.HandlerOptions{
		Level: c.level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			attr := replaceLevelAttr(groups, a)
			return replaceTimeAttrs(groups, attr, c.timeTemplate)
		},
	}

	switch strings.ToLower(c.logType) {
	case "json", "":
		return slog.NewJSONHandler(c.writer, c.handlerOpts)
	case "text":
		return slog.NewTextHandler(c.writer, c.handlerOpts)
	default:
		return slog.NewJSONHandler(c.writer, c.handlerOpts)
	}
}

// replaceTimeAttrs replaces time.Time values with formatted strings.
func replaceTimeAttrs(groups []string, a slog.Attr, timeFormat string) slog.Attr {
	// Default log timestamp.
	if a.Key == slog.TimeKey && len(groups) == 0 {
		if t, ok := a.Value.Any().(time.Time); ok {
			a.Value = slog.StringValue(t.Format(timeFormat))
		}
		return a
	}

	// Common time.Time fields.
	if v, ok := a.Value.Any().(time.Time); ok {
		a.Value = slog.StringValue(v.Format(timeFormat))
		return a
	}

	// Groups handling.
	if a.Value.Kind() == slog.KindGroup {
		newGroup := make([]slog.Attr, len(a.Value.Group()))
		for i, ga := range a.Value.Group() {
			newGroup[i] = replaceTimeAttrs(append(groups, a.Key), ga, timeFormat)
		}
		a.Value = slog.GroupValue(newGroup...)
	}

	return a
}

// replaceLevelAttr replaces slog.Level values with their names.
func replaceLevelAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey && len(groups) == 0 {
		switch level := a.Value.Any().(type) {
		// Common case: slog.Level value as a root-level string.
		case slog.Level:
			if name, exists := levelNames[level]; exists {
				return slog.String(slog.LevelKey, name)
			}
		// In case of receiving level as an int value.
		case int:
			if name, exists := levelNames[slog.Level(level)]; exists {
				return slog.String(slog.LevelKey, name)
			}
		// Unable to recognize the type.
		default:
		}
	}

	return a
}
