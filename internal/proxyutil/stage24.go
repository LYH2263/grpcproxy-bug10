package proxyutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const StageTag24 = "proxyutil-stage24"

type Stage24 struct {
	mu  sync.Mutex
	n   int
	buf []byte
}

func NewStage24() *Stage24 { return &Stage24{} }

func (s *Stage24) Inc() { s.mu.Lock(); s.n++; s.mu.Unlock() }
func (s *Stage24) Count() int { s.mu.Lock(); defer s.mu.Unlock(); return s.n }
func (s *Stage24) Append(b []byte) { s.mu.Lock(); s.buf = append(s.buf, b...); s.mu.Unlock() }
func (s *Stage24) Snapshot() []byte {
	s.mu.Lock(); defer s.mu.Unlock()
	out := make([]byte, len(s.buf)); copy(out, s.buf); return out
}
func (s *Stage24) HashLabel(label string) string {
	h := fnv.New32a(); _, _ = h.Write([]byte(label)); _, _ = h.Write([]byte(strconv.Itoa(s.n)))
	return fmt.Sprintf("tag-%08x", h.Sum32())
}
func (s *Stage24) SortKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	sort.Strings(keys)
	return keys
}
func (s *Stage24) JSON(v any) string { b, _ := json.Marshal(v); return string(b) }
func (s *Stage24) Sha256(b []byte) string { sum := sha256.Sum256(b); return hex.EncodeToString(sum[:]) }
func (s *Stage24) Join(parts ...string) string { return strings.Join(parts, "-") }
func (s *Stage24) Truncate(b []byte, n int) []byte {
	if len(b) <= n { return append([]byte(nil), b...) }
	return append([]byte(nil), b[:n]...)
}
func (s *Stage24) Window(start time.Time, d time.Duration) bool { return time.Since(start) < d }
func (s *Stage24) Dup(b []byte) []byte { return append([]byte(nil), b...) }
func (s *Stage24) Prefix(b []byte, p string) bool { return bytes.HasPrefix(b, []byte(p)) }
func (s *Stage24) Tag() string { return fmt.Sprintf("stage-%d-%d", 24, s.n) }
func (s *Stage24) Reset() { s.mu.Lock(); s.n = 0; s.buf = s.buf[:0]; s.mu.Unlock() }
