# seezlog Specification

## 1. Overview

`seezlog` is an implementation of `slog.Handler` for Go's standard `log/slog` package.

## 2. Purpose

The primary goal is to provide a human-readable log output (pretty print) that is easy for developers to visually parse in a console.

## 3. Guiding Principles

- **Performance**: The implementation must minimize memory allocations to avoid impacting application performance. It utilizes an internal buffer pool (`sync.Pool`) to achieve this.

## 4. Output Format

The log record consists of the following components:

```
[Timestamp] [Log Level] [Source Location] [Message] {Attributes}
```

### 4.1. Component Specifications

#### Timestamp

- **Format**: `YYYY-MM-DD HH:MM:SS.ms` (Go format string: `2006-01-02 15:04:05.000`)

#### Log Level

- Displays the level name (e.g., `DEBUG`, `INFO`, `WARN`, `ERROR`).
- **Fixed-width of 5 characters, right-aligned**. (e.g., `  INFO`, ` ERROR`)

#### Source Location

- This component is only displayed if `AddSource` in `slog.HandlerOptions` is set to `true`.
- It shows the last two elements of the file path (the parent directory and the file name) and the line number, e.g., `directory/file.go:123`.
- The output is **left-aligned** and padded with spaces on the right to a **minimum width of 24 characters**.
- It is followed by a single space to separate it from the message.

#### Message

- The raw message from the `slog.Record`.

#### Attributes

- Key-value pairs from the `slog.Record`, formatted as `{key=value, ...}`.
- If there are no attributes, this entire component (including the curly braces) is omitted.

### 4.2. Attribute Formatting Rules

#### Key

- Keys without spaces or special characters are displayed as-is.
- Keys containing spaces, `=`, or other special characters are enclosed in double quotes (`"`).

#### Value

The format of a value depends on its type.

| Type                   | Example Format               | Notes                                                                          |
| ---------------------- | ---------------------------- | ------------------------------------------------------------------------------ |
| `string`               | `"hello world"`              | Always enclosed in double quotes. Any `"` within the value is escaped as `\"`. |
| `int`, `uint`, `float` | `123`, `3.14`                | The number is displayed as-is.                                                 |
| `bool`                 | `true`, `false`              | The boolean is displayed as-is.                                                |
| `time.Time`            | `2025-06-26T16:20:00.123Z`   | Formatted as RFC3339Nano.                                                      |
| `time.Duration`        | `1m23.45s`                   | Follows the output of Go's `String()` method.                                  |
| `error`                | `"file not found"`           | The string returned by the `Error()` method, enclosed in double quotes.        |
| `slice`, `array`       | `[1, 2, "three"]`            | Enclosed in square brackets `[]`, with elements separated by commas `,`.       |
| `map`                  | `{"key1":1, "key2":"val"}`   | Enclosed in curly braces `{}`. Keys are sorted as strings for stable output.   |
| `struct`               | `{FieldA:"val", FieldB:123}` | Enclosed in curly braces `{}`. **Only exported fields** are displayed.         |
| `nil`                  | `nil`                        | Displayed as the string `nil`.                                                 |

## 5. Optional Features

### 5.1. `slog.HandlerOptions` Support

Standard options can be configured by passing a `*slog.HandlerOptions` struct during initialization with `seezlog.NewHandler`.

- **`Level`**: Sets the minimum log level that the handler will process.
- **`AddSource`**: Toggles the inclusion of the source code location (file and line number) in the log output. See section 4.1 for details.
- **`ReplaceAttr`**: Provides a function to dynamically modify attributes before they are logged. This function also receives the group context from `WithGroup`.

### 5.2. Coloring

- Colorized output can be enabled by setting the `colors` argument to `true` in `seezlog.NewHandler`. It uses 16-color ANSI codes.

| Component               | Color        | Notes                           |
| ----------------------- | ------------ | ------------------------------- |
| Timestamp               | Bright Black |                                 |
| Source Location         | Bright Black |                                 |
| Log Level (DEBUG)       | Green        |                                 |
| Log Level (INFO)        | Blue         |                                 |
| Log Level (WARN)        | Yellow       |                                 |
| Log Level (ERROR)       | Red          |                                 |
| Message                 | Default      | No color                        |
| Attribute Symbols       | Bright Black | `{`, `}`, `=`, `,` etc.         |
| Attribute Key           | Bright Black |                                 |
| Attribute Value         | Default      | No color                        |
| Attribute (error value) | Red          | Only for values of type `error` |

## 6. Example

```go
package main

import (
	"errors"
	"log/slog"
	"os"
	"github.com/kechako/seezlog"
)

func main() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	h := seezlog.NewHandler(os.Stdout, true, opts)
	logger := slog.New(h).With("service", "user-service")

	logger.Info("User logged in", "user_id", 12345)
	logger.Error(
		"Failed to update profile",
		slog.Any("error", errors.New("connection timeout")),
		slog.Group("request", "method", "POST", "path", "/api/users/12345"),
	)
}
```

**Example Output (with colors):**

(This output is actually colorized in the terminal)

```
2025-06-26 16:30:00.123   INFO main.go:18                 User logged in {service="user-service", user_id=12345}
2025-06-26 16:30:00.456  ERROR main.go:19                 Failed to update profile {service="user-service", error="connection timeout", request={method="POST", path="/api/users/12345"}}
```
