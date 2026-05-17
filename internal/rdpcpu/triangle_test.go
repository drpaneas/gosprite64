package rdpcpu

import "testing"

func TestFillTriangleCommandCount(t *testing.T) {
	cmds := FillTriangle([2]float32{160, 40}, [2]float32{80, 200}, [2]float32{240, 200})
	if len(cmds) != 4 {
		t.Fatalf("expected 4 words (edge coeffs), got %d", len(cmds))
	}
}

func TestFillTriangleOpcode(t *testing.T) {
	cmds := FillTriangle([2]float32{160, 40}, [2]float32{80, 200}, [2]float32{240, 200})
	opcode := (cmds[0] >> 56) & 0xFF
	if opcode != 0x08 {
		t.Fatalf("expected opcode 0x08 (fill tri), got 0x%02X", opcode)
	}
}

func TestFillTriangleVertexSorting(t *testing.T) {
	// Regardless of input order, Y coords in the command should be sorted:
	// y1 (top) <= y2 (mid) <= y3 (bottom)
	cmds1 := FillTriangle([2]float32{160, 40}, [2]float32{80, 200}, [2]float32{240, 120})
	cmds2 := FillTriangle([2]float32{240, 120}, [2]float32{160, 40}, [2]float32{80, 200})
	cmds3 := FillTriangle([2]float32{80, 200}, [2]float32{240, 120}, [2]float32{160, 40})

	// All should produce the same edge coefficients regardless of input order
	for i := range cmds1 {
		if cmds1[i] != cmds2[i] || cmds1[i] != cmds3[i] {
			t.Fatalf("word %d: different results for different vertex orders: %016X vs %016X vs %016X",
				i, cmds1[i], cmds2[i], cmds3[i])
		}
	}
}

func TestShadeTriangleCommandCount(t *testing.T) {
	white := [4]float32{1, 1, 1, 1}
	cmds := ShadeTriangle(
		[2]float32{160, 40}, [2]float32{80, 200}, [2]float32{240, 200},
		white, white, white,
	)
	// 4 edge words + 8 shade words = 12
	if len(cmds) != 12 {
		t.Fatalf("expected 12 words (edge + shade), got %d", len(cmds))
	}
}

func TestShadeTriangleOpcode(t *testing.T) {
	red := [4]float32{1, 0, 0, 1}
	green := [4]float32{0, 1, 0, 1}
	blue := [4]float32{0, 0, 1, 1}
	cmds := ShadeTriangle(
		[2]float32{160, 40}, [2]float32{80, 200}, [2]float32{240, 200},
		red, green, blue,
	)
	opcode := (cmds[0] >> 56) & 0xFF
	if opcode != 0x0C {
		t.Fatalf("expected opcode 0x0C (shade tri), got 0x%02X", opcode)
	}
}

func TestFillTriangleDeterministic(t *testing.T) {
	v1 := [2]float32{100, 50}
	v2 := [2]float32{50, 150}
	v3 := [2]float32{200, 180}
	cmds1 := FillTriangle(v1, v2, v3)
	cmds2 := FillTriangle(v1, v2, v3)
	for i := range cmds1 {
		if cmds1[i] != cmds2[i] {
			t.Fatalf("non-deterministic output at word %d", i)
		}
	}
}
