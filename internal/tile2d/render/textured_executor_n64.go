//go:build n64

package render

import (
	"image"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

var blendSrcTiles = rdp.BlendMode{
	P1: rdp.BlenderPMColorCombiner,
	A1: rdp.BlenderAColorCombinerAlpha,
	M1: rdp.BlenderPMColorCombiner,
	B1: rdp.BlenderBZero,
}

func (e TexturedExecutor) EnsurePrepared(tile PreparedTile) bool {
	if e.Framebuffer == nil || e.State == nil {
		return false
	}

	src, ok := e.SourceTexture(tile)
	if !ok {
		return false
	}
	max := rdp.MaxTileSize(src.Format())
	if src.Bounds().Dx() > max.Dx() || src.Bounds().Dy() > max.Dy() {
		return false
	}

	if e.State.Ready && e.State.Source == src {
		return true
	}

	rdp.RDP.SetColorImage(e.Framebuffer)
	rdp.RDP.SetScissor(image.Rectangle{Max: e.Framebuffer.Bounds().Size()}, rdp.InterlaceNone)
	alphaSource := rdp.CombineTex0
	if !src.HasAlpha() {
		alphaSource = rdp.CombineDAlphaOne
	}
	rdp.RDP.SetOtherModes(
		rdp.ForceBlend|rdp.BiLerp0,
		rdp.CycleTypeOne, rdp.RGBDitherNone, rdp.AlphaDitherNone, rdp.ZmodeOpaque, rdp.CvgDestClamp, blendSrcTiles,
	)
	rdp.RDP.SetCombineMode(rdp.CombineMode{
		Two: rdp.CombinePass{
			RGB:   rdp.CombineParams{0, 0, 0, rdp.CombineTex0},
			Alpha: rdp.CombineParams{0, 0, 0, alphaSource},
		},
	})
	rdp.RDP.SetTextureImage(src)
	loadIdx, drawIdx := rdp.RDP.SetTile(rdp.TileDescriptor{
		Format: src.Format(),
		Addr:   0x0,
		Line:   uint16(src.Format().TMEMWords(src.Bounds().Dx())),
	})
	rdp.RDP.LoadTile(loadIdx, src.Bounds())

	e.State.Source = src
	e.State.DrawIdx = drawIdx
	e.State.Ready = true
	return true
}

func (e TexturedExecutor) BlitRun(x, y, tileWidth int, count int) {
	for i := 0; i < count; i++ {
		e.BlitTile(x+(i*tileWidth), y)
	}
}

func (e TexturedExecutor) BlitTile(x, y int) {
	if !e.Ready() {
		return
	}

	srcBounds := e.State.Source.Bounds()
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
	rdp.RDP.TextureRectangle(
		image.Rect(framebufferRect.Min.X, framebufferRect.Min.Y, framebufferRect.Max.X+1, framebufferRect.Max.Y+1),
		srcPt,
		image.Point{1, 1},
		e.State.DrawIdx,
	)
}
