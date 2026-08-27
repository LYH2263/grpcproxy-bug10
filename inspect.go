package grpcproxy

import (
	"encoding/json"
	"fmt"

	"github.com/LYH2263/go-grpcproxy/internal/proxy"
)

// InspectSummary renders a JSON summary for console diagnostics.
func InspectSummary(k *Kit) string {
	s := k.Snapshot()
	out := map[string]any{
		"sessions": s.Sessions,
		"frames":   s.Frames,
		"bytes":    s.Bytes,
		"blocked":  s.Blocked,
		"registry": proxy.RegistrySize(),
	}
	b, _ := json.Marshal(out)
	return string(b)
}

// InspectSession renders one session by id.
func InspectSession(id string) string {
	sess, ok := proxy.LookupSession(id)
	if !ok {
		return fmt.Sprintf(`{"error":"session not found"}`)
	}
	st := sess.Stats()
	out := map[string]any{
		"id": id, "frames_in": st.FramesIn, "frames_out": st.FramesOut,
		"bytes_in": st.BytesIn, "bytes_out": st.BytesOut,
	}
	b, _ := json.Marshal(out)
	return string(b)
}
