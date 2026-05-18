package main

import (
	"encoding/binary"
	"fmt"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/save"
)

type Game struct {
	storage save.Storage
	counter uint32
	status  string
}

func (g *Game) Init() {
	g.storage = save.NewSRAM()
	g.status = "PRESS A TO INCREMENT AND SAVE"
	g.loadCounter()
}

func (g *Game) loadCounter() {
	data, err := save.ReadAll(g.storage)
	if err != nil {
		g.status = "READ ERROR: " + err.Error()
		return
	}

	if len(data) < 8 {
		return
	}

	value := binary.BigEndian.Uint32(data[0:4])
	stored := binary.BigEndian.Uint32(data[4:8])
	expected := save.Checksum(data[0:4])

	if stored == expected {
		g.counter = value
		g.status = "LOADED FROM SRAM"
	} else {
		g.counter = 0
		g.status = "NO VALID SAVE FOUND"
	}
}

func (g *Game) saveCounter() {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], g.counter)
	binary.BigEndian.PutUint32(buf[4:8], save.Checksum(buf[0:4]))

	if err := save.WriteAll(g.storage, buf); err != nil {
		g.status = "WRITE ERROR: " + err.Error()
		return
	}
	g.status = "SAVED!"
}

func (g *Game) Update() {
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		g.counter++
		g.saveCounter()
	}
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonB) {
		g.counter = 0
		g.saveCounter()
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawText("SAVE DATA DEMO", 88, 10, gosprite64.White)
	gosprite64.DrawText("SRAM PERSISTENCE EXAMPLE", 48, 26, gosprite64.LightGray)

	gosprite64.DrawText(
		fmt.Sprintf("COUNTER: %d", g.counter),
		88, 70, gosprite64.Yellow,
	)

	gosprite64.DrawText("A: INCREMENT + SAVE", 60, 120, gosprite64.Green)
	gosprite64.DrawText("B: RESET TO ZERO", 60, 136, gosprite64.Red)

	gosprite64.DrawText(g.status, 20, 180, gosprite64.LightGray)

	gosprite64.DrawText(
		fmt.Sprintf("STORAGE: %s (%d BYTES)",
			g.storage.Type(), g.storage.Size()),
		20, 200, gosprite64.DarkGray,
	)
}

func main() {
	gosprite64.Run(&Game{})
}
