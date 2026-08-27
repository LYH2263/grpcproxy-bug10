package proxy

import (
	"context"
	"io"

	"github.com/LYH2263/go-grpcproxy/internal/frame"
)

// Relay copies framed messages from src to dst through handler chain.
type Relay struct {
	bridge *Bridge
	sess   *Session
}

func NewRelay(b *Bridge, sess *Session) *Relay {
	return &Relay{bridge: b, sess: sess}
}

// Forward reads frames from src, invokes chain, writes to dst.
func (r *Relay) Forward(ctx context.Context, src io.Reader, dst io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	fr := frame.NewReader(r.bridge.codec, src)
	fw := frame.NewWriter(r.bridge.codec, dst)
	for {
		f, err := fr.ReadFrame(ctx)
		if err != nil {
			return err
		}
		r.sess.addIn(int64(len(f.Payload)))
		if r.bridge.chain != nil {
			if err := r.bridge.chain.Invoke(ctx, f); err != nil {
				r.bridge.codec.Release(f)
				return err
			}
		}
		if err := fw.WriteFrame(ctx, f); err != nil {
			r.bridge.codec.Release(f)
			return err
		}
		r.sess.addOut(int64(len(f.Payload)))
		r.bridge.codec.Release(f)
	}
}

// EncodeAndForward encodes raw payload as one frame and forwards.
func (r *Relay) EncodeAndForward(ctx context.Context, streamID uint32, payload []byte, dst io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f := &frame.Frame{
		StreamID: streamID,
		Seq:      r.sess.NextSeq(),
		Flags:    frame.FlagData,
		Payload:  append([]byte(nil), payload...),
	}
	if r.bridge.chain != nil {
		if err := r.bridge.chain.Invoke(ctx, f); err != nil {
			return err
		}
	}
	w := frame.NewWriter(r.bridge.codec, dst)
	if err := w.WriteFrame(ctx, f); err != nil {
		return err
	}
	r.sess.addOut(int64(len(f.Payload)))
	return nil
}
