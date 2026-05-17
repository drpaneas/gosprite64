//go:build n64

package sprite

import (
	"image"
	"image/color"
	"math"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/clktmr/n64/rcp/texture"
	"github.com/drpaneas/gosprite64/gfx"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
	"github.com/drpaneas/gosprite64/internal/rdpcpu"
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
// sprites, a CPU triangle setup path emits two textured RDP triangles.
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

	if scaleX < 1 {
		scaleX = 1
	}
	if scaleY < 1 {
		scaleY = 1
	}

	destW := int(float32(srcW) * scaleX)
	destH := int(float32(srcH) * scaleY)

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

	srcPtX += int(float32(clipOffsetX) / scaleX)
	srcPtY += int(float32(clipOffsetY) / scaleY)

	rdpScaleX := max(1, int(scaleX+0.5))
	rdpScaleY := max(1, int(scaleY+0.5))

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

// renderRotatedSprite draws a rotated sprite as two textured triangles.
func renderRotatedSprite(fb *texture.Texture, tex *texture.Texture,
	drawX, drawY int, flipH, flipV bool, scaleX, scaleY float32,
	blendMode uint8, alpha float32, rotation, originX, originY float32) {

	srcBounds := tex.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	maxTile := rdp.MaxTileSize(tex.Format())
	if srcW > maxTile.Dx() || srcH > maxTile.Dy() {
		return
	}

	rdp.RDP.SetColorImage(fb)
	rdp.RDP.SetScissor(image.Rectangle{Max: fb.Bounds().Size()}, rdp.InterlaceNone)
	setupBlendMode(tex, blendMode, alpha)

	tileDesc := rdp.TileDescriptor{
		Format: tex.Format(),
		Addr:   0x0,
		Line:   uint16(tex.Format().TMEMWords(srcW)),
	}
	rdp.RDP.SetTextureImage(tex)
	loadIdx, drawIdx := rdp.RDP.SetTile(tileDesc)
	rdp.RDP.LoadTile(loadIdx, srcBounds)

	frameOrigin := rendergeom.Origin()
	quad := rotatedQuad(float64(drawX+frameOrigin.X), float64(drawY+frameOrigin.Y),
		float64(srcW), float64(srcH), float64(scaleX), float64(scaleY),
		float64(rotation), float64(originX), float64(originY))
	st := textureCoords(float32(srcW), float32(srcH), flipH, flipV)

	packet1 := rdpcpu.BuildTexturedTriangle(drawIdx, 0,
		rdpcpu.TexVertex{X: float32(quad[0][0]), Y: float32(quad[0][1]), S: st[0][0], T: st[0][1], InvW: 1},
		rdpcpu.TexVertex{X: float32(quad[1][0]), Y: float32(quad[1][1]), S: st[1][0], T: st[1][1], InvW: 1},
		rdpcpu.TexVertex{X: float32(quad[2][0]), Y: float32(quad[2][1]), S: st[2][0], T: st[2][1], InvW: 1},
	)
	packet2 := rdpcpu.BuildTexturedTriangle(drawIdx, 0,
		rdpcpu.TexVertex{X: float32(quad[0][0]), Y: float32(quad[0][1]), S: st[0][0], T: st[0][1], InvW: 1},
		rdpcpu.TexVertex{X: float32(quad[2][0]), Y: float32(quad[2][1]), S: st[2][0], T: st[2][1], InvW: 1},
		rdpcpu.TexVertex{X: float32(quad[3][0]), Y: float32(quad[3][1]), S: st[3][0], T: st[3][1], InvW: 1},
	)
	gfx.PushRaw(packet1...)
	gfx.PushRaw(packet2...)
}

func rotatedQuad(drawX, drawY, srcW, srcH, scaleX, scaleY, rotation, originX, originY float64) [4][2]float64 {
	cos := math.Cos(rotation)
	sin := math.Sin(rotation)
	corners := [4][2]float64{
		{0, 0}, {srcW, 0}, {srcW, srcH}, {0, srcH},
	}
	var quad [4][2]float64
	for i, c := range corners {
		px := (c[0] - originX) * scaleX
		py := (c[1] - originY) * scaleY
		rx := px*cos - py*sin
		ry := px*sin + py*cos
		quad[i] = [2]float64{drawX + rx, drawY + ry}
	}
	return quad
}

func textureCoords(srcW, srcH float32, flipH, flipV bool) [4][2]float32 {
	left, right := float32(0), srcW
	top, bottom := float32(0), srcH
	if flipH {
		left, right = right, left
	}
	if flipV {
		top, bottom = bottom, top
	}
	return [4][2]float32{
		{left, top},
		{right, top},
		{right, bottom},
		{left, bottom},
	}
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
