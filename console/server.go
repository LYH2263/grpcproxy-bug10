package console

import (
	"fmt"
	"net/http"
)

// Server serves the proxy admin UI.
type Server struct {
	addr    string
	mux     *http.ServeMux
	kitHook func() string
}

func New(addr string, kitHook func() string) *Server {
	s := &Server{addr: addr, mux: http.NewServeMux(), kitHook: kitHook}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/health", s.handleHealth)
}

func (s *Server) ListenAndServe() error {
	fmt.Printf("grpcproxy console on http://%s\n", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}
