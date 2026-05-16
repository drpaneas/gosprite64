package render

import (
	"image"
	"testing"

	"github.com/clktmr/n64/rcp/texture"
)

func TestTexturedExecutorSourceTextureAcceptsTextureSources(t *testing.T) {
	src := texture.NewTextureFromImage(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	exec := TexturedExecutor{}

	got, ok := exec.SourceTexture(PreparedTile{Source: src})
	if !ok {
		t.Fatal("expected texture source to be accepted")
	}
	if got != src {
		t.Fatal("expected same texture pointer back")
	}
}

func TestTexturedExecutorSourceTextureRejectsGenericImages(t *testing.T) {
	exec := TexturedExecutor{}
	if _, ok := exec.SourceTexture(PreparedTile{Source: image.NewRGBA(image.Rect(0, 0, 8, 8))}); ok {
		t.Fatal("expected generic image source to be rejected")
	}
}
