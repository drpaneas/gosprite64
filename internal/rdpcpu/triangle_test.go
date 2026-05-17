package rdpcpu

import "testing"

func TestBuildTexturedTrianglePacketSize(t *testing.T) {
	packet := BuildTexturedTriangle(0, 0,
		TexVertex{X: 10, Y: 10, S: 0, T: 0, InvW: 1},
		TexVertex{X: 30, Y: 20, S: 16, T: 0, InvW: 1},
		TexVertex{X: 12, Y: 40, S: 0, T: 16, InvW: 1},
	)
	if len(packet) != 12 {
		t.Fatalf("packet len = %d, want 12", len(packet))
	}
}

func TestBuildTexturedTriangleOpcode(t *testing.T) {
	packet := BuildTexturedTriangle(3, 0,
		TexVertex{X: 10, Y: 10, S: 0, T: 0, InvW: 1},
		TexVertex{X: 30, Y: 20, S: 16, T: 0, InvW: 1},
		TexVertex{X: 12, Y: 40, S: 0, T: 16, InvW: 1},
	)
	opcode := uint8(packet[0] >> 56)
	if opcode != 0xCA && opcode != 0xCB {
		t.Fatalf("opcode = 0x%02X, want 0xCA or 0xCB", opcode)
	}
}

func TestBuildTexturedTriangleSortsByY(t *testing.T) {
	packetA := BuildTexturedTriangle(0, 0,
		TexVertex{X: 10, Y: 10, S: 0, T: 0, InvW: 1},
		TexVertex{X: 30, Y: 20, S: 16, T: 0, InvW: 1},
		TexVertex{X: 12, Y: 40, S: 0, T: 16, InvW: 1},
	)
	packetB := BuildTexturedTriangle(0, 0,
		TexVertex{X: 12, Y: 40, S: 0, T: 16, InvW: 1},
		TexVertex{X: 10, Y: 10, S: 0, T: 0, InvW: 1},
		TexVertex{X: 30, Y: 20, S: 16, T: 0, InvW: 1},
	)
	if len(packetA) != len(packetB) {
		t.Fatalf("packet len mismatch %d vs %d", len(packetA), len(packetB))
	}
	for i := range packetA {
		if packetA[i] != packetB[i] {
			t.Fatalf("packet[%d] mismatch: 0x%016x vs 0x%016x", i, packetA[i], packetB[i])
		}
	}
}

func TestBuildTexturedTriangleDefaultInvW(t *testing.T) {
	packetA := BuildTexturedTriangle(0, 0,
		TexVertex{X: 10, Y: 10, S: 0, T: 0},
		TexVertex{X: 30, Y: 20, S: 16, T: 0},
		TexVertex{X: 12, Y: 40, S: 0, T: 16},
	)
	packetB := BuildTexturedTriangle(0, 0,
		TexVertex{X: 10, Y: 10, S: 0, T: 0, InvW: 1},
		TexVertex{X: 30, Y: 20, S: 16, T: 0, InvW: 1},
		TexVertex{X: 12, Y: 40, S: 0, T: 16, InvW: 1},
	)
	for i := range packetA {
		if packetA[i] != packetB[i] {
			t.Fatalf("packet[%d] mismatch: 0x%016x vs 0x%016x", i, packetA[i], packetB[i])
		}
	}
}

func TestBuildTexturedTriangleNegativeCoordinates(t *testing.T) {
	packet := BuildTexturedTriangle(0, 0,
		TexVertex{X: -20, Y: -10, S: 0, T: 0, InvW: 1},
		TexVertex{X: 12, Y: 5, S: 16, T: 0, InvW: 1},
		TexVertex{X: 5, Y: 42, S: 0, T: 16, InvW: 1},
	)
	if len(packet) != 12 {
		t.Fatalf("packet len = %d, want 12", len(packet))
	}
}
