//go:build n64

package gosprite64

import (
	"image"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

func applyScissor(r DrawRegion) {
	fb := rendergeom.FramebufferBounds()
	sx := r.X * fb.Dx() / 288
	sy := r.Y * fb.Dy() / 216
	sw := r.W * fb.Dx() / 288
	sh := r.H * fb.Dy() / 216
	rdp.RDP.SetScissor(image.Rect(sx, sy, sx+sw, sy+sh), rdp.InterlaceNone)
}

// clearScissor resets the RDP scissor to allow drawing across the full
// framebuffer. The 2x multiplier matches the scissor set in gameloop.go's
// Run(), which uses fb*2 to account for the RDP's 10.2 fixed-point
// coordinate format where scissor bounds are in subpixels.
func clearScissor() {
	fb := rendergeom.FramebufferBounds()
	rdp.RDP.SetScissor(image.Rect(0, 0, fb.Dx()*2, fb.Dy()*2), rdp.InterlaceNone)
}
