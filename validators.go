package logkit

import (
	"reflect"
	"strings"
	"time"
)

// validateLogLevel is a helper that checks if log level is valid.
func validateLogLevel(cfg map[string]any, ve *validationError) {
	if val, ok := cfg["level"]; ok {
		levelStr, ok := val.(string)
		if !ok {
			ve.invalidTypes = append(ve.invalidTypes, "level")
			return
		}
		levelStr = strings.ToLower(levelStr)

		if _, ok := levelValues[levelStr]; !ok {
			ve.invalidValues = append(ve.invalidValues, "level")
			return
		}
	}
}

// validateTimeFormat is a helper that checks if time format is valid.
func validateTimeFormat(cfg map[string]any, ve *validationError) {
	if val, ok := cfg["time_template"]; ok {
		timeTmpl, ok := val.(string)
		if !ok {
			ve.invalidTypes = append(ve.invalidTypes, "time_template")
			return
		}

		if timeTmpl == "" {
			return
		}

		testTime := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
		formatted := testTime.Format(timeTmpl)
		parsedTime, err := time.Parse(timeTmpl, formatted)
		if err != nil || !parsedTime.Equal(testTime) {
			ve.invalidValues = append(ve.invalidValues, "time_template")
		}
	}
}

// validateWriter is a helper that checks if writer is valid.
func validateWriter(cfg map[string]any, ve *validationError) {
	if val, ok := cfg["log_stream"]; ok {
		writerStr, ok := val.(string)
		if !ok {
			ve.invalidTypes = append(ve.invalidTypes, "log_stream")
			return
		}

		switch writerStr {
		case "stdout", "stderr", "":
		default:
			ve.invalidValues = append(ve.invalidValues, "log_stream")
		}
	}
}

// validateLogType is a helper that checks if log type is valid.
func validateLogType(cfg map[string]any, ve *validationError) {
	if val, ok := cfg["format"]; ok {
		logTypeStr, ok := val.(string)
		if !ok {
			ve.invalidTypes = append(ve.invalidTypes, "format")
			return
		}

		switch logTypeStr {
		case "json", "text", "":
		default:
			ve.invalidValues = append(ve.invalidValues, "format")
		}
	}
}

// validateTypes returns missing and wrong type fields found in args.
// optionalFields is a map of field names with their expected types.
func validateTypes(args map[string]any, optionalFields map[string]any) (invalidTypes []string) {
	for field, expectedVal := range optionalFields {
		val, exists := args[field]
		if !exists {
			continue
		}

		// Default type switch will end up with false positive results. E.g., 123.(string) -> ok.
		// Using soft type check for string types, as, i.e., timeFormat in time package is untyped string.
		expectedKind := reflect.TypeOf(expectedVal).Kind()
		actualKind := reflect.TypeOf(val).Kind()
		if expectedKind != actualKind {
			invalidTypes = append(invalidTypes, field)
		}
	}

	return invalidTypes
}

// checkDefaults sets default values for the logger configurations, if they are empty.
func (c *Config) checkDefaults() {
	if c.logType == "" {
		c.logType = DefaultLogType
	}
	if c.writer == nil {
		c.writer = DefaultWriterValue
	}
	if c.timeTemplate == "" {
		c.timeTemplate = DefaultTimeTemplate
	}
	if c.setupLevel {
		c.level = DefaultLevelValue
	}
}
