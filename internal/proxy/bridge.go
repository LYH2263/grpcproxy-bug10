package proxy

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/LYH2263/go-grpcproxy/internal/backpressure"
	"github.com/LYH2263/go-grpcproxy/internal/chain"
	"github.com/LYH2263/go-grpcproxy/internal/frame"
)

var sessCounter uint64

// Bridge orchestrates frame relay between paired streams.
type Bridge struct {
	codec    *frame.Codec
	chain    *chain.Chain
	policy   backpressure.Policy
	credits  int
	sem      *backpressure.Semaphore
}

func NewBridge(codec *frame.Codec, ch *chain.Chain, pol backpressure.Policy, credits int) *Bridge {
	return &Bridge{
		codec: codec, chain: ch, policy: pol, credits: credits,
		sem: backpressure.NewSemaphore(32),
	}
}

func (b *Bridge) NewSession() *Session {
	id := fmt.Sprintf("sess-%d", atomic.AddUint64(&sessCounter, 1))
	return newSession(id, b.credits)
}

// Pipe relays bidirectional traffic through frame codec and handler chain.
func (b *Bridge) Pipe(ctx context.Context, sess *Session, client, server io.ReadWriter) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := b.sem.Acquire(ctx); err != nil {
		return err
	}
	defer b.sem.Release()
	errC := make(chan error, 2)
	go func() {
		errC <- b.relayOne(ctx, sess, client, server, true)
	}()
	go func() {
		errC <- b.relayOne(ctx, sess, server, client, false)
	}()
	var first error
	for i := 0; i < 2; i++ {
		if err := <-errC; err != nil && first == nil {
			first = err
		}
	}
	return first
}

// RelayOneForTest exposes single-direction relay for acceptance tests.
func (b *Bridge) RelayOneForTest(ctx context.Context, sess *Session, in io.Reader, out io.Writer, upstream bool) error {
	return b.relayOne(ctx, sess, in, out, upstream)
}

func (b *Bridge) relayOne(ctx context.Context, sess *Session, in io.Reader, out io.Writer, upstream bool) error {
	r := frame.NewReader(b.codec, in)
	w := frame.NewWriter(b.codec, out)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		f, err := r.ReadFrame(ctx)
		if err != nil {
			return err
		}
		win := sess.Credits().Downstream
		if upstream {
			win = sess.Credits().Upstream
		}
		if err := win.Acquire(ctx, 1); err != nil {
			b.codec.Release(f)
			return err
		}
		if b.chain != nil {
			if err := b.chain.Invoke(ctx, f); err != nil {
				b.codec.Release(f)
				return err
			}
		}
		if err := w.WriteFrame(ctx, f); err != nil {
			b.codec.Release(f)
			return err
		}
		sess.addOut(int64(len(f.Payload)))
		if upstream {
			sess.Credits().ReleaseDownstream(1)
		} else {
			sess.Credits().ReleaseUpstream(1)
		}
		b.codec.Release(f)
	}
}
