//go:build n64

package sprite

import (
	"image"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/clktmr/n64/rcp/texture"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

var blendSrcSprites = rdp.BlendMode{
	P1: rdp.BlenderPMColorCombiner,
	A1: rdp.BlenderAColorCombinerAlpha,
	M1: rdp.BlenderPMColorCombiner,
	B1: rdp.BlenderBZero,
}

// RenderSprite draws a sprite frame via the RDP with flip and scale support.
// Flip is implemented using the RDP's native MirrorS/MirrorT tile flags,
// which requires power-of-2 texture dimensions. For non-power-of-2 textures,
// flip flags are silently ignored.
func RenderSprite(fb *texture.Texture, src image.Image, x, y int,
	flipH, flipV bool, scaleX, scaleY float32) {

	tex, ok := src.(*texture.Texture)
	if !ok || tex == nil || fb == nil {
		return
	}

	srcBounds := tex.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	maxTile := rdp.MaxTileSize(tex.Format())
	if srcW > maxTile.Dx() || srcH > maxTile.Dy() {
		return
	}

	if flipH && !isPowerOf2(srcW) {
		flipH = false
	}
	if flipV && !isPowerOf2(srcH) {
		flipV = false
	}

	rdp.RDP.SetColorImage(fb)
	rdp.RDP.SetScissor(image.Rectangle{Max: fb.Bounds().Size()}, rdp.InterlaceNone)

	alphaSource := rdp.CombineTex0
	if !tex.HasAlpha() {
		alphaSource = rdp.CombineDAlphaOne
	}

	rdp.RDP.SetOtherModes(
		rdp.ForceBlend|rdp.BiLerp0,
		rdp.CycleTypeOne, rdp.RGBDitherNone, rdp.AlphaDitherNone,
		rdp.ZmodeOpaque, rdp.CvgDestClamp, blendSrcSprites,
	)
	rdp.RDP.SetCombineMode(rdp.CombineMode{
		Two: rdp.CombinePass{
			RGB:   rdp.CombineParams{0, 0, 0, rdp.CombineTex0},
			Alpha: rdp.CombineParams{0, 0, 0, alphaSource},
		},
	})

	tileDesc := rdp.TileDescriptor{
		Format: tex.Format(),
		Addr:   0x0,
		Line:   uint16(tex.Format().TMEMWords(srcW)),
	}
	if flipH {
		tileDesc.Flags |= rdp.MirrorS
		tileDesc.MaskS = log2u(srcW)
	}
	if flipV {
		tileDesc.Flags |= rdp.MirrorT
		tileDesc.MaskT = log2u(srcH)
	}

	rdp.RDP.SetTextureImage(tex)
	loadIdx, drawIdx := rdp.RDP.SetTile(tileDesc)
	rdp.RDP.LoadTile(loadIdx, srcBounds)

	rdpScaleX := max(1, int(scaleX+0.5))
	rdpScaleY := max(1, int(scaleY+0.5))

	destW := srcW * rdpScaleX
	destH := srcH * rdpScaleY

	logicalDst := image.Rect(x, y, x+destW, y+destH)
	clipped := logicalDst.Intersect(rendergeom.LogicalBounds())
	if clipped.Empty() {
		return
	}

	framebufferRect, ok2 := rendergeom.MapRectInclusive(image.Rectangle{
		Min: clipped.Min,
		Max: clipped.Max.Sub(image.Pt(1, 1)),
	})
	if !ok2 {
		return
	}

	clipOffsetX := clipped.Min.X - logicalDst.Min.X
	clipOffsetY := clipped.Min.Y - logicalDst.Min.Y

	srcPtX := srcBounds.Min.X
	srcPtY := srcBounds.Min.Y

	if flipH {
		srcPtX = srcBounds.Min.X + srcW
	}
	if flipV {
		srcPtY = srcBounds.Min.Y + srcH
	}

	srcPtX += clipOffsetX / rdpScaleX
	srcPtY += clipOffsetY / rdpScaleY

	rdp.RDP.TextureRectangle(
		image.Rect(
			framebufferRect.Min.X, framebufferRect.Min.Y,
			framebufferRect.Max.X+1, framebufferRect.Max.Y+1,
		),
		image.Pt(srcPtX, srcPtY),
		image.Point{X: rdpScaleX, Y: rdpScaleY},
		drawIdx,
	)
}

func isPowerOf2(n int) bool {
	return n > 0 && n&(n-1) == 0
}

func log2u(n int) uint8 {
	var r uint8
	for n > 1 {
		n >>= 1
		r++
	}
	return r
}
