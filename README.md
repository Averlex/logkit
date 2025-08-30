# logkit

[![Go version](https://img.shields.io/badge/go-1.24.2+-blue.svg)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/Averlex/logkit.svg)](https://pkg.go.dev/github.com/Averlex/logkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/Averlex/logkit)](https://goreportcard.com/report/github.com/Averlex/logkit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A lightweight, configurable logging toolkit for Go's slog with context support and validation.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Custom Configuration](#custom-configuration)
- [Context-Aware Logging](#context-aware-logging)
- [Custom Log Levels](#custom-log-levels)
- [Advanced Usage](#advanced-usage)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Installation](#installation)

## Features

- ✅ **Custom log levels**: `TRACE`, `VERBOSE`, `FATAL` — beyond standard `slog`
- ✅ **Context-aware logging**: automatically inject values from `context.Context` using typed keys
- ✅ **Configurable output**: JSON or text format, custom time template, `stdout`/`stderr` or any custom writer
- ✅ **Functional options**: clean, composable API via `WithConfig`, `WithWriter`, `WithExtraContextFields`
- ✅ **Validation**: config errors are collected and reported clearly
- ✅ **Zero dependencies** beyond Go standard library

## Quick Start

```go
import "github.com/Averlex/logkit"

// Create a default logger (JSON, ERROR level, stdout)
logger, err := logkit.NewLogger()
if err != nil {
    panic(err)
}

logger.Info(context.Background(), "Application started")
// Output: {"time":"05.04.2025 10:00:00.000","level":"INFO","msg":"Application started"}
```

## Custom Configuration

```go
// Configure logger with custom options
logger, err := logkit.NewLogger(
    logkit.WithConfig(map[string]any{
        "format":        "text",
        "level":         "debug",
        "time_template": "15:04:05",
        "log_stream":    "stdout",
    }),
)
if err != nil {
    log.Fatal(err)
}

logger.Debug(ctx, "Debug message")
// Output: time=15:04:05 level=DEBUG msg="Debug message"
```

## Context-Aware Logging

Use `WithExtraContextFields` to automatically include values from `context.Context` in every log entry.

```go
var RequestIDKey struct{}
func (RequestIDKey) String() string { return "request_id" }

logger, _ := logkit.NewLogger(logkit.WithExtraContextFields(RequestIDKey))
ctx := context.WithValue(context.Background(), RequestIDKey, "abc-123")
logger.Info(ctx, "Handling request")
// Output: {"time":"...","level":"INFO","msg":"Handling request","request_id":"abc-123"}
```

### Key Requirements

Only `string` and `fmt.Stringer` keys are supported. This ensures compatibility with `slog`, which requires attribute keys to be strings.

- If the context `key` is a string, it is used directly.
- If the key implements `fmt.Stringer`, its `String()` method provides the attribute name.
- All other key types are ignored.

### Special Case: slog.Attr

If a value in the context is of type `slog.Attr`, it is logged **as-is**, preserving its key, value, and type. This allows full control over the logged attribute.

```go
ctx := context.WithValue(ctx, key, slog.Int("user_id", 42))
logger.Info(ctx, "User action")
// Output includes: "user_id":42 with correct type
```

This makes `logkit` compatible with advanced `slog` patterns while keeping the API simple.

## Custom Log Levels

| Level     | Use case                            |
| --------- | ----------------------------------- |
| `TRACE`   | Very detailed debugging             |
| `DEBUG`   | Standard debug info                 |
| `VERBOSE` | Between DEBUG and INFO              |
| `INFO`    | General information                 |
| `WARN`    | Potential issues                    |
| `ERROR`   | Errors that don't crash the app     |
| `FATAL`   | Critical errors; calls `os.Exit(1)` |

```go
logger.Trace(ctx, "Entering function")
logger.Verbose(ctx, "Processing batch")
logger.Fatal(ctx, "Failed to initialize", "error", err)
```

## Advanced Usage

### Custom Writer

```go
file, _ := os.Create("app.log")
logger, _ := logkit.NewLogger(logkit.WithWriter(file))
```

### Functional Options

Options can be combined:

```go
logger, _ := logkit.NewLogger(
    logkit.WithDefaults(),
    logkit.WithConfig(customConfig),
    logkit.WithExtraContextFields(RequestIDKey, UserIDKey),
)
```

## Error Handling

`NewLogger` returns an error if configuration is invalid:

```go
if logger, err := logkit.NewLogger(invalidConfig); err != nil {
    log.Printf("Logger setup failed: %v", err)
    // Example: "config data is invalid: invalid_type=time_template, invalid_value=level,format"
}
```

> ⚠️ **Note**: Validation errors are accumulated — you’ll see all issues at once, not just the first one.

## Testing

`logkit` is designed to be testable:

- Accepts any `io.Writer` (use `WithWriter` option).
- Supports JSON format - parse easily with `encoding/json`
- No global state.

Here is a quick example:

```go
var buf bytes.Buffer
logger, _ := logkit.NewLogger(logkit.WithWriter(&buf))

logger.Info(context.Background(), "test message")
// Parse buf.String() as JSON in test
```

## Installation

```bash
go get github.com/Averlex/logkit
```
