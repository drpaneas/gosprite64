package dma

import (
	"encoding/binary"
	"errors"
)

var ErrInvalidMIO0 = errors.New("dma: invalid MIO0 data")

// DecompressMIO0 decompresses MIO0-compressed data, the standard compression
// format used in many N64 games (including SM64).
//
// MIO0 format:
//   - 4 bytes: "MIO0" magic
//   - 4 bytes: decompressed size (big-endian)
//   - 4 bytes: offset to compressed data
//   - 4 bytes: offset to uncompressed data
//   - variable: layout bits, compressed data, uncompressed data
func DecompressMIO0(src []byte) ([]byte, error) {
	if len(src) < 16 {
		return nil, ErrInvalidMIO0
	}
	if string(src[0:4]) != "MIO0" {
		return nil, ErrInvalidMIO0
	}

	decompSize := binary.BigEndian.Uint32(src[4:8])
	compOffset := binary.BigEndian.Uint32(src[8:12])
	uncompOffset := binary.BigEndian.Uint32(src[12:16])

	dst := make([]byte, decompSize)
	dstPos := 0

	layoutBits := src[16:]
	compData := src[compOffset:]
	uncompData := src[uncompOffset:]

	layoutIdx := 0
	layoutBit := 0
	compIdx := 0
	uncompIdx := 0

	for dstPos < int(decompSize) {
		if layoutIdx >= len(layoutBits) {
			return nil, ErrInvalidMIO0
		}
		bit := (layoutBits[layoutIdx] >> (7 - layoutBit)) & 1
		layoutBit++
		if layoutBit >= 8 {
			layoutBit = 0
			layoutIdx++
		}

		if bit != 0 {
			if uncompIdx >= len(uncompData) {
				return nil, ErrInvalidMIO0
			}
			dst[dstPos] = uncompData[uncompIdx]
			uncompIdx++
			dstPos++
		} else {
			if compIdx+1 >= len(compData) {
				return nil, ErrInvalidMIO0
			}
			pair := uint16(compData[compIdx])<<8 | uint16(compData[compIdx+1])
			compIdx += 2

			length := int((pair>>12)&0x0F) + 3
			offset := int(pair&0x0FFF) + 1

			for i := 0; i < length; i++ {
				if dstPos >= int(decompSize) {
					break
				}
				dst[dstPos] = dst[dstPos-offset]
				dstPos++
			}
		}
	}

	return dst, nil
}
