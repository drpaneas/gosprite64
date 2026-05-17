package dma

import (
	"encoding/binary"
	"testing"
)

func TestDecompressMIO0InvalidMagic(t *testing.T) {
	header := make([]byte, 16)
	copy(header[0:4], "NOPE")
	_, err := DecompressMIO0(header)
	if err != ErrInvalidMIO0 {
		t.Errorf("expected ErrInvalidMIO0, got %v", err)
	}
}

func TestDecompressMIO0TooShort(t *testing.T) {
	_, err := DecompressMIO0([]byte("MIO0"))
	if err != ErrInvalidMIO0 {
		t.Errorf("expected ErrInvalidMIO0 for 4-byte input, got %v", err)
	}
}

func TestDecompressMIO0AllUncompressed(t *testing.T) {
	// Build a valid MIO0 blob where all 4 output bytes are uncompressed.
	// Layout: header (16 bytes) + 1 layout byte + 4 data bytes = 21 bytes total.
	// Layout byte 0xFF means all 8 bits are 1 (uncompressed); we only consume 4 bits.
	// comp_offset and uncomp_offset both point to byte 17.
	buf := make([]byte, 21)
	copy(buf[0:4], "MIO0")
	binary.BigEndian.PutUint32(buf[4:8], 4)  // decompressed size
	binary.BigEndian.PutUint32(buf[8:12], 17) // compressed data offset (empty)
	binary.BigEndian.PutUint32(buf[12:16], 17) // uncompressed data offset
	buf[16] = 0xFF                              // layout: all bits = uncompressed
	buf[17] = 0xDE
	buf[18] = 0xAD
	buf[19] = 0xBE
	buf[20] = 0xEF

	out, err := DecompressMIO0(buf)
	if err != nil {
		t.Fatalf("DecompressMIO0 returned error: %v", err)
	}
	expected := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	if len(out) != len(expected) {
		t.Fatalf("expected %d bytes, got %d", len(expected), len(out))
	}
	for i := range expected {
		if out[i] != expected[i] {
			t.Errorf("byte %d: expected 0x%02X, got 0x%02X", i, expected[i], out[i])
		}
	}
}
