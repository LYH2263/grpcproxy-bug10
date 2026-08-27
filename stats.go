package grpcproxy

// Snapshot returns a copy of kit stats.
func (k *Kit) Snapshot() Stats {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.stats
}
