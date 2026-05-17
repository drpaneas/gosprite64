package gfx

import "encoding/binary"

// MarshalVtx serializes a Vtx to 16 big-endian bytes matching the N64 format.
func MarshalVtx(v Vtx) [16]byte {
	var buf [16]byte
	binary.BigEndian.PutUint16(buf[0:], uint16(v.X))
	binary.BigEndian.PutUint16(buf[2:], uint16(v.Y))
	binary.BigEndian.PutUint16(buf[4:], uint16(v.Z))
	binary.BigEndian.PutUint16(buf[6:], v.Flag)
	binary.BigEndian.PutUint16(buf[8:], uint16(v.S))
	binary.BigEndian.PutUint16(buf[10:], uint16(v.T))
	buf[12] = v.R
	buf[13] = v.G
	buf[14] = v.B
	buf[15] = v.A
	return buf
}

// MarshalVtxSlice serializes multiple vertices into a byte slice.
func MarshalVtxSlice(verts []Vtx) []byte {
	buf := make([]byte, len(verts)*16)
	for i, v := range verts {
		b := MarshalVtx(v)
		copy(buf[i*16:], b[:])
	}
	return buf
}

// MarshalN64Mtx serializes a math3d.N64Mtx (16 uint32s) to 64 big-endian bytes.
func MarshalN64Mtx(mtx [16]uint32) [64]byte {
	var buf [64]byte
	for i, w := range mtx {
		binary.BigEndian.PutUint32(buf[i*4:], w)
	}
	return buf
}

// MarshalViewport serializes a Viewport to 16 big-endian bytes.
func MarshalViewport(vp Viewport) [16]byte {
	var buf [16]byte
	binary.BigEndian.PutUint16(buf[0:], uint16(vp.ScaleX))
	binary.BigEndian.PutUint16(buf[2:], uint16(vp.ScaleY))
	binary.BigEndian.PutUint16(buf[4:], uint16(vp.ScaleZ))
	binary.BigEndian.PutUint16(buf[6:], uint16(vp.ScalePad))
	binary.BigEndian.PutUint16(buf[8:], uint16(vp.TransX))
	binary.BigEndian.PutUint16(buf[10:], uint16(vp.TransY))
	binary.BigEndian.PutUint16(buf[12:], uint16(vp.TransZ))
	binary.BigEndian.PutUint16(buf[14:], uint16(vp.TransPad))
	return buf
}
