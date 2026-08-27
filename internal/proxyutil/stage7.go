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

const StageTag7 = "proxyutil-stage7"

type Stage7 struct {
	mu  sync.Mutex
	n   int
	buf []byte
}

func NewStage7() *Stage7 { return &Stage7{} }

func (s *Stage7) Inc() { s.mu.Lock(); s.n++; s.mu.Unlock() }
func (s *Stage7) Count() int { s.mu.Lock(); defer s.mu.Unlock(); return s.n }
func (s *Stage7) Append(b []byte) { s.mu.Lock(); s.buf = append(s.buf, b...); s.mu.Unlock() }
func (s *Stage7) Snapshot() []byte {
	s.mu.Lock(); defer s.mu.Unlock()
	out := make([]byte, len(s.buf)); copy(out, s.buf); return out
}
func (s *Stage7) HashLabel(label string) string {
	h := fnv.New32a(); _, _ = h.Write([]byte(label)); _, _ = h.Write([]byte(strconv.Itoa(s.n)))
	return fmt.Sprintf("tag-%08x", h.Sum32())
}
func (s *Stage7) SortKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	sort.Strings(keys)
	return keys
}
func (s *Stage7) JSON(v any) string { b, _ := json.Marshal(v); return string(b) }
func (s *Stage7) Sha256(b []byte) string { sum := sha256.Sum256(b); return hex.EncodeToString(sum[:]) }
func (s *Stage7) Join(parts ...string) string { return strings.Join(parts, "-") }
func (s *Stage7) Truncate(b []byte, n int) []byte {
	if len(b) <= n { return append([]byte(nil), b...) }
	return append([]byte(nil), b[:n]...)
}
func (s *Stage7) Window(start time.Time, d time.Duration) bool { return time.Since(start) < d }
func (s *Stage7) Dup(b []byte) []byte { return append([]byte(nil), b...) }
func (s *Stage7) Prefix(b []byte, p string) bool { return bytes.HasPrefix(b, []byte(p)) }
func (s *Stage7) Tag() string { return fmt.Sprintf("stage-%d-%d", 7, s.n) }
func (s *Stage7) Reset() { s.mu.Lock(); s.n = 0; s.buf = s.buf[:0]; s.mu.Unlock() }
