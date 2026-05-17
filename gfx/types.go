package gfx

// Gfx is a single 64-bit display list command, matching the N64's Gfx union.
type Gfx struct {
	W0, W1 uint32
	Raw    []uint64
}

// Vtx is an N64 vertex with position, texture coords, and color/normal.
// Matches the N64 Vtx_t struct (16 bytes).
type Vtx struct {
	X, Y, Z int16
	Flag    uint16
	S, T    int16 // texture coordinates
	R, G, B, A uint8
}

// VtxN is an N64 vertex with normals instead of colors.
// Matches the N64 Vtx_tn struct (16 bytes).
type VtxN struct {
	X, Y, Z int16
	Flag    uint16
	S, T    int16
	NX, NY, NZ int8
	A          uint8
}

// Viewport matches the N64 Vp_t used by SPViewport.
type Viewport struct {
	ScaleX, ScaleY, ScaleZ, ScalePad int16
	TransX, TransY, TransZ, TransPad int16
}

// NewViewport creates a viewport for the given screen dimensions.
func NewViewport(width, height int) Viewport {
	return Viewport{
		ScaleX: int16(width / 2 * 4),
		ScaleY: int16(height / 2 * 4),
		ScaleZ: int16(0x1FF * 4),
		TransX: int16(width / 2 * 4),
		TransY: int16(height / 2 * 4),
		TransZ: int16(0x1FF * 4),
	}
}
