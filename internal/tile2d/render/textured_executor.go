package render

import (
	"github.com/clktmr/n64/rcp/texture"
)

type TexturedSetupState struct {
	Source  *texture.Texture
	DrawIdx uint8
	Ready   bool
}

type TexturedExecutor struct {
	Framebuffer *texture.Texture
	State       *TexturedSetupState
}

func (e TexturedExecutor) Ready() bool {
	return e.State != nil && e.State.Ready && e.State.Source != nil
}

func (e TexturedExecutor) SourceTexture(tile PreparedTile) (*texture.Texture, bool) {
	src, ok := tile.Source.(*texture.Texture)
	if !ok || src == nil {
		return nil, false
	}
	if src.Palette() != nil {
		return nil, false
	}
	return src, true
}
