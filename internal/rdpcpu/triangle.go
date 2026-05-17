package rdpcpu

import "math"

const (
	opTriBase = 0xC8
	opTriTex  = 0x02
)

// TexVertex is the minimal vertex payload required for a textured triangle.
// S and T are specified in texel coordinates. InvW should normally be 1 for
// orthographic 2D use.
type TexVertex struct {
	X, Y float32
	S, T float32
	InvW float32
}

type edgeData struct {
	hx         float32
	hy         float32
	mx         float32
	my         float32
	fy         float32
	ish        float32
	attrFactor float32
}

// BuildTexturedTriangle assembles a raw RDP textured triangle packet.
// The returned slice contains 12 64-bit RDP words: 4 for edge setup and 8
// for texture coefficients.
func BuildTexturedTriangle(tile, mipmaps uint8, v1, v2, v3 TexVertex) []uint64 {
	v1, v2, v3 = sortByY(v1, v2, v3)

	args := make([]uint32, 0, 24)
	data := writeEdgeArgs(&args, tile, mipmaps, v1, v2, v3)
	writeTexArgs(&args, data, v1, v2, v3)

	out := make([]uint64, 0, len(args)/2)
	base := uint64(opTriBase | opTriTex)
	for i := 0; i < len(args); i += 2 {
		w0 := uint64(args[i])
		w1 := uint64(args[i+1])
		if i == 0 {
			out = append(out, (base<<56)|(w0<<32)|w1)
		} else {
			out = append(out, (w0<<32)|w1)
		}
	}
	return out
}

func sortByY(v1, v2, v3 TexVertex) (TexVertex, TexVertex, TexVertex) {
	if v1.Y > v2.Y {
		v1, v2 = v2, v1
	}
	if v2.Y > v3.Y {
		v2, v3 = v3, v2
	}
	if v1.Y > v2.Y {
		v1, v2 = v2, v1
	}
	return v1, v2, v3
}

func writeEdgeArgs(args *[]uint32, tile, mipmaps uint8, v1, v2, v3 TexVertex) edgeData {
	x1, x2, x3 := v1.X, v2.X, v3.X
	y1 := floorQuarter(v1.Y)
	y2 := floorQuarter(v2.Y)
	y3 := floorQuarter(v3.Y)

	y1f := clampSigned(int32(math.Floor(float64(v1.Y*4))), -4096*4, 4095*4)
	y2f := clampSigned(int32(math.Floor(float64(v2.Y*4))), -4096*4, 4095*4)
	y3f := clampSigned(int32(math.Floor(float64(v3.Y*4))), -4096*4, 4095*4)

	data := edgeData{
		hx: x3 - x1,
		hy: y3 - y1,
		mx: x2 - x1,
		my: y2 - y1,
	}
	lx := x3 - x2
	ly := y3 - y2

	nz := data.hx*data.my - data.hy*data.mx
	if math.Abs(float64(nz)) > math.SmallestNonzeroFloat32 {
		data.attrFactor = -1.0 / nz
	}
	lft := uint32(0)
	if nz < 0 {
		lft = 1
	}

	if math.Abs(float64(data.hy)) > math.SmallestNonzeroFloat32 {
		data.ish = data.hx / data.hy
	}
	var ism, isl float32
	if math.Abs(float64(data.my)) > math.SmallestNonzeroFloat32 {
		ism = data.mx / data.my
	}
	if math.Abs(float64(ly)) > math.SmallestNonzeroFloat32 {
		isl = lx / ly
	}
	data.fy = float32(math.Floor(float64(y1))) - y1

	xh := x1 + data.fy*data.ish
	xm := x1 + data.fy*ism
	xl := x2

	arg0 := carg(lft, 0x1, 23) |
		carg(uint32(mipmapLevel(mipmaps)), 0x7, 19) |
		carg(uint32(tile), 0x7, 16) |
		carg(uint32(y3f), 0x3FFF, 0)
	arg1 := carg(uint32(y2f), 0x3FFF, 16) | carg(uint32(y1f), 0x3FFF, 0)

	*args = append(*args,
		arg0, arg1,
		uint32(floatToS16_16(xl)), uint32(floatToS16_16(isl)),
		uint32(floatToS16_16(xh)), uint32(floatToS16_16(data.ish)),
		uint32(floatToS16_16(xm)), uint32(floatToS16_16(ism)),
	)

	return data
}

func writeTexArgs(args *[]uint32, data edgeData, v1, v2, v3 TexVertex) {
	s1, t1, invw1 := v1.S*32, v1.T*32, defaultInvW(v1.InvW)
	s2, t2, invw2 := v2.S*32, v2.T*32, defaultInvW(v2.InvW)
	s3, t3, invw3 := v3.S*32, v3.T*32, defaultInvW(v3.InvW)

	minw := 1.0 / max3(invw1, invw2, invw3)
	invw1 *= minw
	invw2 *= minw
	invw3 *= minw

	s1 *= invw1
	t1 *= invw1
	s2 *= invw2
	t2 *= invw2
	s3 *= invw3
	t3 *= invw3

	invw1 *= 0x7FFF
	invw2 *= 0x7FFF
	invw3 *= 0x7FFF

	ms, mt, mw := s2-s1, t2-t1, invw2-invw1
	hs, ht, hw := s3-s1, t3-t1, invw3-invw1

	nxS := data.hy*ms - data.my*hs
	nxT := data.hy*mt - data.my*ht
	nxW := data.hy*mw - data.my*hw
	nyS := data.mx*hs - data.hx*ms
	nyT := data.mx*ht - data.hx*mt
	nyW := data.mx*hw - data.hx*mw

	dsDx := nxS * data.attrFactor
	dtDx := nxT * data.attrFactor
	dwDx := nxW * data.attrFactor
	dsDy := nyS * data.attrFactor
	dtDy := nyT * data.attrFactor
	dwDy := nyW * data.attrFactor

	dsDe := dsDy + dsDx*data.ish
	dtDe := dtDy + dtDx*data.ish
	dwDe := dwDy + dwDx*data.ish

	finalS := floatToS16_16(s1 + data.fy*dsDe)
	finalT := floatToS16_16(t1 + data.fy*dtDe)
	finalW := floatToS16_16(invw1 + data.fy*dwDe)

	dsDxFixed := floatToS16_16(dsDx)
	dtDxFixed := floatToS16_16(dtDx)
	dwDxFixed := floatToS16_16(dwDx)
	dsDeFixed := floatToS16_16(dsDe)
	dtDeFixed := floatToS16_16(dtDe)
	dwDeFixed := floatToS16_16(dwDe)
	dsDyFixed := floatToS16_16(dsDy)
	dtDyFixed := floatToS16_16(dtDy)
	dwDyFixed := floatToS16_16(dwDy)

	*args = append(*args,
		(uint32(finalS)&0xffff0000)|(0xffff&(uint32(finalT)>>16)),
		(uint32(finalW) & 0xffff0000),
		(uint32(dsDxFixed)&0xffff0000)|(0xffff&(uint32(dtDxFixed)>>16)),
		(uint32(dwDxFixed) & 0xffff0000),
		(uint32(finalS)<<16)|(uint32(finalT)&0xffff),
		uint32(finalW) << 16,
		(uint32(dsDxFixed)<<16)|(uint32(dtDxFixed)&0xffff),
		uint32(dwDxFixed) << 16,
		(uint32(dsDeFixed)&0xffff0000)|(0xffff&(uint32(dtDeFixed)>>16)),
		(uint32(dwDeFixed) & 0xffff0000),
		(uint32(dsDyFixed)&0xffff0000)|(0xffff&(uint32(dtDyFixed)>>16)),
		(uint32(dwDyFixed) & 0xffff0000),
		(uint32(dsDeFixed)<<16)|(uint32(dtDeFixed)&0xffff),
		uint32(dwDeFixed) << 16,
		(uint32(dsDyFixed)<<16)|(uint32(dtDyFixed)&0xffff),
		uint32(dwDyFixed) << 16,
	)
}

func floorQuarter(v float32) float32 {
	return float32(math.Floor(float64(v*4))) / 4
}

func floatToS16_16(f float32) int32 {
	if f >= 32768 {
		return 0x7FFFFFFF
	}
	if f < -32768 {
		return int32(-2147483648)
	}
	return int32(math.Floor(float64(f * 65536)))
}

func carg(v uint32, mask uint32, shift uint32) uint32 {
	return (v & mask) << shift
}

func mipmapLevel(levels uint8) uint8 {
	if levels == 0 {
		return 0
	}
	return levels - 1
}

func clampSigned(v, lo, hi int32) int32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func defaultInvW(v float32) float32 {
	if v == 0 {
		return 1
	}
	return v
}

func max3(a, b, c float32) float32 {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	return a
}
