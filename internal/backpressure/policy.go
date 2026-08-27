package backpressure

import (
	"context"
	"errors"
)

var ErrBackpressure = errors.New("backpressure: window exhausted")

// Policy decides whether to block or drop when windows exhaust.
type Policy interface {
	OnBlock(ctx context.Context) error
}

type defaultPolicy struct{}

func DefaultPolicy() Policy { return defaultPolicy{} }

func (defaultPolicy) OnBlock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ErrBackpressure
	}
}
