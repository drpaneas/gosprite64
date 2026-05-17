package math3d

// Viewport matches the N64 Vp_t structure used by gSPViewport.
// Scale and Trans have 2 bits of fraction (multiply by 4).
type Viewport struct {
	ScaleX, ScaleY, ScaleZ, ScalePad int16
	TransX, TransY, TransZ, TransPad int16
}

// NewViewport creates a viewport for the given screen dimensions.
// The viewport maps NDC [-1,1] to screen coordinates.
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
