package buffer

import (
	"sync"

	"github.com/kechako/seezlog/internal/color"
)

var bufPool = sync.Pool{
	New: func() any {
		b := make(Buffer, 0, 1024)
		return &b
	},
}

// A Buffer is a byte slice that can be used to build a log record.
// It is used with a sync.Pool to reduce allocations.
type Buffer []byte

func New() *Buffer {
	return bufPool.Get().(*Buffer)
}

func (b *Buffer) Free() {
	if b == nil {
		return
	}
	b.Reset()
	bufPool.Put(b)
}

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) Reset() {
	*b = (*b)[:0]
}

func (b *Buffer) WriteColor(color color.Color) {
	*b = append(*b, color...)
}

func (b *Buffer) ResetColor() {
	*b = append(*b, color.Reset...)
}
