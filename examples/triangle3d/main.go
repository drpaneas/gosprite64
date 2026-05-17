package main

import (
	"image/color"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/gfx"
	"github.com/drpaneas/gosprite64/math3d"
	"github.com/drpaneas/gosprite64/rspq/ucode"
)

const (
	screenW = 320
	screenH = 240
)

type Game struct {
	angle float32
}

func (g *Game) Init() {}

func (g *Game) Update() {
	g.angle += 1.0
	if g.angle >= 360 {
		g.angle -= 360
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(color.RGBA{R: 20, G: 20, B: 40, A: 255})

	dl := gfx.NewDisplayList(256)

	dl.DPPipeSync()
	dl.DPSetScissor(0, 0, 0, screenW, screenH)

	proj, perspNorm := math3d.Perspective(45, float32(screenW)/float32(screenH), 10, 1000, 1.0)
	view := math3d.LookAt(0, 0, 300, 0, 0, 0, 0, 1, 0)
	model := math3d.Rotate(g.angle, 0, 1, 0)
	mv := view.Mul(model)

	projMtx := proj.ToN64Mtx()
	mvMtx := mv.ToN64Mtx()

	dl.SPPerspNormalize(perspNorm)

	// These SPMatrix commands reference physical RDRAM addresses.
	// On real hardware, projMtx and mvMtx would need to be placed in
	// cache-flushed RDRAM and their physical addresses passed here.
	// For now this demonstrates the display list construction.
	_ = projMtx
	_ = mvMtx

	dl.SPSetGeometryMode(gfx.GZBuffer | gfx.GShade | gfx.GCullBack)

	dl.DPFullSync()
	dl.SPEndDisplayList()

	gfx.ExecuteViaRSP(dl, ucode.RSPBoot, ucode.F3DEX2Text, ucode.F3DEX2Data)
}

func main() {
	gosprite64.Run(&Game{})
}
