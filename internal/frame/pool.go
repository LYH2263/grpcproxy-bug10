package frame

import "sync"

// Pool recycles frame payload buffers.
type Pool struct {
	pool sync.Pool
}

func NewPool() *Pool {
	return &Pool{
		pool: sync.Pool{New: func() any { return make([]byte, 0, 4096) }},
	}
}

func (p *Pool) Get(size int) []byte {
	b := p.pool.Get().([]byte)
	if cap(b) < size {
		b = make([]byte, size)
	} else {
		b = b[:size]
	}
	return b
}

func (p *Pool) Put(b []byte) {
	if cap(b) > 1<<20 {
		return
	}
	p.pool.Put(b[:0])
}
