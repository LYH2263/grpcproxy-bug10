package frame

import (
	"context"
	"encoding/binary"
	"hash/crc32"
	"io"
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)

// Frame is one decoded unit on the proxy path.
type Frame struct {
	StreamID uint32
	Seq      uint64
	Flags    uint8
	Payload  []byte
}

// Codec encodes and decodes length-prefixed frames with CRC trailer.
type Codec struct {
	pool       *Pool
	maxPayload int
}

func NewCodec(pool *Pool, maxPayload int) *Codec {
	if pool == nil {
		pool = NewPool()
	}
	if maxPayload <= 0 {
		maxPayload = 64 << 10
	}
	return &Codec{pool: pool, maxPayload: maxPayload}
}

func (c *Codec) MaxPayload() int { return c.maxPayload }

// Encode writes one frame to w honoring ctx cancellation.
func (c *Codec) Encode(ctx context.Context, f *Frame, w io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if f == nil {
		return ErrNilFrame
	}
	if len(f.Payload) > c.maxPayload {
		return ErrPayloadTooLarge
	}
	hdr := Header{StreamID: f.StreamID, Seq: f.Seq, Flags: f.Flags, Length: uint32(len(f.Payload))}
	buf := make([]byte, HeaderLen+len(f.Payload)+4)
	EncodeHeader(buf[:HeaderLen], hdr)
	copy(buf[HeaderLen:], f.Payload)
	crc := crc32.Checksum(buf[:HeaderLen+len(f.Payload)], crcTable)
	binary.BigEndian.PutUint32(buf[HeaderLen+len(f.Payload):], crc)
	_, err := w.Write(buf)
	return err
}

// Decode reads one frame from r honoring ctx cancellation.
func (c *Codec) Decode(ctx context.Context, r io.Reader) (*Frame, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	hdrBuf := make([]byte, HeaderLen)
	if _, err := io.ReadFull(r, hdrBuf); err != nil {
		return nil, err
	}
	hdr, err := DecodeHeader(hdrBuf)
	if err != nil {
		return nil, err
	}
	if hdr.Length > uint32(c.maxPayload) {
		return nil, ErrPayloadTooLarge
	}
	payload := c.pool.Get(int(hdr.Length))
	if _, err := io.ReadFull(r, payload); err != nil {
		c.pool.Put(payload)
		return nil, err
	}
	crcBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, crcBuf); err != nil {
		c.pool.Put(payload)
		return nil, err
	}
	want := binary.BigEndian.Uint32(crcBuf)
	got := crc32.Checksum(append(hdrBuf, payload...), crcTable)
	if got != want {
		c.pool.Put(payload)
		return nil, ErrCorruptFrame
	}
	out := &Frame{StreamID: hdr.StreamID, Seq: hdr.Seq, Flags: hdr.Flags, Payload: payload}
	return out, nil
}

func (c *Codec) Release(f *Frame) {
	if f == nil {
		return
	}
	c.pool.Put(f.Payload)
	f.Payload = nil
}
