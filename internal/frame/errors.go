package frame

import "errors"

var (
	ErrShortHeader     = errors.New("frame: short header")
	ErrBadMagic        = errors.New("frame: bad magic")
	ErrBadVersion      = errors.New("frame: bad version")
	ErrNilFrame        = errors.New("frame: nil frame")
	ErrPayloadTooLarge = errors.New("frame: payload too large")
	ErrCorruptFrame    = errors.New("frame: corrupt frame checksum")
)
