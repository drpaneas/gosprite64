package gosprite64

type ParallaxLayer struct {
	SpeedX float32
	SpeedY float32
}

func (p ParallaxLayer) Offset(cameraX, cameraY int) (int, int) {
	return int(float32(cameraX) * p.SpeedX), int(float32(cameraY) * p.SpeedY)
}

type ParallaxConfig struct {
	Layers []ParallaxLayer
}

func NewParallaxConfig(speeds ...ParallaxLayer) ParallaxConfig {
	return ParallaxConfig{Layers: speeds}
}

func (pc ParallaxConfig) LayerOffset(layer, cameraX, cameraY int) (int, int) {
	if layer < 0 || layer >= len(pc.Layers) {
		return cameraX, cameraY
	}
	return pc.Layers[layer].Offset(cameraX, cameraY)
}
