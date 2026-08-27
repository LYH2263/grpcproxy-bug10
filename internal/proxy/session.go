package proxy

import (
	"sync"
	"sync/atomic"

	"github.com/LYH2263/go-grpcproxy/internal/backpressure"
)

// SessionStats holds relay counters.
type SessionStats struct {
	FramesIn   int64
	FramesOut  int64
	BytesIn    int64
	BytesOut   int64
	Backlogged int64
	Dropped    int64
}

// Session is one proxied stream pair with credit windows.
type Session struct {
	id      string
	credits *backpressure.CreditPair
	stats   SessionStats
	seq     uint64
	mu      sync.Mutex
}

func newSession(id string, credits int) *Session {
	s := &Session{
		id:      id,
		credits: backpressure.NewCreditPair(credits),
	}
	registerSession(s)
	return s
}

func (s *Session) ID() string { return s.id }

func (s *Session) Stats() SessionStats { return s.stats }

func (s *Session) NextSeq() uint64 {
	return atomic.AddUint64(&s.seq, 1)
}

func (s *Session) CurrentSeq() uint64 {
	return atomic.LoadUint64(&s.seq)
}

func (s *Session) ResetSeq() {
	atomic.StoreUint64(&s.seq, 0)
}

func (s *Session) Credits() *backpressure.CreditPair { return s.credits }

func (s *Session) addIn(n int64) {
	atomic.AddInt64(&s.stats.FramesIn, 1)
	atomic.AddInt64(&s.stats.BytesIn, n)
}

func (s *Session) addOut(n int64) {
	atomic.AddInt64(&s.stats.FramesOut, 1)
	atomic.AddInt64(&s.stats.BytesOut, n)
}
