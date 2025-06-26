# seezlog

`seezlog` is a Go library that provides a `slog.Handler` implementation for human-readable, pretty-printed log output.

## Features

- **Human-Readable Format:** Logs are formatted for easy visual parsing, with colors and clear structure.
- **High Performance:** Designed to minimize memory allocations for excellent performance.
- **`slog` Compliant:** Implements the standard `log/slog.Handler` interface and supports `slog.HandlerOptions`, including `AddSource` and `ReplaceAttr`.
- **Customizable:** Supports color output toggling.

## Installation

```bash
go get github.com/kechako/seezlog
```

## Usage

```go
package main

import (
	"log/slog"
	"os"

	"github.com/kechako/seezlog"
)

func main() {
	// Create a new handler with colors enabled and AddSource option
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	h := seezlog.NewHandler(os.Stdout, true, opts)
	
	// Set the new handler as the default
	logger := slog.New(h)
	slog.SetDefault(logger)

	// Log messages
	slog.Debug("Debug message with details", slog.String("user", "debug"))
	slog.Info("User logged in", "user_id", 12345, "status", "active")
	slog.Warn("Disk space is running low", slog.Float64("free_gb", 15.2))
	slog.Error("Failed to connect to database", slog.Any("error", os.ErrPermission))

	// Create a logger with pre-defined attributes and groups
	requestLogger := logger.With("service", "api").WithGroup("request")
	requestLogger.Info("Incoming request", "method", "GET", "path", "/api/v1/users")
}
```

## Log Format

```
[Timestamp] [LEVEL] [Source (optional)] [Message] {key=value, ...}
```

**Example Output:**

```
2025-06-26 16:20:00.123   INFO main.go:25                User logged in {user_id=12345, status="active"}
2025-06-26 16:20:00.123  ERROR main.go:26                Failed to connect to database {error="permission denied"}
2025-06-26 16:20:00.123   INFO main.go:30                Incoming request {service="api", request={method="GET", path="/api/v1/users"}}
```

## Specification

For detailed information on formatting, coloring, and options, please see the [specification document](docs/specification.md).
