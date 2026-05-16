//go:build n64

package sprite

import (
	"image"
	"image/color"
	"math"

	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/rcp/rdp"
	"github.com/clktmr/n64/rcp/texture"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

const (
	blendMasked = 1
	blendAlpha  = 2
)

var blendSrcSprites = rdp.BlendMode{
	P1: rdp.BlenderPMColorCombiner,
	A1: rdp.BlenderAColorCombinerAlpha,
	M1: rdp.BlenderPMColorCombiner,
	B1: rdp.BlenderBZero,
}

var blendOverSprites = rdp.BlendMode{
	P1: rdp.BlenderPMColorCombiner,
	A1: rdp.BlenderAColorCombinerAlpha,
	M1: rdp.BlenderPMFramebuffer,
	B1: rdp.BlenderBOneMinusAlphaA,
}

// RenderSprite draws a sprite frame via the RDP with flip, scale, blend, and rotation.
// For non-rotated sprites, the fast TextureRectangle path is used. For rotated
// sprites, a software rotation fallback composites a transformed NRGBA image
// onto the framebuffer (higher cost, but correct for phase 1).
func RenderSprite(fb *texture.Texture, src image.Image, x, y int,
	flipH, flipV bool, scaleX, scaleY float32, blendMode uint8, alpha float32,
	rotation, originX, originY float32) {

	tex, ok := src.(*texture.Texture)
	if !ok || tex == nil || fb == nil {
		return
	}

	if rotation != 0 {
		renderRotatedSprite(fb, tex, x, y, flipH, flipV, scaleX, scaleY,
			blendMode, alpha, rotation, originX, originY)
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

	setupBlendMode(tex, blendMode, alpha)

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

// renderRotatedSprite creates a software-rotated NRGBA image and composites it
// onto the framebuffer. The RDP cannot draw rotated rectangles directly, so this
// path samples the source texture through an inverse rotation transform.
func renderRotatedSprite(fb *texture.Texture, tex *texture.Texture,
	drawX, drawY int, flipH, flipV bool, scaleX, scaleY float32,
	_ uint8, _ float32, rotation, originX, originY float32) {

	srcBounds := tex.Bounds()
	srcW := float64(srcBounds.Dx())
	srcH := float64(srcBounds.Dy())
	sx := float64(scaleX)
	sy := float64(scaleY)
	ox := float64(originX)
	oy := float64(originY)

	cos := math.Cos(float64(rotation))
	sin := math.Sin(float64(rotation))

	corners := [4][2]float64{
		{0, 0}, {srcW, 0}, {srcW, srcH}, {0, srcH},
	}

	var minX, minY, maxX, maxY float64
	for i, c := range corners {
		px := (c[0] - ox) * sx
		py := (c[1] - oy) * sy
		rx := px*cos - py*sin
		ry := px*sin + py*cos
		if i == 0 {
			minX, maxX = rx, rx
			minY, maxY = ry, ry
		} else {
			minX = min(minX, rx)
			maxX = max(maxX, rx)
			minY = min(minY, ry)
			maxY = max(maxY, ry)
		}
	}

	dstW := int(math.Ceil(maxX-minX)) + 1
	dstH := int(math.Ceil(maxY-minY)) + 1
	if dstW <= 0 || dstH <= 0 {
		return
	}

	rotated := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	invCos := cos
	invSin := -sin

	for dy := 0; dy < dstH; dy++ {
		for dx := 0; dx < dstW; dx++ {
			rx := float64(dx) + minX
			ry := float64(dy) + minY

			px := rx*invCos - ry*invSin
			py := rx*invSin + ry*invCos

			srcX := px/sx + ox
			srcY := py/sy + oy

			if flipH {
				srcX = srcW - srcX
			}
			if flipV {
				srcY = srcH - srcY
			}

			ix := int(srcX)
			iy := int(srcY)
			if ix < 0 || iy < 0 || ix >= int(srcW) || iy >= int(srcH) {
				continue
			}

			c := tex.At(srcBounds.Min.X+ix, srcBounds.Min.Y+iy)
			rotated.Set(dx, dy, c)
		}
	}

	logX := drawX + int(minX)
	logY := drawY + int(minY)
	logicalDst := image.Rect(logX, logY, logX+dstW, logY+dstH)
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
		clipped.Min.X-logicalDst.Min.X,
		clipped.Min.Y-logicalDst.Min.Y,
	)
	n64draw.Over.Draw(
		fb,
		image.Rect(framebufferRect.Min.X, framebufferRect.Min.Y,
			framebufferRect.Max.X+1, framebufferRect.Max.Y+1),
		rotated,
		srcPt,
	)
}

func setupBlendMode(tex *texture.Texture, blendMode uint8, alpha float32) {
	alphaSource := rdp.CombineTex0
	if !tex.HasAlpha() {
		alphaSource = rdp.CombineDAlphaOne
	}

	switch blendMode {
	case blendMasked:
		rdp.RDP.SetOtherModes(
			rdp.AlphaCompare|rdp.ForceBlend|rdp.ImageRead|rdp.BiLerp0,
			rdp.CycleTypeOne, rdp.RGBDitherNone, rdp.AlphaDitherNone,
			rdp.ZmodeOpaque, rdp.CvgDestClamp, blendOverSprites,
		)
		rdp.RDP.SetBlendColor(color.NRGBA{A: 1})
		rdp.RDP.SetCombineMode(rdp.CombineMode{
			Two: rdp.CombinePass{
				RGB:   rdp.CombineParams{0, 0, 0, rdp.CombineTex0},
				Alpha: rdp.CombineParams{0, 0, 0, alphaSource},
			},
		})

	case blendAlpha:
		a := uint8(clampf(alpha, 0, 1) * 255)
		rdp.RDP.SetEnvironmentColor(color.NRGBA{R: 255, G: 255, B: 255, A: a})
		rdp.RDP.SetOtherModes(
			rdp.ForceBlend|rdp.ImageRead|rdp.BiLerp0,
			rdp.CycleTypeOne, rdp.RGBDitherNone, rdp.AlphaDitherNone,
			rdp.ZmodeOpaque, rdp.CvgDestClamp, blendOverSprites,
		)
		rdp.RDP.SetCombineMode(rdp.CombineMode{
			Two: rdp.CombinePass{
				RGB: rdp.CombineParams{0, 0, 0, rdp.CombineTex0},
				Alpha: rdp.CombineParams{
					A: alphaSource,
					B: rdp.CombineAAlphaZero,
					C: rdp.CombineEnvironment,
					D: rdp.CombineDAlphaZero,
				},
			},
		})

	default:
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
	}
}

func clampf(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
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
