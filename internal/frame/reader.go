package frame

import (
	"context"
	"io"
)

// Reader wraps Codec.Decode for streaming reads.
type Reader struct {
	codec *Codec
	src   io.Reader
}

func NewReader(codec *Codec, src io.Reader) *Reader {
	return &Reader{codec: codec, src: src}
}

// ReadFrame decodes the next frame.
func (r *Reader) ReadFrame(ctx context.Context) (*Frame, error) {
	return r.codec.Decode(ctx, r.src)
}

// ReadFrameLimited reads up to maxBytes total payload across frames.
func (r *Reader) ReadFrameLimited(ctx context.Context, maxBytes int) (*Frame, error) {
	f, err := r.codec.Decode(ctx, r.src)
	if err != nil {
		return nil, err
	}
	if len(f.Payload) > maxBytes {
		r.codec.Release(f)
		return nil, ErrPayloadTooLarge
	}
	return f, nil
}
