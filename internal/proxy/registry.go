package proxy

import "sync"

var (
	regMu sync.RWMutex
	reg   = map[string]*Session{}
)

func registerSession(s *Session) {
	regMu.Lock()
	reg[s.id] = s
	regMu.Unlock()
}

func LookupSession(id string) (*Session, bool) {
	regMu.RLock()
	defer regMu.RUnlock()
	s, ok := reg[id]
	return s, ok
}

func RegistrySize() int {
	regMu.RLock()
	defer regMu.RUnlock()
	return len(reg)
}
