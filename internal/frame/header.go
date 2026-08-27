package frame

import "encoding/binary"

const (
	Magic     = "GPXY"
	Version   = 1
	HeaderLen = 4 + 1 + 4 + 8 + 1 + 4
)

const (
	FlagData  uint8 = 1 << 0
	FlagFin   uint8 = 1 << 1
	FlagCont  uint8 = 1 << 2
)

// Header is the fixed-size prefix of a wire frame.
type Header struct {
	StreamID uint32
	Seq      uint64
	Flags    uint8
	Length   uint32
}

func (h Header) WireSize() int { return HeaderLen }

func EncodeHeader(buf []byte, h Header) {
	copy(buf[0:4], Magic)
	buf[4] = Version
	binary.BigEndian.PutUint32(buf[5:9], h.StreamID)
	binary.BigEndian.PutUint64(buf[9:17], h.Seq)
	buf[17] = h.Flags
	binary.BigEndian.PutUint32(buf[18:22], h.Length)
}

func DecodeHeader(buf []byte) (Header, error) {
	if len(buf) < HeaderLen {
		return Header{}, ErrShortHeader
	}
	if string(buf[0:4]) != Magic {
		return Header{}, ErrBadMagic
	}
	if buf[4] != Version {
		return Header{}, ErrBadVersion
	}
	return Header{
		StreamID: binary.BigEndian.Uint32(buf[5:9]),
		Seq:      binary.BigEndian.Uint64(buf[9:17]),
		Flags:    buf[17],
		Length:   binary.BigEndian.Uint32(buf[18:22]),
	}, nil
}
