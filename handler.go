package seezlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"slices"

	"github.com/kechako/seezlog/internal/buffer"
	"github.com/kechako/seezlog/internal/color"
)

var levelToColor = map[slog.Level]color.Color{
	slog.LevelDebug: color.Green,
	slog.LevelInfo:  color.Blue,
	slog.LevelWarn:  color.Yellow,
	slog.LevelError: color.Red,
}

type Handler struct {
	opts   *slog.HandlerOptions
	colors bool
	w      io.Writer
	mu     sync.Mutex
	groups []string
	attrs  []slog.Attr
}

func NewHandler(w io.Writer, colors bool, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		opts:   opts,
		colors: colors,
		w:      w,
	}
}

func (h *Handler) clone() *Handler {
	return &Handler{
		opts:   h.opts,
		colors: h.colors,
		w:      h.w,
		groups: slices.Clone(h.groups),
		attrs:  slices.Clone(h.attrs),
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	b := buffer.New()
	defer b.Free()

	// Time
	if h.colors {
		b.WriteColor(color.BrightBlack)
	}
	*b = r.Time.AppendFormat(*b, "2006-01-02 15:04:05.000")
	if h.colors {
		b.ResetColor()
	}
	b.WriteString(" ")

	// Level
	if h.colors {
		b.WriteColor(levelToColor[r.Level])
	}
	levelStr := r.Level.String()
	padding := 5 - len(levelStr)
	for range padding {
		b.WriteString(" ")
	}
	b.WriteString(levelStr)
	if h.colors {
		b.ResetColor()
	}
	b.WriteString(" ")

	// Source
	if h.opts.AddSource && r.PC != 0 {
		pcs := [1]uintptr{r.PC}
		fs := runtime.CallersFrames(pcs[:])
		f, _ := fs.Next()
		if f.File != "" {
			shortPath := getShortPath(f)
			if h.colors {
				b.WriteColor(color.BrightBlack)
			}
			start := len(*b)
			b.WriteString(shortPath)
			b.WriteString(":")
			*b = strconv.AppendInt(*b, int64(f.Line), 10)
			end := len(*b)

			for i := end - start; i < 24; i++ {
				b.WriteString(" ")
			}
			if h.colors {
				b.ResetColor()
			}
			b.WriteString(" ")
		}
	}

	// Message
	b.WriteString(r.Message)

	// Attributes
	hasAttrs := len(h.attrs) > 0 || r.NumAttrs() > 0
	if hasAttrs {
		b.WriteString(" {")
	}

	for _, g := range h.groups {
		h.appendKey(b, g)
		b.WriteString("={")
	}

	for i, a := range h.attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		h.appendAttr(b, a)
	}

	if len(h.attrs) > 0 && r.NumAttrs() > 0 {
		b.WriteString(", ")
	}

	count := 0
	r.Attrs(func(a slog.Attr) bool {
		if count > 0 {
			b.WriteString(", ")
		}
		h.appendAttr(b, a)
		count++
		return true
	})

	for range h.groups {
		b.WriteString("}")
	}

	if hasAttrs {
		b.WriteString("}")
	}

	*b = append(*b, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*b)
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)
	return h2
}

func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *Handler) appendAttr(b *buffer.Buffer, a slog.Attr) {
	if h.opts.ReplaceAttr != nil {
		a = h.opts.ReplaceAttr(h.groups, a)
	}
	h.appendKey(b, a.Key)
	b.WriteString("=")
	h.appendValue(b, a.Value)
}

func (h *Handler) appendKey(b *buffer.Buffer, key string) {
	if h.colors {
		b.WriteColor(color.BrightBlack)
	}
	if needsQuoting(key) {
		*b = strconv.AppendQuote(*b, key)
	} else {
		b.WriteString(key)
	}
	if h.colors {
		b.ResetColor()
	}
}

func (h *Handler) appendValue(b *buffer.Buffer, v slog.Value) {
	if h.colors && v.Kind() == slog.KindAny {
		if _, ok := v.Any().(error); ok {
			b.WriteColor(color.Red)
			defer b.ResetColor()
		}
	}

	if err, ok := v.Any().(error); ok {
		*b = strconv.AppendQuote(*b, err.Error())
		return
	}

	switch v.Kind() {
	case slog.KindString:
		*b = strconv.AppendQuote(*b, v.String())
	case slog.KindInt64:
		*b = strconv.AppendInt(*b, v.Int64(), 10)
	case slog.KindUint64:
		*b = strconv.AppendUint(*b, v.Uint64(), 10)
	case slog.KindFloat64:
		*b = strconv.AppendFloat(*b, v.Float64(), 'f', -1, 64)
	case slog.KindBool:
		*b = strconv.AppendBool(*b, v.Bool())
	case slog.KindDuration:
		b.WriteString(v.Duration().String())
	case slog.KindTime:
		*b = v.Time().AppendFormat(*b, "2006-01-02T15:04:05.999Z07:00")
	case slog.KindGroup:
		attrs := v.Group()
		b.WriteString("{")
		for i, a := range attrs {
			if i > 0 {
				b.WriteString(", ")
			}
			h.appendAttr(b, a)
		}
		b.WriteString("}")
	case slog.KindAny:
		h.appendReflectedValue(b, reflect.ValueOf(v.Any()))
	}
}

func (h *Handler) appendReflectedValue(b *buffer.Buffer, val reflect.Value) {
	if !val.IsValid() || val.IsZero() {
		b.WriteString("nil")
		return
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		b.WriteString("[")
		for i := range val.Len() {
			if i > 0 {
				b.WriteString(", ")
			}
			h.appendReflectedValue(b, val.Index(i))
		}
		b.WriteString("]")
	case reflect.Map:
		b.WriteString("{")
		keys := make([]string, 0, val.Len())
		for _, k := range val.MapKeys() {
			keys = append(keys, k.String())
		}
		sort.Strings(keys)

		for i, k := range keys {
			if i > 0 {
				b.WriteString(", ")
			}
			h.appendKey(b, k)
			b.WriteString("=")
			h.appendReflectedValue(b, val.MapIndex(reflect.ValueOf(k)))
		}
		b.WriteString("}")
	case reflect.Struct:
		b.WriteString("{")
		t := val.Type()
		count := 0
		for i := range val.NumField() {
			if t.Field(i).IsExported() {
				if count > 0 {
					b.WriteString(", ")
				}
				h.appendKey(b, t.Field(i).Name)
				b.WriteString("=")
				h.appendReflectedValue(b, val.Field(i))
				count++
			}
		}
		b.WriteString("}")
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			b.WriteString("nil")
			return
		}
		h.appendReflectedValue(b, val.Elem())
	default:
		*b = fmt.Append(*b, val.Interface())
	}
}

func needsQuoting(s string) bool {
	for _, r := range s {
		if !strconv.IsPrint(r) || r == ' ' || r == '"' || r == '=' {
			return true
		}
	}
	return false
}

func getShortPath(f runtime.Frame) string {
	// Find the last separator.
	idx := strings.LastIndexByte(f.File, '/')
	if idx == -1 {
		return f.File
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(f.File[:idx], '/')
	if idx == -1 {
		return f.File
	}
	return f.File[idx+1:]
}
