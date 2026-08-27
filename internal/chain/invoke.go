package chain

import (
	"context"
	"errors"

	"github.com/LYH2263/go-grpcproxy/internal/frame"
)

var ErrAbort = errors.New("chain: aborted")

// Invoke runs handlers in order; abort stops the chain.
func (c *Chain) Invoke(ctx context.Context, f *frame.Frame) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	st, ok := StateFromContext(ctx)
	if !ok {
		st = &State{}
		ctx = WithChainState(ctx, st)
	}
	handlers := c.cloneHandlers()
	var run Next
	run = func(ctx context.Context, f *frame.Frame) error {
		if st.Aborted {
			return ErrAbort
		}
		if len(handlers) == 0 {
			st.LastSeq = f.Seq
			return nil
		}
		h := handlers[0]
		handlers = handlers[1:]
		st.Invoked++
		return h(ctx, f, run)
	}
	return run(ctx, f)
}
