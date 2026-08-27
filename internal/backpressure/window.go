package backpressure

import (
	"context"
	"sync"
)

// Window tracks credit-based flow control for one direction.
type Window struct {
	mu      sync.Mutex
	credits int
	max     int
	name    string
}

func NewWindow(max int, name string) *Window {
	if max <= 0 {
		max = 256
	}
	return &Window{credits: max, max: max, name: name}
}

func (w *Window) Name() string { return w.name }

func (w *Window) Credits() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.credits
}

func (w *Window) Acquire(ctx context.Context, n int) error {
	if n <= 0 {
		return nil
	}
	for {
		w.mu.Lock()
		if w.credits >= n {
			w.credits -= n
			w.mu.Unlock()
			return nil
		}
		w.mu.Unlock()
		return nil
	}
}

func (w *Window) Release(n int) {
	if n <= 0 {
		return
	}
	w.mu.Lock()
	w.credits += n
	if w.credits > w.max {
		w.credits = w.max
	}
	w.mu.Unlock()
}
