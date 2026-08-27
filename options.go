package grpcproxy

import (
	"github.com/LYH2263/go-grpcproxy/internal/backpressure"
	"github.com/LYH2263/go-grpcproxy/internal/frame"
)

// Option configures a Kit.
type Option func(*Kit)

// WithMaxPayload sets frame codec max payload bytes.
func WithMaxPayload(n int) Option {
	return func(k *Kit) {
		if n > 0 {
			k.maxPayload = n
		}
	}
}

// WithWindowCredits sets default backpressure window credits.
func WithWindowCredits(n int) Option {
	return func(k *Kit) {
		if n > 0 {
			k.windowCredits = n
		}
	}
}

// WithFramePool overrides the frame buffer pool.
func WithFramePool(p *frame.Pool) Option {
	return func(k *Kit) {
		k.pool = p
	}
}

// WithPolicy overrides backpressure policy.
func WithPolicy(p backpressure.Policy) Option {
	return func(k *Kit) {
		k.policy = p
	}
}
