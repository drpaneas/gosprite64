package rdpcpu

import "math"

// Ported from libdragon's rdpq_tri.c (rdpq_triangle_cpu path).
// Computes RDP edge coefficients on the CPU and packs them into
// raw 64-bit words that can be fed directly to rdp.RDP.PushRaw().

func floatToS16_16(f float32) int32 {
	if f >= 32768.0 {
		return 0x7FFFFFFF
	}
	if f < -32768.0 {
		return -0x80000000
	}
	return int32(math.Floor(float64(f) * 65536.0))
}

func clampI32(v, lo, hi int32) int32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

type edgeData struct {
	hx, hy, mx, my float32
	fy             float32
	ish            float32
	attrFactor     float32
}

func packU64(hi, lo uint32) uint64 {
	return uint64(hi)<<32 | uint64(lo)
}

// FillTriangle computes edge coefficients for a flat-colored triangle
// and returns the raw RDP command words. The RDP opcode is 0x08 (fill tri).
// Vertices are screen-space: v[0]=X, v[1]=Y.
func FillTriangle(v1, v2, v3 [2]float32) []uint64 {
	if v1[1] > v2[1] {
		v1, v2 = v2, v1
	}
	if v2[1] > v3[1] {
		v2, v3 = v3, v2
	}
	if v1[1] > v2[1] {
		v1, v2 = v2, v1
	}

	var data edgeData
	cmds := writeEdgeCoeffs(&data, 0x08, 0, 0, v1, v2, v3)
	return cmds
}

// ShadeTriangle computes edge + shade coefficients for a Gouraud-shaded
// triangle. Vertices: v[0]=X, v[1]=Y. Colors: c[0..3]=R,G,B,A as 0.0-1.0.
// RDP opcode 0x0C (shade tri).
func ShadeTriangle(v1, v2, v3 [2]float32, c1, c2, c3 [4]float32) []uint64 {
	if v1[1] > v2[1] {
		v1, v2 = v2, v1
		c1, c2 = c2, c1
	}
	if v2[1] > v3[1] {
		v2, v3 = v3, v2
		c2, c3 = c3, c2
	}
	if v1[1] > v2[1] {
		v1, v2 = v2, v1
		c1, c2 = c2, c1
	}

	var data edgeData
	cmds := writeEdgeCoeffs(&data, 0x0C, 0, 0, v1, v2, v3)
	cmds = append(cmds, writeShadeCoeffs(&data, c1, c2, c3)...)
	return cmds
}

func writeEdgeCoeffs(data *edgeData, cmdID uint8, tile, mipmaps uint8, v1, v2, v3 [2]float32) []uint64 {
	x1, x2, x3 := v1[0], v2[0], v3[0]
	y1 := float32(math.Floor(float64(v1[1])*4)) / 4
	y2 := float32(math.Floor(float64(v2[1])*4)) / 4
	y3 := float32(math.Floor(float64(v3[1])*4)) / 4

	toFixed := float32(4.0)
	y1f := clampI32(int32(math.Floor(float64(v1[1]*toFixed))), -4096*4, 4095*4)
	y2f := clampI32(int32(math.Floor(float64(v2[1]*toFixed))), -4096*4, 4095*4)
	y3f := clampI32(int32(math.Floor(float64(v3[1]*toFixed))), -4096*4, 4095*4)

	data.hx = x3 - x1
	data.hy = y3 - y1
	data.mx = x2 - x1
	data.my = y2 - y1
	lx := x3 - x2
	ly := y3 - y2

	nz := data.hx*data.my - data.hy*data.mx
	if math.Abs(float64(nz)) > math.SmallestNonzeroFloat32 {
		data.attrFactor = -1.0 / nz
	} else {
		data.attrFactor = 0
	}
	var lft uint32
	if nz < 0 {
		lft = 1
	}

	if math.Abs(float64(data.hy)) > math.SmallestNonzeroFloat32 {
		data.ish = data.hx / data.hy
	} else {
		data.ish = 0
	}
	var ism float32
	if math.Abs(float64(data.my)) > math.SmallestNonzeroFloat32 {
		ism = data.mx / data.my
	}
	var isl float32
	if math.Abs(float64(ly)) > math.SmallestNonzeroFloat32 {
		isl = lx / ly
	}
	data.fy = float32(math.Floor(float64(y1))) - y1

	xh := x1 + data.fy*data.ish
	xm := x1 + data.fy*ism
	xl := x2

	w0 := uint32(cmdID)<<24 | (lft&0x1)<<23 |
		uint32(mipmaps&0x7)<<19 | uint32(tile&0x7)<<16 |
		uint32(y3f)&0x3FFF
	w1 := (uint32(y2f)&0x3FFF)<<16 | uint32(y1f)&0x3FFF

	return []uint64{
		packU64(w0, w1),
		packU64(uint32(floatToS16_16(xl)), uint32(floatToS16_16(isl))),
		packU64(uint32(floatToS16_16(xh)), uint32(floatToS16_16(data.ish))),
		packU64(uint32(floatToS16_16(xm)), uint32(floatToS16_16(ism))),
	}
}

// TexVertex holds screen-space position and texture coordinates for a
// textured triangle vertex.
type TexVertex struct {
	X, Y float32
	S, T float32
	InvW float32
}

// BuildTexturedTriangle computes edge + texture coefficients for an
// RDP textured triangle. Returns raw 64-bit command words (opcode 0x0A).
// tileIdx and mipmaps correspond to the RDP tile descriptor to use.
func BuildTexturedTriangle(tileIdx, mipmaps uint8, v1, v2, v3 TexVertex) []uint64 {
	p1 := [2]float32{v1.X, v1.Y}
	p2 := [2]float32{v2.X, v2.Y}
	p3 := [2]float32{v3.X, v3.Y}
	t1 := [3]float32{v1.S, v1.T, v1.InvW}
	t2 := [3]float32{v2.S, v2.T, v2.InvW}
	t3 := [3]float32{v3.S, v3.T, v3.InvW}

	if p1[1] > p2[1] {
		p1, p2 = p2, p1
		t1, t2 = t2, t1
	}
	if p2[1] > p3[1] {
		p2, p3 = p3, p2
		t2, t3 = t3, t2
	}
	if p1[1] > p2[1] {
		p1, p2 = p2, p1
		t1, t2 = t2, t1
	}

	var data edgeData
	cmds := writeEdgeCoeffs(&data, 0x0A, tileIdx, mipmaps, p1, p2, p3)
	cmds = append(cmds, writeTexCoeffs(&data, t1, t2, t3)...)
	return cmds
}

func writeTexCoeffs(data *edgeData, t1, t2, t3 [3]float32) []uint64 {
	s1 := t1[0] * 32
	tt1 := t1[1] * 32
	invw1 := t1[2]
	s2 := t2[0] * 32
	tt2 := t2[1] * 32
	invw2 := t2[2]
	s3 := t3[0] * 32
	tt3 := t3[1] * 32
	invw3 := t3[2]

	maxW := invw1
	if invw2 > maxW {
		maxW = invw2
	}
	if invw3 > maxW {
		maxW = invw3
	}
	minW := float32(1.0)
	if maxW > 0 {
		minW = 1.0 / maxW
	}

	invw1 *= minW
	invw2 *= minW
	invw3 *= minW
	s1 *= invw1
	tt1 *= invw1
	s2 *= invw2
	tt2 *= invw2
	s3 *= invw3
	tt3 *= invw3
	invw1 *= 0x7FFF
	invw2 *= 0x7FFF
	invw3 *= 0x7FFF

	ms := s2 - s1
	mt := tt2 - tt1
	mw := invw2 - invw1
	hs := s3 - s1
	ht := tt3 - tt1
	hw := invw3 - invw1

	nxS := data.hy*ms - data.my*hs
	nxT := data.hy*mt - data.my*ht
	nxW := data.hy*mw - data.my*hw
	nyS := data.mx*hs - data.hx*ms
	nyT := data.mx*ht - data.hx*mt
	nyW := data.mx*hw - data.hx*mw

	DsDx := nxS * data.attrFactor
	DtDx := nxT * data.attrFactor
	DwDx := nxW * data.attrFactor
	DsDy := nyS * data.attrFactor
	DtDy := nyT * data.attrFactor
	DwDy := nyW * data.attrFactor

	DsDe := DsDy + DsDx*data.ish
	DtDe := DtDy + DtDx*data.ish
	DwDe := DwDy + DwDx*data.ish

	finalS := floatToS16_16(s1 + data.fy*DsDe)
	finalT := floatToS16_16(tt1 + data.fy*DtDe)
	finalW := floatToS16_16(invw1 + data.fy*DwDe)

	DsDxF := floatToS16_16(DsDx)
	DtDxF := floatToS16_16(DtDx)
	DwDxF := floatToS16_16(DwDx)
	DsDeF := floatToS16_16(DsDe)
	DtDeF := floatToS16_16(DtDe)
	DwDeF := floatToS16_16(DwDe)
	DsDyF := floatToS16_16(DsDy)
	DtDyF := floatToS16_16(DtDy)
	DwDyF := floatToS16_16(DwDy)

	hi := func(v int32) uint32 { return uint32(v) & 0xFFFF0000 }
	lo := func(v int32) uint32 { return uint32(v) & 0x0000FFFF }
	shr16 := func(v int32) uint32 { return uint32(v>>16) & 0xFFFF }
	shl16 := func(v int32) uint32 { return uint32(v<<16) & 0xFFFF0000 }

	return []uint64{
		packU64(hi(finalS)|shr16(finalT), hi(finalW)),
		packU64(hi(DsDxF)|shr16(DtDxF), hi(DwDxF)),
		packU64(shl16(finalS)|lo(finalT), shl16(finalW)),
		packU64(shl16(DsDxF)|lo(DtDxF), shl16(DwDxF)),
		packU64(hi(DsDeF)|shr16(DtDeF), hi(DwDeF)),
		packU64(hi(DsDyF)|shr16(DtDyF), hi(DwDyF)),
		packU64(shl16(DsDeF)|lo(DtDeF), shl16(DwDeF)),
		packU64(shl16(DsDyF)|lo(DtDyF), shl16(DwDyF)),
	}
}

func writeShadeCoeffs(data *edgeData, c1, c2, c3 [4]float32) []uint64 {
	mr := (c2[0] - c1[0]) * 255
	mg := (c2[1] - c1[1]) * 255
	mb := (c2[2] - c1[2]) * 255
	ma := (c2[3] - c1[3]) * 255
	hr := (c3[0] - c1[0]) * 255
	hg := (c3[1] - c1[1]) * 255
	hb := (c3[2] - c1[2]) * 255
	ha := (c3[3] - c1[3]) * 255

	nxR := data.hy*mr - data.my*hr
	nxG := data.hy*mg - data.my*hg
	nxB := data.hy*mb - data.my*hb
	nxA := data.hy*ma - data.my*ha
	nyR := data.mx*hr - data.hx*mr
	nyG := data.mx*hg - data.hx*mg
	nyB := data.mx*hb - data.hx*mb
	nyA := data.mx*ha - data.hx*ma

	DrDx := nxR * data.attrFactor
	DgDx := nxG * data.attrFactor
	DbDx := nxB * data.attrFactor
	DaDx := nxA * data.attrFactor
	DrDy := nyR * data.attrFactor
	DgDy := nyG * data.attrFactor
	DbDy := nyB * data.attrFactor
	DaDy := nyA * data.attrFactor

	DrDe := DrDy + DrDx*data.ish
	DgDe := DgDy + DgDx*data.ish
	DbDe := DbDy + DbDx*data.ish
	DaDe := DaDy + DaDx*data.ish

	finalR := floatToS16_16(c1[0]*255 + data.fy*DrDe)
	finalG := floatToS16_16(c1[1]*255 + data.fy*DgDe)
	finalB := floatToS16_16(c1[2]*255 + data.fy*DbDe)
	finalA := floatToS16_16(c1[3]*255 + data.fy*DaDe)

	drDxF := floatToS16_16(DrDx)
	dgDxF := floatToS16_16(DgDx)
	dbDxF := floatToS16_16(DbDx)
	daDxF := floatToS16_16(DaDx)
	drDeF := floatToS16_16(DrDe)
	dgDeF := floatToS16_16(DgDe)
	dbDeF := floatToS16_16(DbDe)
	daDeF := floatToS16_16(DaDe)
	drDyF := floatToS16_16(DrDy)
	dgDyF := floatToS16_16(DgDy)
	dbDyF := floatToS16_16(DbDy)
	daDyF := floatToS16_16(DaDy)

	hi := func(v int32) uint32 { return uint32(v) & 0xFFFF0000 }
	lo := func(v int32) uint32 { return uint32(v) & 0x0000FFFF }
	shr16 := func(v int32) uint32 { return uint32(v>>16) & 0xFFFF }
	shl16 := func(v int32) uint32 { return uint32(v<<16) & 0xFFFF0000 }

	return []uint64{
		packU64(hi(finalR)|shr16(finalG), hi(finalB)|shr16(finalA)),
		packU64(hi(drDxF)|shr16(dgDxF), hi(dbDxF)|shr16(daDxF)),
		packU64(shl16(finalR)|lo(finalG), shl16(finalB)|lo(finalA)),
		packU64(shl16(drDxF)|lo(dgDxF), shl16(dbDxF)|lo(daDxF)),
		packU64(hi(drDeF)|shr16(dgDeF), hi(dbDeF)|shr16(daDeF)),
		packU64(hi(drDyF)|shr16(dgDyF), hi(dbDyF)|shr16(daDyF)),
		packU64(shl16(drDeF)|lo(dgDeF), shl16(dbDeF)|lo(daDeF)),
		packU64(shl16(drDyF)|lo(dgDyF), shl16(dbDyF)|lo(daDyF)),
	}
}
