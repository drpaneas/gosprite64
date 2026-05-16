//go:build !n64

package sprite

import (
	"image"
	"image/draw"

	"github.com/clktmr/n64/rcp/texture"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

// RenderSprite draws a sprite using software rendering on the host.
// Flip, scale, and rotation options are ignored; the image is drawn at the
// given position. Rotation is not implemented on the host fallback path.
func RenderSprite(fb *texture.Texture, src image.Image, x, y int,
	_ bool, _ bool, _ float32, _ float32, blendMode uint8, _ float32,
	_ float32, _ float32, _ float32) {

	if fb == nil || src == nil {
		return
	}

	srcBounds := src.Bounds()
	logicalDst := image.Rect(x, y, x+srcBounds.Dx(), y+srcBounds.Dy())
	clipped := logicalDst.Intersect(rendergeom.LogicalBounds())
	if clipped.Empty() {
		return
	}

	framebufferRect, ok := rendergeom.MapRectInclusive(image.Rectangle{
		Min: clipped.Min,
		Max: clipped.Max.Sub(image.Pt(1, 1)),
	})
	if !ok {
		return
	}

	srcPt := image.Pt(
		srcBounds.Min.X+(clipped.Min.X-logicalDst.Min.X),
		srcBounds.Min.Y+(clipped.Min.Y-logicalDst.Min.Y),
	)

	op := draw.Src
	if blendMode >= 1 {
		op = draw.Over
	}
	op.Draw(
		fb,
		image.Rect(framebufferRect.Min.X, framebufferRect.Min.Y,
			framebufferRect.Max.X+1, framebufferRect.Max.Y+1),
		src,
		srcPt,
	)
}
