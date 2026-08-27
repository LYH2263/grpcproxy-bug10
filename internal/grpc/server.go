package grpcx

import (
	"context"
	"fmt"
	"net"
	"sync"
)

// Server wraps a TCP listener stub for proxy admin hooks.
type Server struct {
	addr string
	ln   net.Listener
	mu   sync.Mutex
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.ln = ln
	s.mu.Unlock()
	fmt.Printf("grpcproxy grpc listening on %s\n", s.addr)
	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		_ = c.Close()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	ln := s.ln
	s.mu.Unlock()
	if ln == nil {
		return nil
	}
	done := make(chan struct{})
	go func() {
		_ = ln.Close()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
