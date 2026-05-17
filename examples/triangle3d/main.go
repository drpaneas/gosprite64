//go:build n64

package main

import (
	"image"
	"image/color"
	"math"
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

const (
	screenW = 320
	screenH = 240
)

var blendSrc = rdp.BlendMode{
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
	var angle float32

	for {
		fb := disp.Swap()
		drawFrame(fb, angle)
		angle += 1.0
		if angle >= 360 {
			angle -= 360
		}
		time.Sleep(time.Second / 60)
	}
}

var bgColor = image.NewUniform(color.RGBA{R: 20, G: 20, B: 40, A: 255})

var combineMode = rdp.CombineMode{
	Two: rdp.CombinePass{
		RGB:   rdp.CombineParams{0, 0, 0, rdp.CombinePrimitive},
		Alpha: rdp.CombineParams{0, 0, 0, rdp.CombineDAlphaOne},
	},
}

const (
	halfW   = float32(screenW) / 2
	halfH   = float32(screenH) / 2
	cameraZ = float32(300)
	focal   = float32(250)
)

type vert3 struct{ x, y, z float32 }

var worldVerts = [3]vert3{
	{0, 80, -50},
	{-80, -50, 50},
	{80, -50, 30},
}

func drawFrame(fb *texture.Texture, angle float32) {
	n64draw.Src.Draw(fb, fb.Bounds(), bgColor, image.Point{})

	rdp.RDP.SetColorImage(fb)
	rdp.RDP.SetScissor(image.Rectangle{Max: fb.Bounds().Size()}, rdp.InterlaceNone)

	rdp.RDP.SetOtherModes(
		rdp.ForceBlend,
		rdp.CycleTypeOne,
		rdp.RGBDitherNone,
		rdp.AlphaDitherNone,
		rdp.ZmodeOpaque,
		rdp.CvgDestClamp,
		blendSrc,
	)
	rdp.RDP.SetCombineMode(combineMode)

	rad := angle * math.Pi / 180.0
	sinA := float32(math.Sin(float64(rad)))
	cosA := float32(math.Cos(float64(rad)))

	var sv [3][2]float32
	allVisible := true
	for i, v := range worldVerts {
		rotX := v.x*cosA + v.z*sinA
		rotZ := -v.x*sinA + v.z*cosA
		viewZ := rotZ + cameraZ
		if viewZ < 1 {
			allVisible = false
			break
		}
		sv[i][0] = halfW + (rotX/viewZ)*focal
		sv[i][1] = halfH - (v.y/viewZ)*focal
	}

	if allVisible {
		rdp.RDP.SetPrimitiveColor(color.RGBA{R: 255, G: 50, B: 50, A: 255})
		rdp.RDP.Push(rdp.SyncPipe)
		cmds := rdpcpu.FillTriangle(sv[0], sv[1], sv[2])
		gfx.PushRaw(cmds...)
	}

	rdp.RDP.Flush()
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
