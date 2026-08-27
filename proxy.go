package grpcproxy

import (
	"sync"
	"sync/atomic"

	"github.com/LYH2263/go-grpcproxy/internal/backpressure"
	"github.com/LYH2263/go-grpcproxy/internal/chain"
	"github.com/LYH2263/go-grpcproxy/internal/frame"
	"github.com/LYH2263/go-grpcproxy/internal/proxy"
)

// Kit coordinates stream proxy pipelines.
type Kit struct {
	bridge        *proxy.Bridge
	pool          *frame.Pool
	codec         *frame.Codec
	policy        backpressure.Policy
	maxPayload    int
	windowCredits int
	chain         *chain.Chain

	mu     sync.RWMutex
	closed atomic.Bool
	stats  Stats
}

// New constructs a Kit with default dependencies.
func New(opts ...Option) *Kit {
	k := &Kit{
		pool:          frame.NewPool(),
		maxPayload:    64 << 10,
		windowCredits: 256,
		policy:        backpressure.DefaultPolicy(),
	}
	for _, o := range opts {
		o(k)
	}
	k.codec = frame.NewCodec(k.pool, k.maxPayload)
	k.chain = chain.New()
	k.bridge = proxy.NewBridge(k.codec, k.chain, k.policy, k.windowCredits)
	return k
}

// Close stops accepting new work.
func (k *Kit) Close() error {
	if !k.closed.CompareAndSwap(false, true) {
		return nil
	}
	return nil
}

// Bridge exposes the underlying stream bridge for diagnostics.
func (k *Kit) Bridge() *proxy.Bridge {
	return k.bridge
}

// Chain returns the handler chain bound to this kit.
func (k *Kit) Chain() *chain.Chain {
	return k.chain
}

func (k *Kit) bumpFrames(n int64) {
	k.mu.Lock()
	k.stats.Frames += n
	k.mu.Unlock()
}

func (k *Kit) bumpBytes(n int64) {
	k.mu.Lock()
	k.stats.Bytes += n
	k.mu.Unlock()
}
