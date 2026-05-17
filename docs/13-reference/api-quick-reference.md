# API Quick Reference

A compact listing of every public function, type, and constant in the GoSprite64 API.

## Core

| Symbol | Description |
|--------|-------------|
| `Run(g Game)` | Starts the game loop with fixed-step timing at 60 FPS |
| `Game` (interface) | Implement `Init()`, `Update()`, `Draw()` to define your game |
| `TargetFPS` (var, int) | Target frame rate, defaults to 60 |
| `RegisterAssetFS(f cartfs.FS)` | Registers the embedded cartridge filesystem for asset loading |
| `LoadFromCartridge(filename string) ([]byte, error)` | Reads a raw file from the cartridge filesystem |

## Drawing

| Symbol | Description |
|--------|-------------|
| `ClearScreen()` | Fills the screen with Black |
| `ClearScreenWith(c color.Color)` | Fills the screen with the given color |
| `FillRect(x1, y1, x2, y2 int, c color.Color)` | Draws a filled rectangle (inclusive coordinates) |
| `DrawRect(x1, y1, x2, y2 int, c color.Color)` | Draws a rectangle outline |
| `DrawLine(x1, y1, x2, y2 int, c color.Color)` | Draws a 1-pixel line (Bresenham for diagonals) |
| `DrawImage(src image.Image, x, y int)` | Draws a Go `image.Image` at screen coordinates |
| `DrawWorldImage(src image.Image, worldX, worldY int, cam *Camera)` | Draws an image in world space, offset by camera |

## Colors

All 16 predefined palette colors are exported as `color.Color` variables:

| Variable | RGB |
|----------|-----|
| `Black` | (0, 0, 0) |
| `DarkBlue` | (29, 43, 83) |
| `DarkPurple` | (126, 37, 83) |
| `DarkGreen` | (0, 135, 81) |
| `Brown` | (171, 82, 54) |
| `DarkGray` | (95, 87, 79) |
| `LightGray` | (194, 195, 199) |
| `White` | (255, 241, 232) |
| `Red` | (255, 0, 77) |
| `Orange` | (255, 163, 0) |
| `Yellow` | (255, 236, 39) |
| `Green` | (0, 228, 54) |
| `Blue` | (41, 173, 255) |
| `Indigo` | (131, 118, 156) |
| `Pink` | (255, 119, 168) |
| `Peach` | (255, 204, 170) |

## Sprites

| Symbol | Description |
|--------|-------------|
| `SpriteSheet` (struct) | A loaded sprite sheet containing animation frames |
| `LoadSpriteSheet(path string) (*SpriteSheet, error)` | Loads a sprite sheet from the cartridge filesystem |
| `(*SpriteSheet).FrameCount() int` | Returns the total number of frames |
| `(*SpriteSheet).FrameWidth() int` | Returns the width of each frame in pixels |
| `(*SpriteSheet).FrameHeight() int` | Returns the height of each frame in pixels |
| `DrawSprite(sheet *SpriteSheet, frame int, x, y float32)` | Draws a sprite frame at screen coordinates |
| `DrawSpriteWithOptions(sheet *SpriteSheet, frame int, x, y float32, opts DrawSpriteOptions)` | Draws a sprite with flip, scale, rotation, blend |
| `DrawWorldSprite(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera)` | Draws a sprite in world space |
| `DrawWorldSpriteWithOptions(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera, opts DrawSpriteOptions)` | World-space sprite with options |
| `DrawSpriteOptions` (struct) | FlipH, FlipV, ScaleX, ScaleY, Rotation, OriginX, OriginY, Blend, Alpha |
| `BlendMode` (uint8) | Blend mode enum type |
| `BlendNone` | No blending (fastest, opaque blit) |
| `BlendMasked` | Binary alpha (pixels are fully opaque or fully transparent) |
| `BlendAlpha` | Per-pixel alpha blending |

## Sheet & Tile

| Symbol | Description |
|--------|-------------|
| `Sheet` (struct) | A tile sheet loaded from a bundle |
| `SheetInfo` (struct) | TileWidth, TileHeight, TileCount, AtlasWidth, AtlasHeight |
| `(*Sheet).Info() SheetInfo` | Returns tile dimensions and atlas size |
| `(*Sheet).Tile(tileID uint16) image.Image` | Returns the image for a single tile |

## Animation

| Symbol | Description |
|--------|-------------|
| `AnimationSet` (struct) | A named collection of animation clips |
| `AnimationClip` (struct) | Name, FPS, and frame indices for one clip |
| `(*AnimationSet).Name() string` | Returns the animation set's name |
| `(*AnimationSet).Clips() []AnimationClip` | Returns all clips in the set |
| `(*AnimationSet).Clip(name string) (AnimationClip, bool)` | Looks up a clip by name |
| `AnimationPlayer` (struct) | Drives frame-by-frame playback of a clip |
| `NewAnimationPlayer() *AnimationPlayer` | Creates a new player (stopped) |
| `(*AnimationPlayer).Play(clip AnimationClip)` | Starts playing a clip from frame 0 |
| `(*AnimationPlayer).Pause()` | Pauses playback |
| `(*AnimationPlayer).Resume()` | Resumes a paused player |
| `(*AnimationPlayer).Stop()` | Stops playback and resets to frame 0 |
| `(*AnimationPlayer).Restart()` | Replays the current clip from the start |
| `(*AnimationPlayer).SetLoop(loop bool)` | Enables or disables looping |
| `(*AnimationPlayer).Advance(ticks int)` | Advances playback by the given number of ticks |
| `(*AnimationPlayer).Frame() int` | Returns the current frame index |
| `(*AnimationPlayer).Playing() bool` | True if the player is actively playing |
| `(*AnimationPlayer).Done() bool` | True if playback has stopped |

## Fonts & Text

| Symbol | Description |
|--------|-------------|
| `DrawText(str string, x, y int, c color.Color)` | Draws text using the built-in 8x8 font |
| `Font` (struct) | Custom font with per-glyph metrics from a sprite sheet |
| `Glyph` (struct) | Frame, Width, Advance, OffsetX, OffsetY |
| `NewFont(sheet *SpriteSheet, glyphs map[rune]Glyph, lineHeight int) *Font` | Creates a custom font |
| `(*Font).LineHeight() int` | Vertical distance between baselines |
| `(*Font).GlyphFor(r rune) (Glyph, bool)` | Looks up a glyph, falls back to `Font.Fallback` |
| `(*Font).MeasureText(text string) (width, height int)` | Measures pixel size of rendered text |
| `(*Font).DrawTextEx(text string, x, y int, align TextAlign)` | Draws text with alignment |
| `(*Font).WrapText(text string, maxWidth int) string` | Inserts newlines to wrap text at a pixel width |
| `FormatScore(score int, width int) string` | Formats an integer with leading zeros |
| `TextAlign` (int) | Horizontal alignment enum |
| `AlignLeft` | Left-aligned (default) |
| `AlignCenter` | Center-aligned |
| `AlignRight` | Right-aligned |

## Input

| Symbol | Description |
|--------|-------------|
| `IsButtonDown(button ButtonMask) bool` | True if button is held (port 0) |
| `IsButtonJustPressed(button ButtonMask) bool` | True on the frame a button transitions to pressed (port 0) |
| `StickPosition(deadzone float64) (float64, float64)` | Analog stick X/Y in [-1, 1] (port 0) |
| `PlayerButtonDown(port int, button ButtonMask) bool` | Per-port button held check |
| `PlayerButtonJustPressed(port int, button ButtonMask) bool` | Per-port just-pressed check |
| `PlayerStickPosition(port int, deadzone float64) (float64, float64)` | Per-port analog stick |
| `IsControllerConnected(port int) bool` | True if a controller is plugged into the given port |
| `ConnectedControllers() int` | Number of connected controllers |
| `SetRumble(port int, enabled bool)` | Enables or disables the rumble pak |
| `ButtonMask` (type alias) | Bitmask type for button constants |
| `MaxControllers` (const, 4) | Number of controller ports |

**Button constants:** `ButtonA`, `ButtonB`, `ButtonZ`, `ButtonStart`, `ButtonDPadUp`, `ButtonDPadDown`, `ButtonDPadLeft`, `ButtonDPadRight`, `ButtonL`, `ButtonR`, `ButtonCUp`, `ButtonCDown`, `ButtonCLeft`, `ButtonCRight`

## Input Replay

| Symbol | Description |
|--------|-------------|
| `FrameInput` (struct) | Buttons, StickX, StickY for one player/frame |
| `ReplayData` (struct) | PlayerCount, FrameCount, and recorded frames |
| `InputRecorder` (struct) | Records per-frame input during gameplay |
| `NewInputRecorder(playerCount int) *InputRecorder` | Creates a recorder |
| `(*InputRecorder).CaptureFrame(player int, input FrameInput)` | Records one frame |
| `(*InputRecorder).Finish() *ReplayData` | Finalizes recording |
| `InputPlayer` (struct) | Replays recorded input |
| `NewInputPlayer(data *ReplayData) *InputPlayer` | Creates a replay player |
| `(*InputPlayer).NextFrame(player int) (FrameInput, bool)` | Gets next frame for a player |
| `(*InputPlayer).Done() bool` | True when all frames consumed |
| `(*InputPlayer).Reset()` | Restarts playback from the beginning |
| `(*InputPlayer).CurrentFrame() int` | Current playback position |

## Audio

| Symbol | Description |
|--------|-------------|
| `AudioAsset` (struct) | Describes a single audio asset (ID, rate, loop points, etc.) |
| `AudioBundle` (struct) | Collection of audio assets with data and name resolver |
| `RegisterAudioBundle(bundle AudioBundle)` | Registers audio assets before `Run()` |
| `PlaySoundEffect(id sfx.ID) bool` | Plays a sound effect, returns true if queued |
| `PlayMusic(id music.ID) bool` | Starts music playback |
| `StopMusic()` | Stops the current music track |
| `SetSoundEffectVolume(v float32)` | Sets SFX volume (0.0-1.0) |
| `SetMusicVolume(v float32)` | Sets music volume (0.0-1.0) |
| `DefaultAudioOutputRate` (const, 48000) | DAC output rate in Hz |

## Tile Scenes

| Symbol | Description |
|--------|-------------|
| `Bundle` (struct) | A loaded asset bundle containing sheets, maps, and animations |
| `OpenBundle(path string) (*Bundle, error)` | Opens a bundle from the cartridge filesystem |
| `(*Bundle).LoadSheet(name string) (*Sheet, error)` | Loads a named sheet from the bundle |
| `(*Bundle).LoadMap(name string) (*Map, error)` | Loads a named map from the bundle |
| `(*Bundle).LoadAnimation(name string) (*AnimationSet, error)` | Loads a named animation set |
| `Scene` (struct) | A fully loaded tile scene with map, sheets, and animations |
| `LoadScene(bundle *Bundle) (*Scene, error)` | Loads all assets from a bundle into a renderable scene |
| `(*Scene).Draw(cam *Camera)` | Renders all visible layers with the given camera |
| `(*Scene).Map() *Map` | Returns the scene's map |
| `(*Scene).Sheet(index int) *Sheet` | Returns a sheet by index |
| `(*Scene).SheetByID(id uint16) *Sheet` | Returns a sheet by 1-based ID |
| `(*Scene).SheetCount() int` | Number of sheets in the scene |
| `(*Scene).Animation(index int) *AnimationSet` | Returns an animation by index |
| `(*Scene).AnimationByName(name string) *AnimationSet` | Looks up an animation set by name |
| `(*Scene).AnimationCount() int` | Number of animation sets |
| `(*Scene).LayerSheet(layer int) (*Sheet, bool)` | Returns the sheet assigned to a layer |
| `(*Scene).LayerAssets(layer int) (MapLayerInfo, *Sheet, bool)` | Returns layer info and sheet together |
| `(*Scene).LayerSheetInfo(layer int) (SheetInfo, bool)` | Returns the SheetInfo for a layer |
| `(*Scene).Stats() RuntimeStats` | Returns rendering statistics (allocation-free) |
| `Map` (struct) | Tile map with layers of cell data |
| `MapLayerInfo` (struct) | SheetID, NonZeroTiles |
| `(*Map).Width() int` | Map width in tiles |
| `(*Map).Height() int` | Map height in tiles |
| `(*Map).TileWidth() int` | Width of each tile in pixels |
| `(*Map).TileHeight() int` | Height of each tile in pixels |
| `(*Map).PixelWidth() int` | Total map width in pixels |
| `(*Map).PixelHeight() int` | Total map height in pixels |
| `(*Map).LayerCount() int` | Number of layers |
| `(*Map).LayerInfo(layer int) (MapLayerInfo, bool)` | Returns cached info for a layer (O(1)) |
| `(*Map).LayerSheetID(layer int) (uint16, bool)` | Returns the sheet ID for a layer |
| `(*Map).TileAt(layer, x, y int) (uint16, bool)` | Returns the tile ID at a grid position |
| `RuntimeStats` (struct) | SheetRAMBytes, MapRAMBytes, CachedChunks, VisibleTiles, SheetCount, LayerCount, UploadCount |

## Game Systems

| Symbol | Description |
|--------|-------------|
| `GameState` (interface) | Enter, Update, Draw, Exit for screen management |
| `StateMachine` (struct) | Stack-based game state manager |
| `NewStateMachine(initial GameState) *StateMachine` | Creates a state machine |
| `(*StateMachine).Init()` | Triggers `Enter()` on the initial state |
| `(*StateMachine).Update()` | Delegates to the top state's `Update()` |
| `(*StateMachine).Draw()` | Delegates to the top state's `Draw()` |
| `(*StateMachine).Switch(state GameState)` | Replaces the top state |
| `(*StateMachine).Push(state GameState)` | Overlays a new state (for pause menus, dialogs) |
| `(*StateMachine).Pop()` | Removes the top state |
| `(*StateMachine).Current() GameState` | Returns the active state |
| `(*StateMachine).Depth() int` | Number of states on the stack |
| `MenuItem` (struct) | Label, Disabled, OnConfirm callback |
| `Menu` (struct) | D-pad-navigated list with cursor tracking |
| `NewMenu(items []MenuItem) *Menu` | Creates a menu |
| `(*Menu).HandleInput() bool` | Reads D-pad/A button, returns true on confirm |
| `(*Menu).Draw()` | Renders the menu using `DrawText` |
| `(*Menu).MoveUp()` | Moves cursor up (skips disabled items) |
| `(*Menu).MoveDown()` | Moves cursor down (skips disabled items) |
| `(*Menu).Confirm()` | Triggers the selected item's callback |
| `(*Menu).Cursor() int` | Returns the cursor position |
| `(*Menu).SetCursor(index int)` | Sets the cursor position |
| `(*Menu).Selected() MenuItem` | Returns the highlighted item |
| `(*Menu).Count() int` | Number of items |
| `Timer` (struct) | Counts down a fixed number of frames |
| `NewTimer(durationFrames int) *Timer` | Creates a one-shot timer |
| `(*Timer).Tick() bool` | Advances one frame, returns true when it finishes |
| `(*Timer).Done() bool` | True when the timer has expired |
| `(*Timer).Progress() float32` | Elapsed/duration ratio (0-1) |
| `(*Timer).Elapsed() int` | Frames elapsed |
| `(*Timer).Remaining() int` | Frames left |
| `(*Timer).Duration() int` | Total frame count |
| `(*Timer).Reset()` | Restarts with the same duration |
| `(*Timer).ResetWith(durationFrames int)` | Restarts with a new duration |
| `RepeatingTimer` (struct) | Triggers at a fixed interval, counts triggers |
| `NewRepeatingTimer(intervalFrames int) *RepeatingTimer` | Creates a repeating timer |
| `(*RepeatingTimer).Tick() bool` | Advances one frame, returns true on trigger |
| `(*RepeatingTimer).Count() int` | Number of times triggered |
| `(*RepeatingTimer).Reset()` | Clears elapsed time and count |

## Parallax

| Symbol | Description |
|--------|-------------|
| `ParallaxLayer` (struct) | SpeedX, SpeedY scroll multipliers |
| `(ParallaxLayer).Offset(cameraX, cameraY int) (int, int)` | Computes the layer offset for a camera position |
| `ParallaxConfig` (struct) | Holds a slice of `ParallaxLayer` |
| `NewParallaxConfig(speeds ...ParallaxLayer) ParallaxConfig` | Creates a parallax config from layer speeds |
| `(ParallaxConfig).LayerOffset(layer, cameraX, cameraY int) (int, int)` | Returns the scroll offset for a specific layer |

## Transitions

| Symbol | Description |
|--------|-------------|
| `TransitionStyle` (int) | Enum for transition types |
| `FadeToBlack` | Screen fades to black |
| `FadeFromBlack` | Screen fades in from black |
| `Transition` (struct) | Style, Duration, active state |
| `StartTransition(style TransitionStyle, durationFrames int) *Transition` | Begins a screen transition |
| `(*Transition).Advance()` | Advances the transition by one frame |
| `(*Transition).Done() bool` | True when the transition has completed |
| `(*Transition).Active() bool` | True while the transition is running |
| `(*Transition).Stop()` | Cancels the transition |
| `(*Transition).Draw()` | Renders the transition overlay |

## Draw Regions

| Symbol | Description |
|--------|-------------|
| `DrawRegion` (struct) | X, Y, W, H defining a sub-rectangle of the screen |
| `(DrawRegion).Active() bool` | True if the region restricts drawing |
| `(DrawRegion).Offset(x, y int) (int, int)` | Translates local coordinates to screen space |
| `(DrawRegion).Clip(x1, y1, x2, y2 int) (int, int, int, int, bool)` | Clips a rectangle to the region, returns false if outside |
| `(DrawRegion).ContainsPoint(x, y int) bool` | True if a local point is inside the region |
| `SetDrawRegion(x, y, w, h int)` | Restricts drawing to a screen sub-rectangle (nestable) |
| `ResetDrawRegion()` | Removes the most recent draw region |

## Camera

| Symbol | Description |
|--------|-------------|
| `Camera` (struct) | X, Y, Width, Height, Zoom, FollowTarget, FollowSpeed, Bounds |
| `(*Camera).EffectiveZoom() float32` | Returns Zoom, defaulting to 1.0 if unset |
| `(*Camera).WorldToScreen(worldX, worldY float32) (float32, float32)` | Converts world coordinates to screen coordinates |
| `(*Camera).UpdateFollow()` | Moves the camera toward `FollowTarget` with smooth lerp |
| `(*Camera).ClampToBounds()` | Restricts camera position to stay within `Bounds` |
| `(*Camera).AddTrauma(amount float32)` | Adds screen shake intensity (0-1) |
| `(*Camera).UpdateShake()` | Decays trauma each frame |
| `(*Camera).ShakeOffset() (int, int)` | Returns the current frame's shake displacement |

## Sub-Packages

GoSprite64 also provides these sub-packages for specialized functionality:

| Package | Import Path | Purpose |
|---------|-------------|---------|
| `math2d` | `gosprite64/math2d` | 2D vectors, rectangles, collision, easing, grid utilities, random numbers |
| `math3d` | `gosprite64/math3d` | 3D vectors, 4x4 matrices, viewport projection |
| `save` | `gosprite64/save` | EEPROM, SRAM, and FlashRAM save data |
| `gfx` | `gosprite64/gfx` | Low-level display list construction and execution |
| `dma` | `gosprite64/dma` | DMA transfer helpers, MIO0 decompression, memory pool |
| `rspq` | `gosprite64/rspq` | RSP task queue, microcode loading, OS task submission |
| `n64os` | `gosprite64/n64os` | N64 OS primitives: scheduler, timers, events, messages |
| `scene3d` | `gosprite64/scene3d` | 3D scene graph, camera, mesh, LOD, triangle rendering |
| `audio/sfx` | `gosprite64/audio/sfx` | Sound effect ID type |
| `audio/music` | `gosprite64/audio/music` | Music track ID type |
| `audio/bank` | `gosprite64/audio/bank` | Instrument bank loading |
| `audio/sequence` | `gosprite64/audio/sequence` | Sequence player for MIDI-style playback |
