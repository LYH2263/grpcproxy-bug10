package chain

import "context"

type ctxKey struct{}

// WithChainState attaches mutable chain state to ctx.
func WithChainState(ctx context.Context, st *State) context.Context {
	return context.WithValue(ctx, ctxKey{}, st)
}

// StateFromContext returns chain state if present.
func StateFromContext(ctx context.Context) (*State, bool) {
	st, ok := ctx.Value(ctxKey{}).(*State)
	return st, ok
}

// State tracks per-session chain progression.
type State struct {
	Invoked int
	Aborted bool
	LastSeq uint64
}
