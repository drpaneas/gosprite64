//go:build ignore

package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	generateCharacter()
	generateTiles()
	generateLevel()
	generateAnims()
}

func generateCharacter() {
	img := image.NewRGBA(image.Rect(0, 0, 64, 16))
	colors := []color.RGBA{
		{R: 255, G: 0, B: 77, A: 255},
		{R: 0, G: 228, B: 54, A: 255},
		{R: 41, G: 173, B: 255, A: 255},
		{R: 255, G: 236, B: 39, A: 255},
	}
	for i, c := range colors {
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				img.SetRGBA(i*16+x, y, c)
			}
		}
	}
	writeImage("character.png", img)
}

func generateTiles() {
	img := image.NewRGBA(image.Rect(0, 0, 16, 8))
	grass := color.RGBA{R: 0, G: 180, B: 60, A: 255}
	dirt := color.RGBA{R: 160, G: 100, B: 50, A: 255}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.SetRGBA(x, y, grass)
		}
		for x := 8; x < 16; x++ {
			img.SetRGBA(x, y, dirt)
		}
	}
	writeImage("tiles.png", img)
}

type levelJSON struct {
	Width      int        `json:"width"`
	Height     int        `json:"height"`
	LayerCount int        `json:"layer_count"`
	CellBits   int        `json:"cell_bits"`
	ChunkW     int        `json:"chunk_width"`
	ChunkH     int        `json:"chunk_height"`
	Layers     []layerDef `json:"layers"`
}

type layerDef struct {
	SheetID int   `json:"sheet_id"`
	Cells   []int `json:"cells"`
}

func generateLevel() {
	const w, h = 48, 36
	cells := make([]int, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < 2 || x >= w-2 || y < 2 || y >= h-2 {
				cells[y*w+x] = 1
			} else if (x >= 10 && x < 15 && y >= 8 && y < 12) ||
				(x >= 25 && x < 32 && y >= 16 && y < 22) ||
				(x >= 8 && x < 14 && y >= 24 && y < 30) ||
				(x >= 34 && x < 40 && y >= 6 && y < 10) ||
				(x >= 20 && x < 28 && y >= 28 && y < 33) {
				cells[y*w+x] = 2
			}
		}
	}
	lev := levelJSON{
		Width: w, Height: h,
		LayerCount: 1, CellBits: 16,
		ChunkW: 8, ChunkH: 8,
		Layers: []layerDef{{SheetID: 1, Cells: cells}},
	}
	data, err := json.Marshal(lev)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("level.json", data, 0o644); err != nil {
		panic(err)
	}
}

type animJSON struct {
	Clips []clipDef `json:"clips"`
}

type clipDef struct {
	Name   string `json:"name"`
	FPS    int    `json:"fps"`
	Frames []int  `json:"frames"`
}

func generateAnims() {
	a := animJSON{
		Clips: []clipDef{
			{Name: "idle", FPS: 4, Frames: []int{0, 1}},
			{Name: "walk", FPS: 8, Frames: []int{0, 1, 2, 3}},
		},
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("anims.json", data, 0o644); err != nil {
		panic(err)
	}
}

func writeImage(name string, img image.Image) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
