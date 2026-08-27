package grpcproxy

import (
	"context"
	"io"

	"github.com/LYH2263/go-grpcproxy/internal/proxy"
)

// Bridge pipes bidirectional byte streams through frame codec and handler chain.
type Bridge struct {
	kit    *Kit
	inner  *proxy.Bridge
	sess   *proxy.Session
}

// NewBridge creates a bridge bound to kit defaults.
func (k *Kit) NewBridge() (*Bridge, error) {
	if k.closed.Load() {
		return nil, ErrClosed
	}
	sess := k.bridge.NewSession()
	k.mu.Lock()
	k.stats.Sessions++
	k.mu.Unlock()
	return &Bridge{kit: k, inner: k.bridge, sess: sess}, nil
}

// Pipe relays frames from client to server stream and back.
func (b *Bridge) Pipe(ctx context.Context, client, server io.ReadWriter) error {
	if b.kit.closed.Load() {
		return ErrClosed
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	err := b.inner.Pipe(ctx, b.sess, client, server)
	if err != nil {
		return err
	}
	b.kit.bumpFrames(b.sess.Stats().FramesOut)
	b.kit.bumpBytes(b.sess.Stats().BytesOut)
	return nil
}

// SessionStats returns relay counters for this bridge session.
func (b *Bridge) SessionStats() SessionStats {
	s := b.sess.Stats()
	return SessionStats{
		FramesIn: s.FramesIn, FramesOut: s.FramesOut,
		BytesIn: s.BytesIn, BytesOut: s.BytesOut,
		Backlogged: s.Backlogged, Dropped: s.Dropped,
	}
}
