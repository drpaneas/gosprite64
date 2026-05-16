package format

import (
	"encoding/binary"
	"fmt"
)

const headerSize = 16

type Header struct {
	Magic        [4]byte
	Version      uint16
	Flags        uint16
	HeaderBytes  uint32
	PayloadBytes uint32
}

func ParseHeader(raw []byte, want string) (Header, error) {
	var h Header
	if len(raw) < headerSize {
		return h, fmt.Errorf("format: header too short: got %d bytes", len(raw))
	}

	copy(h.Magic[:], raw[:4])
	if string(h.Magic[:]) != want {
		return h, fmt.Errorf("format: bad magic %q, want %q", string(h.Magic[:]), want)
	}

	h.Version = binary.LittleEndian.Uint16(raw[4:6])
	h.Flags = binary.LittleEndian.Uint16(raw[6:8])
	h.HeaderBytes = binary.LittleEndian.Uint32(raw[8:12])
	h.PayloadBytes = binary.LittleEndian.Uint32(raw[12:16])

	if h.Version != 1 {
		return h, fmt.Errorf("format: unsupported version %d", h.Version)
	}
	if h.HeaderBytes < headerSize {
		return h, fmt.Errorf("format: invalid header size %d", h.HeaderBytes)
	}
	if int(h.HeaderBytes) > len(raw) {
		return h, fmt.Errorf("format: header size %d exceeds input %d", h.HeaderBytes, len(raw))
	}

	payloadEnd := int(h.HeaderBytes + h.PayloadBytes)
	if payloadEnd > len(raw) {
		return h, fmt.Errorf("format: payload size %d exceeds input %d", h.PayloadBytes, len(raw))
	}

	return h, nil
}
