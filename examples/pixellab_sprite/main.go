package main

const (
	logicalWidth  = 288
	logicalHeight = 216
)

func centeredHeroPosition(frameWidth, frameHeight int) (float32, float32) {
	return float32((logicalWidth - frameWidth) / 2), float32((logicalHeight - frameHeight) / 2)
}
