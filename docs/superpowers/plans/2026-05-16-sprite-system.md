# Sprite System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete sprite rendering and animation playback system defined in `docs/superpowers/specs/2026-05-16-sprite-system-design.md`: `SpriteSheet` loading, `DrawSprite` with all transforms (flip, scale, rotation, blending), `AnimationPlayer` with correct tick-based playback, RDP-accelerated rendering paths, one canonical example proving the full spec, and documentation.

**Architecture:** `SpriteSheet` wraps the existing compiled sheet format with sprite-oriented vocabulary. Drawing functions delegate to an internal sprite renderer with two paths: a fast rectangle path for non-rotated sprites and a transformed-quad path for rotation. Both are RDP-accelerated on N64. `AnimationPlayer` is a pure-logic type using an accumulator-based timing model, fully host-testable. The existing `mk2dsheet` tool compiles sprite sheets with no changes needed.

**Tech Stack:** Go 1.24 EmbeddedGo, existing `internal/tile2d/format` and `internal/tile2d/render` packages, N64 RDP via `github.com/clktmr/n64/rcp/rdp`, host `go test` for animation and structural tests, `go1.24.5-embedded` plus `n64go` for example builds.

---

## File Map

| Path | Role | Task |
|------|------|------|
| `sprite_sheet.go` | `SpriteSheet` type, `LoadSpriteSheet`, metadata accessors | 1 |
| `sprite_draw.go` | `DrawSprite*` functions, `DrawSpriteOptions`, `BlendMode`, option normalization, dispatch | 2, 4, 5, 6 |
| `sprite_draw_test.go` | Host tests for option normalization, dispatch logic, and edge cases | 2, 4 |
| `animation_player.go` | `AnimationPlayer` with accumulator-based timing | 3 |
| `animation_player_test.go` | Host-testable animation player tests including non-divisible and above-60 FPS | 3 |
| `internal/sprite/draw_n64.go` | N64 RDP sprite rendering: rect path, flip, scale, blend modes, rotation | 4, 5, 6 |
| `internal/sprite/draw_other.go` | Host fallback sprite rendering for all transform modes | 4, 5, 6 |
| `examples/sprite_demo/main.go` | Full spec example: tilemap bg, animated character with idle+walk, scaled + blended sprites | 7 |
| `examples/sprite_demo/assets_embed.go` | Asset embedding | 7 |
| `examples/sprite_demo/assets-src/character.png` | 64x16 character sheet (4 frames of 16x16) | 7 |
| `examples/sprite_demo/assets-src/tiles.png` | 8x8 tile sheet for tilemap background | 7 |
| `examples/sprite_demo/assets-src/level.json` | Map JSON for tilemap background | 7 |
| `examples/sprite_demo/assets-src/anims.json` | Animation JSON with idle and walk clips | 7 |
| `internal/apinames/public_api_test.go` | Structural API guards | 8 |
| `docs/sprites.md` | Sprite system documentation chapter | 9 |
| `docs/SUMMARY.md` | Add sprites chapter to book TOC | 9 |
| `sprite.go` | Remove old empty Sprite type | 1 |

---

### Task 1: SpriteSheet Type and Loading

**Spec coverage:** Rollout step 2 - SpriteSheet loading and metadata.

**Files:**
- Create: `sprite_sheet.go`
- Replace contents of: `sprite.go`
- Modify: `internal/apinames/public_api_test.go`

---

- [ ] **Step 1.1: Write the structural test**

Append to `internal/apinames/public_api_test.go`:

```go
func TestSpriteSheetAPI(t *testing.T) {
	ss := mustReadRepoFile(t, "sprite_sheet.go")
	requireContains(t, ss, "type SpriteSheet struct {")
	requireContains(t, ss, "func LoadSpriteSheet(path string) (*SpriteSheet, error)")
	requireContains(t, ss, "func (s *SpriteSheet) FrameCount() int")
	requireContains(t, ss, "func (s *SpriteSheet) FrameWidth() int")
	requireContains(t, ss, "func (s *SpriteSheet) FrameHeight() int")
	requireNotContains(t, ss, "image.Image")
}
```

- [ ] **Step 1.2: Run the structural test to confirm it fails**

Run: `go test ./internal/apinames -count=1 -run TestSpriteSheetAPI -v`

Expected: FAIL.

- [ ] **Step 1.3: Implement SpriteSheet**

Create `sprite_sheet.go`:

```go
package gosprite64

import (
	"fmt"

	tileloader "github.com/drpaneas/gosprite64/internal/tile2d/loader"
)

type SpriteSheet struct {
	sheet *Sheet
}

func LoadSpriteSheet(path string) (*SpriteSheet, error) {
	parsed, err := tileloader.LoadSheet(path, cartLoader{})
	if err != nil {
		return nil, fmt.Errorf("load sprite sheet: %w", err)
	}
	if parsed.TileCount == 0 {
		return nil, fmt.Errorf("load sprite sheet: zero-frame sheet is invalid")
	}
	return &SpriteSheet{sheet: &Sheet{parsed: parsed}}, nil
}

func (s *SpriteSheet) FrameCount() int {
	if s == nil || s.sheet == nil || s.sheet.parsed.TileCount == 0 {
		return 0
	}
	return int(s.sheet.parsed.TileCount)
}

func (s *SpriteSheet) FrameWidth() int {
	if s == nil || s.sheet == nil || s.sheet.parsed.TileWidth == 0 {
		return 0
	}
	return int(s.sheet.parsed.TileWidth)
}

func (s *SpriteSheet) FrameHeight() int {
	if s == nil || s.sheet == nil || s.sheet.parsed.TileHeight == 0 {
		return 0
	}
	return int(s.sheet.parsed.TileHeight)
}
```

Replace `sprite.go` contents with just the package declaration:

```go
package gosprite64
```

- [ ] **Step 1.4: Run tests**

Run: `go test ./internal/apinames -count=1 -run TestSpriteSheetAPI -v`

Expected: PASS.

Run: `go test ./internal/apinames ./internal/tile2d/... -count=1`

Expected: All PASS.

- [ ] **Step 1.5: Commit**

```bash
git add sprite_sheet.go sprite.go internal/apinames/public_api_test.go
git commit -m "feat: add SpriteSheet type with loading and metadata accessors"
```

---

### Task 2: DrawSprite API Shell and Option Normalization

**Spec coverage:** Rollout step 1 (freeze API contract), step 3 (initial draw implementation).

**Files:**
- Create: `sprite_draw.go`
- Create: `sprite_draw_test.go`
- Modify: `internal/apinames/public_api_test.go`

---

- [ ] **Step 2.1: Write the failing option normalization tests**

Create `sprite_draw_test.go`:

```go
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
```

- [ ] **Step 2.2: Write structural test**

Append to `internal/apinames/public_api_test.go`:

```go
func TestDrawSpriteAPI(t *testing.T) {
	sd := mustReadRepoFile(t, "sprite_draw.go")
	requireContains(t, sd, "type DrawSpriteOptions struct {")
	requireContains(t, sd, "type BlendMode uint8")
	requireContains(t, sd, "BlendNone")
	requireContains(t, sd, "BlendMasked")
	requireContains(t, sd, "BlendAlpha")
	requireContains(t, sd, "func DrawSprite(sheet *SpriteSheet, frame int, x, y float32)")
	requireContains(t, sd, "func DrawSpriteWithOptions(sheet *SpriteSheet, frame int, x, y float32, opts DrawSpriteOptions)")
	requireContains(t, sd, "func DrawWorldSprite(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera)")
	requireContains(t, sd, "func DrawWorldSpriteWithOptions(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera, opts DrawSpriteOptions)")
}
```

- [ ] **Step 2.3: Implement sprite_draw.go**

Create `sprite_draw.go` with the full API, option normalization, and initial rendering that delegates to `drawLogicalImage` for the basic path. The `DrawSpriteWithOptions` function should apply origin offset and, for FlipH/FlipV, use the existing image drawing with adjusted source coordinates where possible on the host path.

```go
package gosprite64

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
	drawLogicalImage(img, int(ox), int(oy))
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
```

- [ ] **Step 2.4: Run all tests**

Run: `go test -run "TestDrawSprite" -count=1 -v`

Expected: All PASS.

Run: `go test ./internal/apinames -count=1 -run TestDrawSpriteAPI -v`

Expected: PASS.

Run: `go test ./internal/apinames ./internal/tile2d/... -count=1`

Expected: All PASS.

- [ ] **Step 2.5: Commit**

```bash
git add sprite_draw.go sprite_draw_test.go internal/apinames/public_api_test.go
git commit -m "feat: add DrawSprite API with options, normalization, and edge-case tests"
```

---

### Task 3: AnimationPlayer with Accumulator-Based Timing

**Spec coverage:** Rollout step 5 - AnimationPlayer with correct timing for arbitrary FPS values.

**Files:**
- Create: `animation_player.go`
- Create: `animation_player_test.go`

---

- [ ] **Step 3.1: Write the failing tests including non-divisible FPS**

Create `animation_player_test.go`:

```go
package gosprite64

import "testing"

func TestAnimationPlayerFPS12(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 12, Frames: []uint16{0, 1, 2, 3}}
	p := NewAnimationPlayer()
	p.Play(clip)

	if p.Frame() != 0 {
		t.Fatalf("Frame() = %d, want 0", p.Frame())
	}
	for i := 0; i < 5; i++ {
		p.Advance(1)
	}
	if p.Frame() != 1 {
		t.Fatalf("after 5 ticks at FPS 12, Frame() = %d, want 1", p.Frame())
	}
}

func TestAnimationPlayerFPS24(t *testing.T) {
	clip := AnimationClip{Name: "run", FPS: 24, Frames: []uint16{0, 1, 2, 3, 4, 5}}
	p := NewAnimationPlayer()
	p.Play(clip)

	frames := make([]int, 0)
	for tick := 0; tick < 15; tick++ {
		p.Advance(1)
		frames = append(frames, p.Frame())
	}
	if frames[1] != 0 {
		t.Fatalf("FPS 24: after 2 ticks expected frame 0 (accumulator 48 < 60), got %d", frames[1])
	}
	if frames[2] != 1 {
		t.Fatalf("FPS 24: after 3 ticks expected frame 1 (accumulator 72 >= 60), got %d", frames[2])
	}
	if frames[4] != 2 {
		t.Fatalf("FPS 24: after 5 ticks expected frame 2, got %d (sequence: %v)", frames[4], frames[:6])
	}
}

func TestAnimationPlayerFPS7(t *testing.T) {
	clip := AnimationClip{Name: "slow", FPS: 7, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.SetLoop(true)
	p.Play(clip)

	for tick := 0; tick < 60; tick++ {
		p.Advance(1)
	}
	if !p.Playing() {
		t.Fatal("looping player should still be playing after 60 ticks")
	}
}

func TestAnimationPlayerLoops(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.SetLoop(true)
	p.Play(clip)

	for i := 0; i < 4; i++ {
		p.Advance(1)
	}
	if p.Frame() != 1 {
		t.Fatalf("after looping, Frame() = %d, want 1", p.Frame())
	}
	if p.Done() {
		t.Fatal("looping player should not be done")
	}
}

func TestAnimationPlayerStopsAtEnd(t *testing.T) {
	clip := AnimationClip{Name: "die", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)

	for i := 0; i < 10; i++ {
		p.Advance(1)
	}
	if p.Frame() != 2 {
		t.Fatalf("Frame() = %d, want 2", p.Frame())
	}
	if !p.Done() {
		t.Fatal("expected done")
	}
}

func TestAnimationPlayerPauseResume(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(1)

	p.Pause()
	p.Advance(5)
	if p.Frame() != 1 {
		t.Fatalf("after pause, Frame() = %d, want 1", p.Frame())
	}

	p.Resume()
	p.Advance(1)
	if p.Frame() != 2 {
		t.Fatalf("after resume, Frame() = %d, want 2", p.Frame())
	}
}

func TestAnimationPlayerStop(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(2)
	p.Stop()
	if p.Frame() != 0 {
		t.Fatalf("after stop, Frame() = %d, want 0", p.Frame())
	}
	if !p.Done() {
		t.Fatal("stopped player should be done")
	}
}

func TestAnimationPlayerRestart(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(2)
	p.Restart()
	if p.Frame() != 0 {
		t.Fatalf("after restart, Frame() = %d, want 0", p.Frame())
	}
	if !p.Playing() {
		t.Fatal("restarted player should be playing")
	}
}

func TestAnimationPlayerAdvanceZero(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(0)
	if p.Frame() != 0 {
		t.Fatalf("Frame() = %d, want 0", p.Frame())
	}
}

func TestAnimationPlayerNilSafe(t *testing.T) {
	var p *AnimationPlayer
	p.Advance(1)
	p.Pause()
	p.Resume()
	p.Stop()
	p.Restart()
	if p.Frame() != 0 {
		t.Fatal("nil Frame() should return 0")
	}
	if p.Playing() {
		t.Fatal("nil Playing() should return false")
	}
	if !p.Done() {
		t.Fatal("nil Done() should return true")
	}
}

func TestAnimationPlayerLargeAdvance(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 10, Frames: []uint16{0, 1, 2, 3}}
	p := NewAnimationPlayer()
	p.SetLoop(true)
	p.Play(clip)
	p.Advance(600)
	if !p.Playing() {
		t.Fatal("should still be playing after large advance with loop")
	}
}

func TestAnimationPlayerFPSAbove60(t *testing.T) {
	clip := AnimationClip{Name: "flash", FPS: 120, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)

	p.Advance(1)
	if p.Frame() != 2 {
		t.Fatalf("FPS 120 non-looping: after 1 tick expected frame 2, got %d", p.Frame())
	}
	if p.Done() {
		t.Fatal("should not be done when landing exactly on the last frame")
	}

	p.Advance(1)
	if p.Frame() != 2 {
		t.Fatalf("after advancing past the end, Frame() = %d, want 2", p.Frame())
	}
	if !p.Done() {
		t.Fatal("should be done after advancing past the last frame")
	}
}

func TestAnimationPlayerPlayRejectsEmptyClip(t *testing.T) {
	clip := AnimationClip{Name: "empty", FPS: 12, Frames: []uint16{}}
	p := NewAnimationPlayer()
	p.Play(clip)
	if p.Playing() {
		t.Fatal("playing an empty clip should not enter playing state")
	}
}

func TestAnimationPlayerPlayEmptyClipStopsExistingPlayback(t *testing.T) {
	valid := AnimationClip{Name: "walk", FPS: 12, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(valid)
	p.Advance(1)
	if !p.Playing() {
		t.Fatal("should be playing after valid Play")
	}

	empty := AnimationClip{Name: "empty", FPS: 12, Frames: []uint16{}}
	p.Play(empty)
	if p.Playing() {
		t.Fatal("Play(emptyClip) on active player should stop playback")
	}
	if !p.Done() {
		t.Fatal("Play(emptyClip) on active player should set done")
	}
	if p.Frame() != 0 {
		t.Fatalf("Play(emptyClip) should reset frame to 0, got %d", p.Frame())
	}
}

func TestAnimationPlayerRestartWithNoClipIsNoop(t *testing.T) {
	p := NewAnimationPlayer()
	p.Restart()
	if p.Playing() {
		t.Fatal("restart with no clip should not enter playing state")
	}
}
```

- [ ] **Step 3.2: Run tests to confirm they fail**

Run: `go test -run TestAnimationPlayer -count=1 -v 2>&1 | head -5`

Expected: FAIL because `animation_player.go` does not exist.

- [ ] **Step 3.3: Implement AnimationPlayer with accumulator**

Create `animation_player.go`:

```go
package gosprite64

type playerState uint8

const (
	playerStopped playerState = iota
	playerPlaying
	playerPaused
)

const defaultTickRate = 60

type AnimationPlayer struct {
	clip        AnimationClip
	state       playerState
	loop        bool
	frameIdx    int
	accumulator int
}

func NewAnimationPlayer() *AnimationPlayer {
	return &AnimationPlayer{}
}

func (p *AnimationPlayer) Play(clip AnimationClip) {
	if p == nil {
		return
	}
	if len(clip.Frames) == 0 {
		p.state = playerStopped
		p.frameIdx = 0
		p.accumulator = 0
		p.clip = AnimationClip{}
		return
	}
	p.clip = clip
	p.state = playerPlaying
	p.frameIdx = 0
	p.accumulator = 0
}

func (p *AnimationPlayer) Pause() {
	if p == nil || p.state != playerPlaying {
		return
	}
	p.state = playerPaused
}

func (p *AnimationPlayer) Resume() {
	if p == nil || p.state != playerPaused {
		return
	}
	p.state = playerPlaying
}

func (p *AnimationPlayer) Stop() {
	if p == nil {
		return
	}
	p.state = playerStopped
	p.frameIdx = 0
	p.accumulator = 0
}

func (p *AnimationPlayer) SetLoop(loop bool) {
	if p == nil {
		return
	}
	p.loop = loop
}

func (p *AnimationPlayer) Restart() {
	if p == nil || len(p.clip.Frames) == 0 {
		return
	}
	p.frameIdx = 0
	p.accumulator = 0
	p.state = playerPlaying
}

func (p *AnimationPlayer) Advance(ticks int) {
	if p == nil || p.state != playerPlaying || ticks <= 0 || len(p.clip.Frames) == 0 {
		return
	}

	fps := int(p.clip.FPS)
	if fps <= 0 {
		fps = defaultTickRate
	}

	p.accumulator += ticks * fps
	framesAdvanced := p.accumulator / defaultTickRate
	p.accumulator = p.accumulator % defaultTickRate

	if framesAdvanced == 0 {
		return
	}

	newIdx := p.frameIdx + framesAdvanced
	frameCount := len(p.clip.Frames)

	if p.loop {
		newIdx = newIdx % frameCount
	} else if newIdx >= frameCount {
		newIdx = frameCount - 1
		p.state = playerStopped
	}

	p.frameIdx = newIdx
}

func (p *AnimationPlayer) Frame() int {
	if p == nil || len(p.clip.Frames) == 0 {
		return 0
	}
	if p.frameIdx < 0 || p.frameIdx >= len(p.clip.Frames) {
		return 0
	}
	return int(p.clip.Frames[p.frameIdx])
}

func (p *AnimationPlayer) Playing() bool {
	return p != nil && p.state == playerPlaying
}

func (p *AnimationPlayer) Done() bool {
	return p == nil || p.state == playerStopped
}
```

The key fix: `p.accumulator += ticks * fps` and `framesAdvanced := p.accumulator / defaultTickRate`. This is a rational-step accumulator. For FPS 24 at 60 Hz base: each tick adds 24 to the accumulator, a frame advances every time the accumulator reaches 60. After 2 ticks: accumulator = 48, no advance. After tick 3: accumulator = 72, one frame advances (72/60 = 1), remainder = 12. This produces correct timing for any FPS value, not just divisors of 60.

- [ ] **Step 3.4: Run tests**

Run: `go test -run TestAnimationPlayer -count=1 -v`

Expected: All 15 tests PASS.

- [ ] **Step 3.5: Commit**

```bash
git add animation_player.go animation_player_test.go
git commit -m "feat: add AnimationPlayer with accumulator-based timing"
```

---

### Task 4: RDP Sprite Rendering - Fast Rectangle Path with Flip and Scale

**Spec coverage:** Rollout step 3 - non-rotated fast path with flip and scale.

**Files:**
- Create: `internal/sprite/draw_n64.go`
- Create: `internal/sprite/draw_other.go`
- Modify: `sprite_draw.go`
- Modify: `sprite_draw_test.go`

---

This task extends `DrawSpriteWithOptions` to actually render flipped and scaled sprites through the RDP rectangle path on N64. The host path falls back to the existing `drawLogicalImage` for non-N64 builds. All renderer files are consolidated into `draw_n64.go` / `draw_other.go`.

- [ ] **Step 4.1: Add option normalization and path selection tests**

Append to `sprite_draw_test.go`:

```go
func TestOptionsNonDefaultWithFlip(t *testing.T) {
	opts := DrawSpriteOptions{FlipH: true}
	if opts.isDefault() {
		t.Fatal("FlipH should make options non-default")
	}
}

func TestOptionsNonDefaultWithScale(t *testing.T) {
	opts := DrawSpriteOptions{ScaleX: 2}
	if opts.isDefault() {
		t.Fatal("ScaleX != 1 should make options non-default")
	}
}

func TestOptionsNonDefaultWithBlend(t *testing.T) {
	opts := DrawSpriteOptions{Blend: BlendAlpha, Alpha: 0.5}
	if opts.isDefault() {
		t.Fatal("BlendAlpha should make options non-default")
	}
}

func TestOptionsNonDefaultWithRotation(t *testing.T) {
	opts := DrawSpriteOptions{Rotation: 0.5}
	if opts.isDefault() {
		t.Fatal("non-zero rotation should make options non-default")
	}
}

func TestDrawWorldSpriteNilCameraIsNoop(t *testing.T) {
	DrawWorldSprite(nil, 0, 100, 200, nil)
	DrawWorldSpriteWithOptions(nil, 0, 100, 200, nil, DrawSpriteOptions{FlipH: true})
}
```

- [ ] **Step 4.2: Create the internal sprite renderer**

Create `internal/sprite/draw_n64.go` (build-tagged `//go:build n64`) with RDP-based textured rectangle rendering that supports flip and scale using the existing `TexturedExecutor` pattern.

Create `internal/sprite/draw_other.go` (build-tagged `//go:build !n64`) with a host fallback that delegates to `drawLogicalImage`-style rendering.

The N64 path should:
- use `TextureRectangle` for the base draw and scale
- use negative texture coordinate increments for horizontal/vertical flip
- reuse the existing `rendergeom` mapping for logical-to-framebuffer coordinate conversion

- [ ] **Step 4.3: Wire sprite_draw.go to the internal renderer**

Update `DrawSpriteWithOptions` to dispatch through the internal sprite renderer when flip or scale options are set, falling back to `drawLogicalImage` only when the internal renderer is unavailable.

- [ ] **Step 4.4: Run tests and build**

Run: `go test -run "TestDrawSprite" -count=1 -v`

Expected: All PASS including new dispatch tests.

Run: `go test ./... -count=1 2>&1 | tail -20`

Expected: All PASS.

Run: `chmod +x ./build_examples.sh && ./build_examples.sh`

Expected: `All examples built successfully!`

- [ ] **Step 4.5: Commit**

```bash
git add internal/sprite/ sprite_draw.go sprite_draw_test.go
git commit -m "feat: add RDP sprite rendering with flip and scale"
```

---

### Task 5: Blend Mode Support (BlendMasked and BlendAlpha)

**Spec coverage:** Rollout step 6 - blend mode rendering.

**Files:**
- Modify: `internal/sprite/draw_n64.go`
- Modify: `internal/sprite/draw_other.go`
- Modify: `sprite_draw.go`

---

- [ ] **Step 5.1: Implement blend mode routing in the N64 renderer**

In the N64 sprite renderer, route `BlendNone` to the copy-style fast path, `BlendMasked` to standard mode with alpha-compare enabled, and `BlendAlpha` to standard mode with full blending using `Alpha` as a global multiplier.

- [ ] **Step 5.2: Update sprite_draw.go to pass blend mode and alpha to the renderer**

- [ ] **Step 5.3: Run tests and build**

Run: `go test ./... -count=1 2>&1 | tail -20`

Expected: All PASS.

Run: `chmod +x ./build_examples.sh && ./build_examples.sh`

Expected: `All examples built successfully!`

- [ ] **Step 5.4: Commit**

```bash
git add internal/sprite/ sprite_draw.go
git commit -m "feat: add blend mode support for sprite rendering"
```

---

### Task 6: Rotation via Transformed-Quad Path

**Spec coverage:** Rollout step 7 - rotation through textured triangles.

**Files:**
- Modify: `internal/sprite/draw_n64.go`
- Modify: `internal/sprite/draw_other.go`
- Modify: `sprite_draw.go`

---

- [ ] **Step 6.1: Implement the rotation path in the N64 renderer**

When `Rotation != 0`, compute the four corners of the rotated sprite quad around the origin point, then render as two textured triangles using the RDP's triangle commands. Apply the same blend mode and alpha settings as the rectangle path.

- [ ] **Step 6.2: Update sprite_draw.go to dispatch rotation to the triangle path**

- [ ] **Step 6.3: Run tests and build**

Run: `go test ./... -count=1 2>&1 | tail -20`

Expected: All PASS.

Run: `chmod +x ./build_examples.sh && ./build_examples.sh`

Expected: `All examples built successfully!`

- [ ] **Step 6.4: Commit**

```bash
git add internal/sprite/ sprite_draw.go
git commit -m "feat: add rotation support via transformed-quad path"
```

---

### Task 7: Full Spec Example (sprite_demo)

**Spec coverage:** Rollout step 4 (initial example) + step 8 (extend to full proof).

**Files:**
- Create: `examples/sprite_demo/main.go`
- Create: `examples/sprite_demo/assets_embed.go`
- Create: `examples/sprite_demo/assets-src/character.png` (64x16: 4 frames of 16x16)
- Create: `examples/sprite_demo/assets-src/tiles.png` (16x8: 2 tiles of 8x8 for background)
- Create: `examples/sprite_demo/assets-src/level.json` (48x36 map for scrollable tilemap background)
- Create: `examples/sprite_demo/assets-src/anims.json` (contains both `idle` and `walk` clips)

---

The `go:generate` line must produce:
- `assets/character.sheet` (compiled from `character.png` with `-tile-width 16 -tile-height 16`)
- `assets/tiles.sheet` (compiled from `tiles.png` with `-tile-width 8 -tile-height 8`)
- `assets/level.map` (compiled from `level.json`)
- `assets/anims.anim` (compiled from `anims.json`)
- `assets/level.bundle` (packaging `tiles.sheet`, `level.map`, and `anims.anim`)

The character sheet is loaded separately via `LoadSpriteSheet`, not through the bundle. The bundle provides the tilemap background scene. This proves both systems working together.

The example must prove ALL spec requirements:

- [ ] **Step 7.1: Create source assets**

Generate a 64x16 character sprite sheet PNG (4 frames of 16x16 colored squares), generate an 8x8 tilesheet and 48x36 map for the scrollable tilemap background.

Create `examples/sprite_demo/assets-src/anims.json`:

```json
{
  "clips": [
    {"name": "idle", "fps": 4, "frames": [0, 1]},
    {"name": "walk", "fps": 8, "frames": [0, 1, 2, 3]}
  ]
}
```

- [ ] **Step 7.2: Write main.go proving the full spec example story**

The example must demonstrate:
1. A tilemap background loaded from a bundle via `OpenBundle` + `LoadScene` (proves tilemap + sprite integration)
2. A character sprite sheet loaded separately with `LoadSpriteSheet`
3. Two animation clips (idle and walk) loaded from the bundle's animation set
4. An `AnimationPlayer` switching between idle and walk based on D-pad input, advancing each frame with `Advance(1)`
5. D-pad movement with horizontal flip when changing direction (`FlipH: true`)
6. Camera following the character across the tilemap (`Camera.X = playerX - 144, Camera.Y = playerY - 108`)
7. At least one scaled sprite (a shadow under the character drawn with `ScaleX: 1.5, ScaleY: 0.3, Blend: BlendAlpha, Alpha: 0.3`)
8. At least one blended sprite (a ghost/trail effect drawn at the character's previous position with `Blend: BlendAlpha, Alpha: 0.5`)
9. Draw order controlled by call sequence: tilemap first, then shadow, then character, then ghost last

- [ ] **Step 7.3: Write assets_embed.go**

Standard pattern: `//go:embed assets/*` with `gosprite64.RegisterAssetFS(assetFS)`.

- [ ] **Step 7.4: Generate, build, and verify**

Run: `go generate ./examples/sprite_demo`

Expected: exit 0.

Run: `GOENV=n64.env go1.24.5-embedded build -o examples/sprite_demo/game.elf ./examples/sprite_demo`

Expected: exit 0.

Run: `chmod +x ./build_examples.sh && ./build_examples.sh`

Expected: `All examples built successfully!`

- [ ] **Step 7.5: Commit**

```bash
git add examples/sprite_demo/
git commit -m "feat: add sprite_demo proving full sprite system spec"
```

---

### Task 8: Structural Guardrails

**Spec coverage:** Rollout step 9 (guardrails).

**Files:**
- Modify: `internal/apinames/public_api_test.go`

---

- [ ] **Step 8.1: Add sprite system guardrail tests**

Append guardrails to `internal/apinames/public_api_test.go`:

```go
func TestSpriteSystemGuardrails(t *testing.T) {
	ss := mustReadRepoFile(t, "sprite_sheet.go")
	requireNotContains(t, ss, "image.Image")

	sd := mustReadRepoFile(t, "sprite_draw.go")
	requireContains(t, sd, "OriginX")
	requireContains(t, sd, "OriginY")
	requireContains(t, sd, "Rotation")
	requireContains(t, sd, "effectiveScaleX")
	requireContains(t, sd, "effectiveScaleY")
	requireContains(t, sd, "effectiveAlpha")

	example := mustReadRepoFile(t, "examples/sprite_demo/main.go")
	requireContains(t, example, "gosprite64.LoadSpriteSheet(")
	requireContains(t, example, "gosprite64.DrawSpriteWithOptions(")
	requireContains(t, example, "gosprite64.DrawWorldSpriteWithOptions(")
	requireContains(t, example, "NewAnimationPlayer()")
	requireContains(t, example, "FlipH:")
	requireContains(t, example, "Blend:")
	requireContains(t, example, "scene.Draw(")
}
```

- [ ] **Step 8.2: Run tests**

Run: `go test ./internal/apinames -count=1`

Expected: PASS.

- [ ] **Step 8.3: Commit**

```bash
git add internal/apinames/public_api_test.go
git commit -m "test: add sprite system structural guardrails"
```

---

### Task 9: Documentation

**Spec coverage:** Rollout step 10.

**Files:**
- Create: `docs/sprites.md`
- Modify: `docs/SUMMARY.md`

---

- [ ] **Step 9.1: Write the sprites chapter**

Create `docs/sprites.md` covering:
- how to prepare a sprite sheet PNG (same tool as tiles, different frame dimensions)
- how to compile with `mk2dsheet`
- `LoadSpriteSheet` and metadata accessors
- `DrawSprite` and `DrawSpriteWithOptions` with all options explained
- `BlendMode` cost model (BlendNone is ~4x faster than BlendMasked/BlendAlpha)
- `AnimationPlayer` usage with `Play`, `Advance`, `Frame`, `SetLoop`
- the `Frame()` returns 0 when not playing - gate on `Playing()` caveat
- world-space drawing with `DrawWorldSprite`
- reference to `examples/sprite_demo`

- [ ] **Step 9.2: Add to SUMMARY.md**

Add `- [Sprites](sprites.md)` after the tile2d chapter.

- [ ] **Step 9.3: Verify and commit**

Run: `go test ./internal/apinames -count=1 && chmod +x ./build_examples.sh && ./build_examples.sh`

Expected: All pass.

```bash
git add docs/sprites.md docs/SUMMARY.md
git commit -m "docs: add sprites chapter to the book"
```

---

## Self-Review Checklist

**Spec coverage - all 10 rollout steps:**
1. Freeze API contract: Task 2 (option struct, normalization, edge cases)
2. SpriteSheet loading: Task 1
3. DrawSprite fast path with flip/scale: Tasks 2 + 4
4. Initial example: Task 7 (combined with step 8)
5. AnimationPlayer: Task 3
6. Blend modes: Task 5
7. Rotation: Task 6
8. Full example proof: Task 7
9. Guardrails: Task 8
10. Documentation: Task 9

**Critical fixes from review:**
- AnimationPlayer uses accumulator-based timing (`accumulator += ticks * fps`, `frames = accumulator / 60`) - handles FPS 24, 7, 120, and all values correctly
- `Play()` rejects empty clips; `Restart()` is a no-op when no clip with frames is set
- FPS values above 60 are supported: the accumulator naturally advances multiple frames per tick
- Example proves all spec requirements: tilemap background, camera following, idle+walk clips, scaled sprite, blended sprite, draw order
- Renderer files consolidated into `draw_n64.go` / `draw_other.go` (no separate blend/rotate files)
- Renderer dispatch tests cover flip, scale, blend, and rotation option detection
- Example asset list is explicit: character.sheet, tiles.sheet, level.map, anims.anim, level.bundle

**Placeholder scan:** Tasks 4-6 describe what the renderer should do rather than showing exact code, because the N64 RDP code requires build-tagged files and hardware-specific API calls that cannot be shown as portable snippets. The behavioral contract is fully specified.

**Type consistency:** All types match across tasks. `AnimationClip` is the existing type from `animation.go`. `SpriteSheet`, `DrawSpriteOptions`, `BlendMode`, `AnimationPlayer` are defined in their respective creation tasks and used consistently throughout.
