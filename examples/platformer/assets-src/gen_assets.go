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
	generateTiles()
	generateCharacter()
	generateLevel()
	generateAnims()
}

func generateTiles() {
	img := image.NewRGBA(image.Rect(0, 0, 16, 8))
	sky := color.RGBA{R: 92, G: 148, B: 252, A: 255}
	grass := color.RGBA{R: 56, G: 183, B: 100, A: 255}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.SetRGBA(x, y, sky)
		}
		for x := 8; x < 16; x++ {
			img.SetRGBA(x, y, grass)
		}
	}
	writeImage("tiles.png", img)
}

func generateCharacter() {
	img := image.NewRGBA(image.Rect(0, 0, 64, 16))
	body := color.RGBA{R: 220, G: 80, B: 60, A: 255}
	head := color.RGBA{R: 255, G: 200, B: 150, A: 255}
	feet := color.RGBA{R: 80, G: 60, B: 40, A: 255}

	for frame := 0; frame < 4; frame++ {
		ox := frame * 16
		for y := 0; y < 5; y++ {
			for x := 5; x < 11; x++ {
				img.SetRGBA(ox+x, y, head)
			}
		}
		for y := 5; y < 12; y++ {
			for x := 4; x < 12; x++ {
				img.SetRGBA(ox+x, y, body)
			}
		}
		legOffset := frame % 2
		for y := 12; y < 16; y++ {
			img.SetRGBA(ox+5+legOffset, y, feet)
			img.SetRGBA(ox+6+legOffset, y, feet)
			img.SetRGBA(ox+9-legOffset, y, feet)
			img.SetRGBA(ox+10-legOffset, y, feet)
		}
	}
	writeImage("character.png", img)
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
	const w, h = 48, 32
	cells := make([]int, w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cells[y*w+x] = 1
		}
	}

	// Ground floor (bottom 3 rows)
	for y := h - 3; y < h; y++ {
		for x := 0; x < w; x++ {
			cells[y*w+x] = 2
		}
	}

	// Platform 1: low left
	for x := 6; x < 12; x++ {
		cells[(h-7)*w+x] = 2
	}

	// Platform 2: mid center
	for x := 16; x < 24; x++ {
		cells[(h-11)*w+x] = 2
	}

	// Platform 3: high right
	for x := 30; x < 38; x++ {
		cells[(h-15)*w+x] = 2
	}

	// Platform 4: top left
	for x := 8; x < 16; x++ {
		cells[(h-19)*w+x] = 2
	}

	// Platform 5: top center
	for x := 22; x < 28; x++ {
		cells[(h-23)*w+x] = 2
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
