package grpcproxy

import "errors"

var (
	ErrClosed        = errors.New("grpcproxy: kit closed")
	ErrCorruptFrame  = errors.New("grpcproxy: corrupt frame")
	ErrBackpressure  = errors.New("grpcproxy: backpressure exceeded")
	ErrPartialFrame  = errors.New("grpcproxy: partial frame read")
	ErrChainAbort    = errors.New("grpcproxy: handler chain aborted")
	ErrInvalidWindow = errors.New("grpcproxy: invalid credit window")
)
