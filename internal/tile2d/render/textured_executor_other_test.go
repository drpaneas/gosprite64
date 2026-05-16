//go:build !n64

package render

import (
	"image"
	"testing"

	"github.com/clktmr/n64/rcp/texture"
)

func TestTexturedExecutorEnsurePreparedFallsBackOnHost(t *testing.T) {
	state := &TexturedSetupState{}
	exec := TexturedExecutor{
		State: state,
	}
	src := texture.NewTextureFromImage(image.NewRGBA(image.Rect(0, 0, 8, 8)))

	if ok := exec.EnsurePrepared(PreparedTile{Source: src}); ok {
		t.Fatal("expected host EnsurePrepared to report fallback")
	}
	if state.Ready {
		t.Fatal("host EnsurePrepared must not mark textured state ready")
	}
	if state.Source != nil {
		t.Fatal("host EnsurePrepared must not retain prepared source")
	}
}
