package logkit

import "strings"

// validationError accumulates multiple validation failures (invalid types or values) during config parsing.
// Used internally by WithConfig to report all issues at once.
type validationError struct {
	invalidTypes  []string
	invalidValues []string
}

// Error returns a string representation of validation error.
func (e *validationError) Error() string {
	var b strings.Builder
	if len(e.invalidTypes) > 0 {
		b.WriteString("invalid_type=" + strings.Join(e.invalidTypes, ","))
	}
	if len(e.invalidValues) > 0 {
		if b.Len() > 0 {
			b.WriteString(", ")
		}
		b.WriteString("invalid_value=" + strings.Join(e.invalidValues, ","))
	}
	return b.String()
}

// hasErrors returns true if there are any validation errors.
func (e *validationError) hasErrors() bool {
	return len(e.invalidTypes) > 0 || len(e.invalidValues) > 0
}
