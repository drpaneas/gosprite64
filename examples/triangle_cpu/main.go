//go:build n64

package main

import (
	"image"
	"image/color"
	"time"

	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/drivers/display"
	"github.com/clktmr/n64/machine"
	"github.com/clktmr/n64/rcp/rdp"
	"github.com/clktmr/n64/rcp/texture"
	"github.com/clktmr/n64/rcp/video"
	"github.com/drpaneas/gosprite64/gfx"
	"github.com/drpaneas/gosprite64/internal/rdpcpu"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

var blendSrcTriangles = rdp.BlendMode{
	P1: rdp.BlenderPMColorCombiner,
	A1: rdp.BlenderAColorCombinerAlpha,
	M1: rdp.BlenderPMColorCombiner,
	B1: rdp.BlenderBZero,
}

func main() {
	video.Setup(false)
	video.SetScale(squarePixelPresentationRect())

	fbBounds := rendergeom.FramebufferBounds()
	disp := display.NewDisplay(fbBounds.Size(), video.BPP16)
	tex := makeCheckerTexture()
	var angle float32

	for {
		fb := disp.Swap()
		drawFrame(fb, tex, angle)
		angle += 1.5
		if angle >= 360 {
			angle -= 360
		}
		time.Sleep(time.Second / 60)
	}
}

func drawFrame(fb *texture.Texture, tex *texture.Texture, angle float32) {
	n64draw.Src.Draw(fb, fb.Bounds(), image.NewUniform(color.RGBA{R: 20, G: 20, B: 40, A: 255}), image.Point{})

	rdp.RDP.SetColorImage(fb)
	rdp.RDP.SetScissor(image.Rectangle{Max: fb.Bounds().Size()}, rdp.InterlaceNone)
	rdp.RDP.SetOtherModes(
		rdp.ForceBlend|rdp.BiLerp0,
		rdp.CycleTypeOne,
		rdp.RGBDitherNone,
		rdp.AlphaDitherNone,
		rdp.ZmodeOpaque,
		rdp.CvgDestClamp,
		blendSrcTriangles,
	)
	rdp.RDP.SetCombineMode(rdp.CombineMode{
		Two: rdp.CombinePass{
			RGB:   rdp.CombineParams{0, 0, 0, rdp.CombineTex0},
			Alpha: rdp.CombineParams{0, 0, 0, rdp.CombineTex0},
		},
	})
	rdp.RDP.SetTextureImage(tex)
	loadIdx, drawIdx := rdp.RDP.SetTile(rdp.TileDescriptor{
		Format: tex.Format(),
		Addr:   0x0,
		Line:   uint16(tex.Format().TMEMWords(tex.Bounds().Dx())),
	})
	rdp.RDP.LoadTile(loadIdx, tex.Bounds())

	verts := buildProjectedTriangle(angle)
	packet := rdpcpu.BuildTexturedTriangle(drawIdx, 0, verts[0], verts[1], verts[2])
	gfx.PushRaw(packet...)
	rdp.RDP.Flush()
}

func makeCheckerTexture() *texture.Texture {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	c0 := color.RGBA{R: 255, G: 96, B: 96, A: 255}
	c1 := color.RGBA{R: 255, G: 240, B: 96, A: 255}
	c2 := color.RGBA{R: 96, G: 192, B: 255, A: 255}
	c3 := color.RGBA{R: 96, G: 255, B: 160, A: 255}
	palette := [4]color.RGBA{c0, c1, c2, c3}
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			idx := ((x / 4) + (y / 4)) & 0x3
			img.SetRGBA(x, y, palette[idx])
		}
	}
	return texture.NewTextureFromImage(img)
}

func squarePixelPresentationRect() image.Rectangle {
	outputSize := rendergeom.FramebufferBounds().Size().Mul(2)
	switch machine.VideoType {
	case machine.VideoPAL:
		return rendergeom.CenteredRect(image.Rect(128, 45, 128+640, 45+576), outputSize)
	case machine.VideoMPAL, machine.VideoNTSC:
		return rendergeom.CenteredRect(image.Rect(108, 35, 108+640, 35+480), outputSize)
	default:
		return rendergeom.CenteredRect(video.Scale(), outputSize)
	}
}
