package chain

import (
	"sync"
)

// Chain wires ordered handlers around frame forwarding.
type Chain struct {
	mu       sync.RWMutex
	handlers []Handler
}

func New() *Chain { return &Chain{} }

func (c *Chain) Use(h Handler) {
	c.mu.Lock()
	c.handlers = append(c.handlers, h)
	c.mu.Unlock()
}

func (c *Chain) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.handlers)
}

func (c *Chain) cloneHandlers() []Handler {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Handler, len(c.handlers))
	copy(out, c.handlers)
	return out
}
