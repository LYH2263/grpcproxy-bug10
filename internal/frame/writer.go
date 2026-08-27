package frame

import (
	"context"
	"io"
)

// Writer wraps Codec.Encode for streaming writes.
type Writer struct {
	codec *Codec
	dst   io.Writer
}

func NewWriter(codec *Codec, dst io.Writer) *Writer {
	return &Writer{codec: codec, dst: dst}
}

func (w *Writer) WriteFrame(ctx context.Context, f *Frame) error {
	return w.codec.Encode(ctx, f, w.dst)
}
