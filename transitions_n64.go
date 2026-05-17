//go:build n64

package gosprite64

import (
	"image"
	"image/color"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

func drawTransitionOverlay(c color.RGBA) {
	video := currentVideo()
	if video == nil || video.Framebuffer == nil {
		return
	}

	fbBounds := rendergeom.FramebufferBounds()

	rdp.RDP.SetOtherModes(
		rdp.ForceBlend,
		rdp.CycleTypeOne,
		rdp.RGBDitherNone,
		rdp.AlphaDitherNone,
		rdp.ZmodeOpaque,
		rdp.CvgDestClamp,
		rdp.BlendMode{
			P1: rdp.BlenderPMColorCombiner,
			A1: rdp.BlenderAColorCombinerAlpha,
			M1: rdp.BlenderPMFramebuffer,
			B1: rdp.BlenderBOneMinusAlphaA,
		},
	)
	rdp.RDP.SetCombineMode(rdp.CombineMode{
		Two: rdp.CombinePass{
			RGB:   rdp.CombineParams{0, 0, 0, rdp.CombinePrimitive},
			Alpha: rdp.CombineParams{0, 0, 0, rdp.CombinePrimitive},
		},
	})
	rdp.RDP.SetPrimitiveColor(c)
	rdp.RDP.Push(rdp.SyncPipe)
	rdp.RDP.FillRectangle(image.Rect(0, 0, fbBounds.Dx(), fbBounds.Dy()))
	rdp.RDP.Push(rdp.SyncPipe)

	rdp.RDP.SetOtherModes(
		0,
		rdp.CycleTypeOne,
		rdp.RGBDitherNone,
		rdp.AlphaDitherNone,
		rdp.ZmodeOpaque,
		rdp.CvgDestClamp,
		rdp.BlendMode{},
	)
}
