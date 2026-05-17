package main

import (
	"math"

	"github.com/drpaneas/gosprite64/internal/rdpcpu"
)

const (
	screenW = 320
	screenH = 240
)

func buildProjectedTriangle(angle float32) [3]rdpcpu.TexVertex {
	src := [3]struct {
		x float32
		y float32
		z float32
		s   float32
		t   float32
	}{
		{x: 0, y: 80, z: 0, s: 0, t: 0},
		{x: -70, y: -40, z: 0, s: 16, t: 0},
		{x: 70, y: -40, z: 0, s: 0, t: 16},
	}

	rad := float32(angle * math.Pi / 180.0)
	sinA := float32(math.Sin(float64(rad)))
	cosA := float32(math.Cos(float64(rad)))
	cameraZ := float32(260)
	focal := float32(120)

	var out [3]rdpcpu.TexVertex
	for i, v := range src {
		rotX := v.x*cosA + v.z*sinA
		rotZ := -v.x*sinA + v.z*cosA
		viewZ := rotZ + cameraZ
		if viewZ <= 0 {
			continue
		}
		screenX := float32(screenW)/2 + (rotX/viewZ)*focal
		screenY := float32(screenH)/2 - (v.y/viewZ)*focal
		out[i] = rdpcpu.TexVertex{
			X:    screenX,
			Y:    screenY,
			S:    v.s,
			T:    v.t,
			InvW: 1 / viewZ,
		}
	}
	return out
}
