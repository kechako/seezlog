package seezlog

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, false, nil)
	logger := slog.New(h)

	logger.Info("test message", "key1", "value1", "key2", 123)

	got := buf.String()
	if !strings.HasSuffix(got, "  INFO test message {key1=\"value1\", key2=123}\n") {
		t.Errorf("unexpected log output:\ngot:  %q", got)
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, false, nil)
	logger := slog.New(h.WithAttrs([]slog.Attr{slog.String("attr1", "val1")}))

	logger.Info("test message", "key1", "value1")

	got := buf.String()
	if !strings.HasSuffix(got, "  INFO test message {attr1=\"val1\", key1=\"value1\"}\n") {
		t.Errorf("unexpected log output:\ngot:  %q", got)
	}
}

func TestHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, false, nil)
	logger := slog.New(h.WithGroup("group1"))

	logger.Info("test message", "key1", "value1")

	got := buf.String()
	if !strings.HasSuffix(got, "  INFO test message {group1={key1=\"value1\"}}\n") {
		t.Errorf("unexpected log output:\ngot:  %q", got)
	}
}

func TestHandler_WithGroupAndAttrs(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, false, nil)
	logger := slog.New(h.WithAttrs([]slog.Attr{slog.String("attr1", "val1")}).WithGroup("group1"))

	logger.Info("test message", "key1", "value1")

	got := buf.String()
	if !strings.HasSuffix(got, "  INFO test message {group1={attr1=\"val1\", key1=\"value1\"}}\n") {
		t.Errorf("unexpected log output:\ngot:  %q", got)
	}
}

func TestHandler_AddSource(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, false, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(h)

	logger.Info("test message")

	got := buf.String()
	if !strings.Contains(got, "seezlog/handler_test.go:") {
		t.Errorf("log output should contain short source file info: %q", got)
	}
}

func TestHandler_ReplaceAttr(t *testing.T) {
	var buf bytes.Buffer
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "key1" {
			return slog.String("replaced_key", "replaced_value")
		}
		return a
	}
	h := NewHandler(&buf, false, &slog.HandlerOptions{ReplaceAttr: replace})
	logger := slog.New(h)

	logger.Info("test message", "key1", "value1", "key2", 123)

	got := buf.String()
	if !strings.Contains(got, "{replaced_key=\"replaced_value\", key2=123}") {
		t.Errorf("unexpected log output: %q", got)
	}
}

func TestHandler_Error(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, false, nil)
	logger := slog.New(h)

	err := errors.New("this is an error")
	logger.Error("error message", slog.Any("err", err))

	got := buf.String()
	if !strings.Contains(got, `err="this is an error"`) {
		t.Errorf("unexpected error format: got %q", got)
	}
}