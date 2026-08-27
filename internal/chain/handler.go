package chain

import (
	"context"

	"github.com/LYH2263/go-grpcproxy/internal/frame"
)

// Next continues the handler chain.
type Next func(ctx context.Context, f *frame.Frame) error

// Handler mutates or inspects a frame before forwarding.
type Handler func(ctx context.Context, f *frame.Frame, next Next) error
