package gfx

import "testing"

func TestDisplayListPipeSync(t *testing.T) {
	dl := NewDisplayList(8)
	dl.DPPipeSync()
	if dl.Len() != 1 {
		t.Fatalf("expected 1 command, got %d", dl.Len())
	}
	cmd := dl.Commands()[0]
	opcode := cmd.W0 >> 24
	if opcode != OpDPPipeSync {
		t.Errorf("opcode = 0x%02X, want 0x%02X", opcode, OpDPPipeSync)
	}
	if cmd.W1 != 0 {
		t.Errorf("W1 = 0x%08X, want 0", cmd.W1)
	}
}

func TestDisplayListFullSync(t *testing.T) {
	dl := NewDisplayList(8)
	dl.DPFullSync()
	if dl.Len() != 1 {
		t.Fatalf("expected 1 command, got %d", dl.Len())
	}
	cmd := dl.Commands()[0]
	opcode := cmd.W0 >> 24
	if opcode != OpDPFullSync {
		t.Errorf("opcode = 0x%02X, want 0x%02X", opcode, OpDPFullSync)
	}
}

func TestDisplayListTriangle(t *testing.T) {
	dl := NewDisplayList(8)
	dl.SP1Triangle(1, 2, 3, 0)
	if dl.Len() != 1 {
		t.Fatalf("expected 1 command, got %d", dl.Len())
	}
	cmd := dl.Commands()[0]
	opcode := cmd.W0 >> 24
	if opcode != OpSP1Triangle {
		t.Errorf("opcode = 0x%02X, want 0x%02X", opcode, OpSP1Triangle)
	}
	// Vertex indices are encoded as v*10.
	v0 := uint8((cmd.W1 >> 16) & 0xFF)
	v1 := uint8((cmd.W1 >> 8) & 0xFF)
	v2 := uint8(cmd.W1 & 0xFF)
	if v0 != 1*10 {
		t.Errorf("v0 = %d, want %d", v0, 1*10)
	}
	if v1 != 2*10 {
		t.Errorf("v1 = %d, want %d", v1, 2*10)
	}
	if v2 != 3*10 {
		t.Errorf("v2 = %d, want %d", v2, 3*10)
	}
}

func TestDisplayListReset(t *testing.T) {
	dl := NewDisplayList(8)
	dl.DPPipeSync()
	dl.DPFullSync()
	if dl.Len() != 2 {
		t.Fatalf("expected 2 commands, got %d", dl.Len())
	}
	dl.Reset()
	if dl.Len() != 0 {
		t.Errorf("after Reset, Len() = %d, want 0", dl.Len())
	}
}

func TestDisplayListSetColorImage(t *testing.T) {
	dl := NewDisplayList(8)
	addr := uint32(0x00100000)
	dl.DPSetColorImage(FmtRGBA, Siz16b, 320, addr)
	if dl.Len() != 1 {
		t.Fatalf("expected 1 command, got %d", dl.Len())
	}
	cmd := dl.Commands()[0]
	opcode := cmd.W0 >> 24
	if opcode != OpDPSetColorImage {
		t.Errorf("opcode = 0x%02X, want 0x%02X", opcode, OpDPSetColorImage)
	}
	if cmd.W1 != addr {
		t.Errorf("W1 = 0x%08X, want 0x%08X", cmd.W1, addr)
	}
}

func TestDisplayListRawPacket(t *testing.T) {
	dl := NewDisplayList(4)
	dl.DPRaw(0xc800000000000001, 0x123456789abcdef0)
	if dl.Len() != 1 {
		t.Fatalf("expected 1 raw packet entry, got %d", dl.Len())
	}
	cmd := dl.Commands()[0]
	if len(cmd.Raw) != 2 {
		t.Fatalf("raw packet len = %d, want 2", len(cmd.Raw))
	}
	if cmd.Raw[0] != 0xc800000000000001 {
		t.Fatalf("raw[0] = 0x%016x, want 0xc800000000000001", cmd.Raw[0])
	}
	if cmd.Raw[1] != 0x123456789abcdef0 {
		t.Fatalf("raw[1] = 0x%016x, want 0x123456789abcdef0", cmd.Raw[1])
	}
}
