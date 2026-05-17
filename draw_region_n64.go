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

func clearScissor() {
	fb := rendergeom.FramebufferBounds()
	rdp.RDP.SetScissor(image.Rect(0, 0, fb.Dx()*2, fb.Dy()*2), rdp.InterlaceNone)
}
