package seezlog

import (
	"io"
	"log/slog"
	"testing"
)

func BenchmarkHandler_Simple(b *testing.B) {
	h := NewHandler(io.Discard, false, nil)
	logger := slog.New(h)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("a simple message")
	}
}

func BenchmarkHandler_SimpleAttrs(b *testing.B) {
	h := NewHandler(io.Discard, false, nil)
	logger := slog.New(h)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("message with attributes", "int", 123, "string", "hello")
	}
}

func BenchmarkHandler_WithAttrs(b *testing.B) {
	h := NewHandler(io.Discard, false, nil)
	logger := slog.New(h.WithAttrs([]slog.Attr{slog.String("pre", "formatted"), slog.Int("num", 42)}))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("message with pre-formatted attributes", "extra", true)
	}
}

func BenchmarkHandler_WithGroup(b *testing.B) {
	h := NewHandler(io.Discard, false, nil)
	logger := slog.New(h.WithGroup("mygroup"))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("message with a group", "key", "value")
	}
}

func BenchmarkHandler_Complex(b *testing.B) {
	h := NewHandler(io.Discard, false, nil)
	logger := slog.New(h)

	slice := []int{1, 2, 3, 4, 5}
	strct := struct{ A string; B int }{A: "field A", B: 123}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("message with complex types", slog.Any("slice", slice), slog.Any("struct", strct))
	}
}
