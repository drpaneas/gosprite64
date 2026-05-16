//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles_overlay.png -out assets/tiles_overlay.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2danim -in assets-src/idle.json -out assets/idle.anim && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -sheet assets/tiles_overlay.sheet -map assets/level.map -anim assets/idle.anim -out assets/level.bundle"

package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	scene  *gosprite64.Scene
	camera *gosprite64.Camera
	tick   int
}

func (g *Game) Init() {
	bundle, err := gosprite64.OpenBundle("assets/level.bundle")
	if err != nil {
		panic(err)
	}

	scene, err := gosprite64.LoadScene(bundle)
	if err != nil {
		panic(err)
	}

	g.scene = scene
	g.camera = &gosprite64.Camera{Width: 64, Height: 64}
}

func (g *Game) Update() {
	if g.camera == nil || g.scene == nil {
		return
	}

	g.tick++
	scene := g.scene
	m := scene.Map()
	if m == nil {
		return
	}

	g.camera.X = pingPong(g.tick/2, max(0, m.PixelWidth()-g.camera.Width))
	g.camera.Y = pingPong(g.tick/3, max(0, m.PixelHeight()-g.camera.Height))
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	scene := g.scene
	camera := g.camera
	if scene == nil || camera == nil {
		return
	}

	scene.Draw(camera)

	baseLayer, baseSheet, _ := scene.LayerAssets(0)
	overlayLayer, _, _ := scene.LayerAssets(1)
	baseSheetInfo, _ := scene.LayerSheetInfo(0)
	overlaySheetInfo, _ := scene.LayerSheetInfo(1)
	stats := scene.Stats()
	gosprite64.DrawText(fmt.Sprintf("vis:%d", stats.VisibleTiles), 2, 2, gosprite64.White)
	gosprite64.DrawText(fmt.Sprintf("up:%d", stats.UploadCount), 2, 12, gosprite64.White)
	gosprite64.DrawText(fmt.Sprintf("ram:%d/%d", stats.SheetRAMBytes, stats.MapRAMBytes), 2, 22, gosprite64.White)
	gosprite64.DrawText(fmt.Sprintf("sh:%d/%d", baseLayer.SheetID, overlayLayer.SheetID), 2, 32, gosprite64.White)
	gosprite64.DrawText(fmt.Sprintf("nz:%d/%d", baseLayer.NonZeroTiles, overlayLayer.NonZeroTiles), 2, 42, gosprite64.White)
	gosprite64.DrawText(fmt.Sprintf("tc:%d/%d", baseSheetInfo.TileCount, overlaySheetInfo.TileCount), 2, 52, gosprite64.White)

	anim := scene.AnimationByName("idle")
	if clip, ok := anim.Clip("idle"); ok && len(clip.Frames) > 0 {
		frameIdx, frameTile := animationCursor(g.tick, clip)
		// Animation frames are zero-based clip entries; sheet tiles are 1-based.
		animTile := baseSheet.Tile(frameTile + 1)
		gosprite64.DrawWorldImage(animTile, 56, 24, camera)
		gosprite64.DrawText(fmt.Sprintf("an:%s %d/%d f:%d@%d", clip.Name, frameIdx+1, len(clip.Frames), frameTile, clip.FPS), 2, 62, gosprite64.White)
	}
}

func main() {
	gosprite64.Run(&Game{})
}

func pingPong(step, limit int) int {
	if limit <= 0 {
		return 0
	}

	period := limit * 2
	pos := step % period
	if pos > limit {
		pos = period - pos
	}
	return pos
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func animationCursor(tick int, clip gosprite64.AnimationClip) (int, uint16) {
	if len(clip.Frames) == 0 {
		return 0, 0
	}
	stepEvery := 1
	if clip.FPS > 0 {
		stepEvery = max(1, 60/int(clip.FPS))
	}
	frameIdx := (tick / stepEvery) % len(clip.Frames)
	return frameIdx, clip.Frames[frameIdx]
}
