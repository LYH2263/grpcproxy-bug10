package grpcproxy

// FrameMeta describes one proxied frame on the wire.
type FrameMeta struct {
	StreamID uint32
	Seq      uint64
	Flags    uint8
	Size     int
}

// SessionStats aggregates relay counters for one bridge session.
type SessionStats struct {
	FramesIn   int64
	FramesOut  int64
	BytesIn    int64
	BytesOut   int64
	Backlogged int64
	Dropped    int64
}

// Stats exposes kit-level counters.
type Stats struct {
	Sessions  int64
	Frames    int64
	Bytes     int64
	Blocked   int64
	Released  int64
}
