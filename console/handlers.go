package console

import (
	"net/http"

	"github.com/LYH2263/go-grpcproxy"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "console/static/index.html")
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.kitHook != nil {
		_, _ = w.Write([]byte(s.kitHook()))
		return
	}
	_, _ = w.Write([]byte(grpcproxy.InspectSummary(grpcproxy.New())))
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true}`))
}
