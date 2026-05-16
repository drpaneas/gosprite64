package gosprite64

import (
	"github.com/drpaneas/gosprite64/internal/sprite"
)

type BlendMode uint8

const (
	BlendNone   BlendMode = iota
	BlendMasked
	BlendAlpha
)

type DrawSpriteOptions struct {
	FlipH    bool
	FlipV    bool
	ScaleX   float32
	ScaleY   float32
	Rotation float32
	OriginX  float32
	OriginY  float32
	Blend    BlendMode
	Alpha    float32
}

func (o DrawSpriteOptions) effectiveScaleX() float32 {
	if o.ScaleX == 0 {
		return 1
	}
	return o.ScaleX
}

func (o DrawSpriteOptions) effectiveScaleY() float32 {
	if o.ScaleY == 0 {
		return 1
	}
	return o.ScaleY
}

func (o DrawSpriteOptions) effectiveAlpha() float32 {
	if o.Alpha == 0 {
		return 1
	}
	return o.Alpha
}

func (o DrawSpriteOptions) isDefault() bool {
	return !o.FlipH && !o.FlipV &&
		o.effectiveScaleX() == 1 && o.effectiveScaleY() == 1 &&
		o.Rotation == 0 &&
		o.OriginX == 0 && o.OriginY == 0 &&
		o.Blend == BlendNone
}

func DrawSprite(sheet *SpriteSheet, frame int, x, y float32) {
	if sheet == nil || frame < 0 || frame >= sheet.FrameCount() {
		return
	}
	img := sheet.sheet.tileImage(uint16(frame + 1))
	if img == nil {
		return
	}
	drawLogicalImage(img, int(x), int(y))
}

func DrawSpriteWithOptions(sheet *SpriteSheet, frame int, x, y float32, opts DrawSpriteOptions) {
	if sheet == nil || frame < 0 || frame >= sheet.FrameCount() {
		return
	}
	if opts.isDefault() {
		DrawSprite(sheet, frame, x, y)
		return
	}
	img := sheet.sheet.tileImage(uint16(frame + 1))
	if img == nil {
		return
	}
	sx := opts.effectiveScaleX()
	sy := opts.effectiveScaleY()
	ox := x - opts.OriginX*sx
	oy := y - opts.OriginY*sy

	video := currentVideo()
	if video == nil || video.Framebuffer == nil {
		return
	}
	sprite.RenderSprite(video.Framebuffer, img, int(ox), int(oy),
		opts.FlipH, opts.FlipV, sx, sy)
}

func DrawWorldSprite(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera) {
	if cam == nil {
		DrawSprite(sheet, frame, worldX, worldY)
		return
	}
	DrawSprite(sheet, frame, worldX-float32(cam.X), worldY-float32(cam.Y))
}

func DrawWorldSpriteWithOptions(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera, opts DrawSpriteOptions) {
	if cam == nil {
		DrawSpriteWithOptions(sheet, frame, worldX, worldY, opts)
		return
	}
	DrawSpriteWithOptions(sheet, frame, worldX-float32(cam.X), worldY-float32(cam.Y), opts)
}
