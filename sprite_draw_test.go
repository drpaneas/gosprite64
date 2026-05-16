package gosprite64

import "testing"

func TestDrawSpriteOptionsDefaults(t *testing.T) {
	var opts DrawSpriteOptions
	if opts.effectiveScaleX() != 1 {
		t.Fatalf("zero ScaleX should default to 1, got %f", opts.effectiveScaleX())
	}
	if opts.effectiveScaleY() != 1 {
		t.Fatalf("zero ScaleY should default to 1, got %f", opts.effectiveScaleY())
	}
	if opts.effectiveAlpha() != 1 {
		t.Fatalf("zero Alpha should default to 1, got %f", opts.effectiveAlpha())
	}
}

func TestDrawSpriteOptionsExplicitValues(t *testing.T) {
	opts := DrawSpriteOptions{ScaleX: 2, ScaleY: 0.5, Alpha: 0.7}
	if opts.effectiveScaleX() != 2 {
		t.Fatalf("ScaleX=2 should return 2, got %f", opts.effectiveScaleX())
	}
	if opts.effectiveScaleY() != 0.5 {
		t.Fatalf("ScaleY=0.5 should return 0.5, got %f", opts.effectiveScaleY())
	}
	if opts.effectiveAlpha() != 0.7 {
		t.Fatalf("Alpha=0.7 should return 0.7, got %f", opts.effectiveAlpha())
	}
}

func TestDrawSpriteNilSheetIsNoop(t *testing.T) {
	DrawSprite(nil, 0, 10, 20)
	DrawSpriteWithOptions(nil, 0, 10, 20, DrawSpriteOptions{})
	DrawWorldSprite(nil, 0, 10, 20, nil)
	DrawWorldSpriteWithOptions(nil, 0, 10, 20, nil, DrawSpriteOptions{})
}

func TestDrawSpriteOutOfRangeFrameIsNoop(t *testing.T) {
	DrawSprite(&SpriteSheet{sheet: &Sheet{}}, 999, 10, 20)
	DrawSpriteWithOptions(&SpriteSheet{sheet: &Sheet{}}, -1, 10, 20, DrawSpriteOptions{})
}
