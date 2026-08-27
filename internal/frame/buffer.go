package frame

import "sync"

// Buffer holds reusable encode/decode scratch space.
type Buffer struct {
	mu  sync.Mutex
	raw []byte
}

func NewBuffer(cap int) *Buffer {
	return &Buffer{raw: make([]byte, 0, cap)}
}

func (b *Buffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]byte(nil), b.raw...)
}

func (b *Buffer) Reset() {
	b.mu.Lock()
	b.raw = b.raw[:0]
	b.mu.Unlock()
}

func (b *Buffer) Append(p []byte) {
	b.mu.Lock()
	b.raw = append(b.raw, p...)
	b.mu.Unlock()
}
