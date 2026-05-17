# Feature Overview

A complete reference of every feature GoSprite64 provides, organized by section.

## Core (`gosprite64`)

| Feature | Description | Docs |
|---------|-------------|------|
| `Game` interface | `Init()`, `Update()`, `Draw()` lifecycle for your game | [Game Loop](../04-core-concepts/game-loop.md) |
| `Run(g Game)` | Starts the 60 FPS game loop with fixed timestep | [Game Loop](../04-core-concepts/game-loop.md) |
| 288x216 canvas | Fixed logical resolution for all drawing APIs | [Fixed Canvas](../04-core-concepts/fixed-canvas.md) |
| 16-color palette | Built-in named colors: `Black`, `White`, `Red`, etc. | [Colors](../04-core-concepts/colors.md) |

## Drawing Functions

| Feature | Description | Docs |
|---------|-------------|------|
| `ClearScreen()` | Fills the screen with black | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `ClearScreenWith(c)` | Fills the screen with any color | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `FillRect(x1,y1,x2,y2,c)` | Draws a filled rectangle | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `DrawRect(x1,y1,x2,y2,c)` | Draws a rectangle outline | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `DrawLine(x1,y1,x2,y2,c)` | Draws a 1px line (Bresenham for diagonals) | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `DrawText(str,x,y,c)` | Draws text using the built-in 8x8 font | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `DrawImage(src,x,y)` | Draws a Go `image.Image` at logical coordinates | [Drawing Primitives](../05-graphics/drawing-primitives.md) |
| `DrawWorldImage(src,x,y,cam)` | Draws an image offset by camera position | [Drawing Primitives](../05-graphics/drawing-primitives.md) |

## Sprites

| Feature | Description | Docs |
|---------|-------------|------|
| `LoadSpriteSheet(path)` | Loads a `.sheet` file from the cartridge filesystem | [Sprite Sheets](../05-graphics/sprite-sheets.md) |
| `SpriteSheet.FrameCount()` | Returns the number of frames in the sheet | [Sprite Sheets](../05-graphics/sprite-sheets.md) |
| `SpriteSheet.FrameWidth()` | Returns the pixel width of a single frame | [Sprite Sheets](../05-graphics/sprite-sheets.md) |
| `SpriteSheet.FrameHeight()` | Returns the pixel height of a single frame | [Sprite Sheets](../05-graphics/sprite-sheets.md) |
| `DrawSprite(sheet,frame,x,y)` | Draws a sprite frame at logical coordinates | [Sprites](../05-graphics/sprites.md) |
| `DrawSpriteWithOptions(...)` | Draws with flip, scale, rotation, blend, and alpha | [Sprites](../05-graphics/sprites.md) |
| `DrawWorldSprite(...)` | Draws a sprite offset by camera position | [Sprites](../05-graphics/sprites.md) |
| `DrawWorldSpriteWithOptions(...)` | World-space sprite with full draw options | [Sprites](../05-graphics/sprites.md) |
| `DrawSpriteOptions` | Struct: `FlipH`, `FlipV`, `ScaleX/Y`, `Rotation`, `Origin`, `Blend`, `Alpha` | [Sprites](../05-graphics/sprites.md) |
| `BlendNone`, `BlendMasked`, `BlendAlpha` | Blend mode constants for sprite drawing | [Sprites](../05-graphics/sprites.md) |

## Animation

| Feature | Description | Docs |
|---------|-------------|------|
| `AnimationSet` | Collection of named animation clips loaded from `.anim` files | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationClip` | A single animation: name, FPS, and frame indices | [Animation Player](../05-graphics/animation-player.md) |
| `NewAnimationPlayer()` | Creates a player that drives sprite frame changes | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationPlayer.Play(clip)` | Starts playing a clip from the beginning | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationPlayer.Advance(ticks)` | Advances the animation by N ticks (call each frame with 1) | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationPlayer.Frame()` | Returns the current sprite sheet frame index | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationPlayer.SetLoop(bool)` | Enables or disables looping | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationPlayer.Pause/Resume/Stop/Restart` | Playback control | [Animation Player](../05-graphics/animation-player.md) |
| `AnimationPlayer.Playing()` / `Done()` | Status queries | [Animation Player](../05-graphics/animation-player.md) |

## Custom Fonts

| Feature | Description | Docs |
|---------|-------------|------|
| `NewFont(sheet, glyphs, lineHeight)` | Creates a font from a sprite sheet and glyph map | [Custom Fonts](../05-graphics/custom-fonts.md) |
| `Font.DrawTextEx(text, x, y, align)` | Draws text with left/center/right alignment | [Text Alignment](../05-graphics/text-alignment.md) |
| `Font.MeasureText(text)` | Returns pixel width and height of rendered text | [Custom Fonts](../05-graphics/custom-fonts.md) |
| `Font.WrapText(text, maxWidth)` | Inserts newlines to fit text within a pixel width | [Text Alignment](../05-graphics/text-alignment.md) |
| `FormatScore(score, width)` | Formats an integer with leading zeros | [Custom Fonts](../05-graphics/custom-fonts.md) |
| `AlignLeft`, `AlignCenter`, `AlignRight` | Text alignment constants | [Text Alignment](../05-graphics/text-alignment.md) |

## Parallax Scrolling

| Feature | Description | Docs |
|---------|-------------|------|
| `NewParallaxConfig(speeds...)` | Configures multi-layer parallax with speed factors | [Parallax Scrolling](../05-graphics/parallax.md) |
| `ParallaxConfig.LayerOffset(layer, camX, camY)` | Returns the scroll offset for a given layer and camera | [Parallax Scrolling](../05-graphics/parallax.md) |
| `ParallaxLayer` | Defines `SpeedX` and `SpeedY` multipliers (0.0 = static, 1.0 = full speed) | [Parallax Scrolling](../05-graphics/parallax.md) |

## Screen Transitions

| Feature | Description | Docs |
|---------|-------------|------|
| `StartTransition(style, frames)` | Begins a fade transition over N frames | [Transitions](../05-graphics/transitions.md) |
| `FadeToBlack`, `FadeFromBlack` | Transition style constants | [Transitions](../05-graphics/transitions.md) |
| `Transition.Advance()` | Steps the transition forward one frame | [Transitions](../05-graphics/transitions.md) |
| `Transition.Draw()` | Renders the transition overlay | [Transitions](../05-graphics/transitions.md) |
| `Transition.Done()` / `Active()` / `Stop()` | Status and control | [Transitions](../05-graphics/transitions.md) |

## Draw Regions

| Feature | Description | Docs |
|---------|-------------|------|
| `SetDrawRegion(x, y, w, h)` | Restricts drawing to a sub-rectangle (for split-screen) | [Draw Regions](../05-graphics/draw-regions.md) |
| `ResetDrawRegion()` | Pops the most recent draw region | [Draw Regions](../05-graphics/draw-regions.md) |
| `DrawRegion.Clip(...)` | Offsets and clips coordinates to region bounds | [Draw Regions](../05-graphics/draw-regions.md) |
| `DrawRegion.ContainsPoint(x, y)` | Hit-tests a local coordinate against the region | [Draw Regions](../05-graphics/draw-regions.md) |

## Input

| Feature | Description | Docs |
|---------|-------------|------|
| `IsButtonDown(button)` | Returns true while a button is held (port 0) | [D-Pad and Buttons](../06-input/buttons-and-dpad.md) |
| `IsButtonJustPressed(button)` | Returns true on the frame a button is first pressed (port 0) | [D-Pad and Buttons](../06-input/buttons-and-dpad.md) |
| `StickPosition(deadzone)` | Returns analog stick X/Y in [-1.0, 1.0] (port 0) | [Analog Stick](../06-input/analog-stick.md) |
| `PlayerButtonDown(port, button)` | Per-port button check for multiplayer | [Multi-Controller](../06-input/multi-controller.md) |
| `PlayerButtonJustPressed(port, button)` | Per-port just-pressed check | [Multi-Controller](../06-input/multi-controller.md) |
| `PlayerStickPosition(port, deadzone)` | Per-port analog stick | [Multi-Controller](../06-input/multi-controller.md) |
| `IsControllerConnected(port)` | Checks if a controller is plugged in | [Multi-Controller](../06-input/multi-controller.md) |
| `ConnectedControllers()` | Returns the number of connected controllers | [Multi-Controller](../06-input/multi-controller.md) |
| `SetRumble(port, enabled)` | Turns the Rumble Pak on or off | [Rumble](../06-input/rumble.md) |
| Button constants | `ButtonA`, `ButtonB`, `ButtonZ`, `ButtonStart`, `ButtonDPadUp/Down/Left/Right`, `ButtonL`, `ButtonR`, `ButtonCUp/Down/Left/Right` | [D-Pad and Buttons](../06-input/buttons-and-dpad.md) |

## Input Recording and Replay

| Feature | Description | Docs |
|---------|-------------|------|
| `NewInputRecorder(playerCount)` | Creates a recorder that captures per-frame controller state | [Input Replay](../06-input/input-replay.md) |
| `InputRecorder.CaptureFrame(player, input)` | Records one frame of input | [Input Replay](../06-input/input-replay.md) |
| `InputRecorder.Finish()` | Finalizes recording into `ReplayData` | [Input Replay](../06-input/input-replay.md) |
| `NewInputPlayer(data)` | Creates a player that replays recorded input | [Input Replay](../06-input/input-replay.md) |
| `InputPlayer.NextFrame(player)` | Returns the next frame of recorded input | [Input Replay](../06-input/input-replay.md) |
| `InputPlayer.Done()` / `Reset()` | Playback status and restart | [Input Replay](../06-input/input-replay.md) |

## Audio

| Feature | Description | Docs |
|---------|-------------|------|
| `RegisterAudioBundle(bundle)` | Registers VADPCM audio assets before the game loop starts | [SFX and Music](../07-audio/sfx-and-music.md) |
| `PlaySoundEffect(id)` | Triggers a one-shot sound effect | [SFX and Music](../07-audio/sfx-and-music.md) |
| `PlayMusic(id)` | Starts background music playback | [SFX and Music](../07-audio/sfx-and-music.md) |
| `StopMusic()` | Stops the current music track | [SFX and Music](../07-audio/sfx-and-music.md) |
| `SetSoundEffectVolume(v)` | Sets SFX volume (0.0 to 1.0) | [SFX and Music](../07-audio/sfx-and-music.md) |
| `SetMusicVolume(v)` | Sets music volume (0.0 to 1.0) | [SFX and Music](../07-audio/sfx-and-music.md) |
| `sequence.NewPlayer()` | Creates a MIDI-like sequence player | [Sequence Player](../07-audio/sequence-player.md) |
| `sequence.Player.Play/Stop/Pause/Resume` | Sequence playback control | [Sequence Player](../07-audio/sequence-player.md) |
| `sequence.Player.SetTempo(bpm)` | Sets playback tempo | [Sequence Player](../07-audio/sequence-player.md) |
| `sequence.Player.SetLoop(start, count)` | Configures loop points | [Sequence Player](../07-audio/sequence-player.md) |

## Tile Scene Pipeline

| Feature | Description | Docs |
|---------|-------------|------|
| `OpenBundle(path)` | Opens a `.bundle` file containing sheets, maps, and animations | [Bundles and Loading](../08-tile-scenes/bundles-and-loading.md) |
| `LoadScene(bundle)` | Loads all assets from a bundle into a renderable scene | [Pipeline Overview](../08-tile-scenes/pipeline-overview.md) |
| `Scene.Draw(cam)` | Renders visible tiles to the screen through the camera | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |
| `Scene.Map()` | Returns the scene's `Map` for tile queries | [Tile Sheets and Maps](../08-tile-scenes/tile-sheets-and-maps.md) |
| `Map.Width()` / `Height()` | Map dimensions in tiles | [Tile Sheets and Maps](../08-tile-scenes/tile-sheets-and-maps.md) |
| `Map.TileWidth()` / `TileHeight()` | Tile dimensions in pixels | [Tile Sheets and Maps](../08-tile-scenes/tile-sheets-and-maps.md) |
| `Map.PixelWidth()` / `PixelHeight()` | Total map size in pixels | [Tile Sheets and Maps](../08-tile-scenes/tile-sheets-and-maps.md) |
| `Map.TileAt(layer, col, row)` | Returns the tile ID at a grid cell | [Tile Sheets and Maps](../08-tile-scenes/tile-sheets-and-maps.md) |
| `Scene.Stats()` | Returns `RuntimeStats` with visible tile count and upload count | [Pipeline Overview](../08-tile-scenes/pipeline-overview.md) |

## Camera

| Feature | Description | Docs |
|---------|-------------|------|
| `Camera` struct | Position, size, zoom, follow target, bounds, and screen shake | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |
| `Camera.WorldToScreen(x, y)` | Converts world coordinates to screen space | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |
| `Camera.UpdateFollow()` | Smoothly moves camera toward the follow target | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |
| `Camera.ClampToBounds()` | Prevents the camera from leaving the world bounds | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |
| `Camera.AddTrauma(amount)` | Adds screen shake intensity (0 to 1) | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |
| `Camera.ShakeOffset()` | Returns the current shake displacement for drawing | [Camera and Scrolling](../08-tile-scenes/camera-and-scrolling.md) |

## Game Systems

### State Machine

| Feature | Description | Docs |
|---------|-------------|------|
| `GameState` interface | `Enter()`, `Update()`, `Draw()`, `Exit()` for each screen | [State Machine](../09-game-systems/state-machine.md) |
| `NewStateMachine(initial)` | Creates a state machine with an initial state | [State Machine](../09-game-systems/state-machine.md) |
| `StateMachine.Switch(state)` | Replaces the top state (calls Exit then Enter) | [State Machine](../09-game-systems/state-machine.md) |
| `StateMachine.Push(state)` | Overlays a new state (for pause menus, dialogs) | [State Machine](../09-game-systems/state-machine.md) |
| `StateMachine.Pop()` | Removes the top state and returns to the one below | [State Machine](../09-game-systems/state-machine.md) |
| `StateMachine.Update()` / `Draw()` | Delegates to the top state | [State Machine](../09-game-systems/state-machine.md) |

### Timers

| Feature | Description | Docs |
|---------|-------------|------|
| `NewTimer(frames)` | Creates a countdown timer | [Timers](../09-game-systems/timers.md) |
| `Timer.Tick()` | Advances by one frame; returns true on the finishing frame | [Timers](../09-game-systems/timers.md) |
| `Timer.Done()` / `Progress()` / `Remaining()` | Status queries | [Timers](../09-game-systems/timers.md) |
| `Timer.Reset()` / `ResetWith(frames)` | Restart the timer | [Timers](../09-game-systems/timers.md) |
| `NewRepeatingTimer(interval)` | Creates a timer that fires every N frames | [Timers](../09-game-systems/timers.md) |
| `RepeatingTimer.Tick()` | Returns true on trigger frames | [Timers](../09-game-systems/timers.md) |
| `RepeatingTimer.Count()` | Returns how many times it has triggered | [Timers](../09-game-systems/timers.md) |

### Menus

| Feature | Description | Docs |
|---------|-------------|------|
| `NewMenu(items)` | Creates a D-pad-navigated menu from `MenuItem` entries | [Menus](../09-game-systems/menus.md) |
| `Menu.HandleInput()` | Reads the controller and moves the cursor; returns true on confirm | [Menus](../09-game-systems/menus.md) |
| `Menu.Draw()` | Renders the menu with cursor indicator | [Menus](../09-game-systems/menus.md) |
| `Menu.MoveUp()` / `MoveDown()` | Manual cursor movement (skips disabled items) | [Menus](../09-game-systems/menus.md) |
| `MenuItem` | Struct: `Label`, `Disabled`, `OnConfirm` callback | [Menus](../09-game-systems/menus.md) |

### Save Data

| Feature | Description | Docs |
|---------|-------------|------|
| `save.Storage` interface | Uniform API for EEPROM, SRAM, and FlashRAM | [Save Data](../09-game-systems/save-data.md) |
| `save.ReadAll(s)` / `WriteAll(s, data)` | Read or write the entire save storage | [Save Data](../09-game-systems/save-data.md) |
| `save.Checksum(data)` | Additive checksum for save integrity | [Save Data](../09-game-systems/save-data.md) |
| Storage types | `StorageEEPROM4K` (512B), `StorageEEPROM16K` (2KB), `StorageSRAM` (32KB), `StorageFlashRAM` (128KB) | [Save Data](../09-game-systems/save-data.md) |

## 2D Math (`math2d`)

### Vectors

| Feature | Description | Docs |
|---------|-------------|------|
| `Vec2` | 2D vector with `X`, `Y` float32 fields | [Vectors](../10-math/vectors.md) |
| `Add`, `Sub`, `Scale`, `Negate`, `Abs` | Arithmetic operations | [Vectors](../10-math/vectors.md) |
| `Length`, `LengthSq`, `Normalize` | Magnitude and normalization | [Vectors](../10-math/vectors.md) |
| `Dot`, `Distance`, `DistanceSq` | Products and distances | [Vectors](../10-math/vectors.md) |
| `Lerp`, `Rotate`, `Angle` | Interpolation and rotation | [Vectors](../10-math/vectors.md) |
| `Min`, `Max` | Component-wise min/max | [Vectors](../10-math/vectors.md) |

### Rectangles

| Feature | Description | Docs |
|---------|-------------|------|
| `Rect` | Axis-aligned rectangle: `X`, `Y`, `W`, `H` | [Rectangles](../10-math/rectangles.md) |
| `RectFromCenter(center, w, h)` | Creates a rect centered on a point | [Rectangles](../10-math/rectangles.md) |
| `ContainsPoint`, `ContainsRect`, `Overlaps` | Spatial queries | [Rectangles](../10-math/rectangles.md) |
| `Intersection`, `Expand`, `Center` | Rect operations | [Rectangles](../10-math/rectangles.md) |

### Collision Detection

| Feature | Description | Docs |
|---------|-------------|------|
| `AABBOverlap(a, b)` | Returns true if two rects overlap | [Collision Detection](../10-math/collision-detection.md) |
| `AABBPenetration(a, b)` | Returns the minimum penetration vector | [Collision Detection](../10-math/collision-detection.md) |
| `AABBResolve(a, b)` | Pushes rect `a` out of rect `b` | [Collision Detection](../10-math/collision-detection.md) |
| `AABBSweep(a, vel, b)` | Swept AABB test: returns hit time and normal | [Collision Detection](../10-math/collision-detection.md) |
| `Layer`, `Collider`, `ColliderOverlap` | Layer-masked collision filtering | [Collision Detection](../10-math/collision-detection.md) |

### Easing Functions

| Feature | Description | Docs |
|---------|-------------|------|
| `Clamp(v, lo, hi)` | Restricts a value to a range | [Easing Functions](../10-math/easing-functions.md) |
| `Lerp(a, b, t)` | Linear interpolation | [Easing Functions](../10-math/easing-functions.md) |
| `InvLerp(a, b, v)` | Inverse lerp (returns 0-1 ratio) | [Easing Functions](../10-math/easing-functions.md) |
| `Remap(v, inMin, inMax, outMin, outMax)` | Maps a value from one range to another | [Easing Functions](../10-math/easing-functions.md) |
| `MoveToward(current, target, maxDelta)` | Moves toward target by at most maxDelta | [Easing Functions](../10-math/easing-functions.md) |
| `EaseInQuad`, `EaseOutQuad`, `EaseInOutQuad` | Quadratic easing | [Easing Functions](../10-math/easing-functions.md) |
| `EaseInCubic`, `EaseOutCubic`, `EaseInOutCubic` | Cubic easing | [Easing Functions](../10-math/easing-functions.md) |
| `SmoothStep(edge0, edge1, x)` | Hermite interpolation | [Easing Functions](../10-math/easing-functions.md) |

### Grid Utilities

| Feature | Description | Docs |
|---------|-------------|------|
| `NewGrid[T](cols, rows)` | Creates a generic 2D grid | [Grid Utilities](../10-math/grid-utilities.md) |
| `Get`, `Set`, `Clear`, `Fill` | Cell access and bulk operations | [Grid Utilities](../10-math/grid-utilities.md) |
| `ScanRow`, `ScanCol` | Run-length scanning for matching groups | [Grid Utilities](../10-math/grid-utilities.md) |
| `CountValue`, `FindAll`, `Neighbors4` | Queries and spatial helpers | [Grid Utilities](../10-math/grid-utilities.md) |

### Random Numbers

| Feature | Description | Docs |
|---------|-------------|------|
| `NewRand(seed)` | Creates a deterministic xoshiro128** PRNG | [Random Numbers](../10-math/random-numbers.md) |
| `Uint32`, `Intn(n)`, `Float32` | Raw random values | [Random Numbers](../10-math/random-numbers.md) |
| `RangeInt(min, max)`, `RangeFloat32(min, max)` | Range-bounded random values | [Random Numbers](../10-math/random-numbers.md) |
| `Bool()` | Random boolean | [Random Numbers](../10-math/random-numbers.md) |

## 3D Graphics

| Feature | Description | Docs |
|---------|-------------|------|
| `math3d.Mat4` | 4x4 matrix with multiply, perspective, ortho, lookAt, translate, rotate, scale | [3D Math](../11-3d-graphics/3d-math.md) |
| `math3d.Vec3` / `Vec4` | 3D and 4D vectors with arithmetic, dot, cross, normalize | [3D Math](../11-3d-graphics/3d-math.md) |
| `math3d.Viewport` | Maps clip-space to screen coordinates | [3D Math](../11-3d-graphics/3d-math.md) |
| `scene3d.NewScene()` | Creates a 3D scene graph | [Scene Graph](../11-3d-graphics/scene-graph.md) |
| `scene3d.Node` | Scene graph node with transform, children, and render function | [Scene Graph](../11-3d-graphics/scene-graph.md) |
| `scene3d.NewMeshNode(name, dl)` | Creates a node that renders a display list | [Scene Graph](../11-3d-graphics/scene-graph.md) |
| `scene3d.NewPerspectiveCamera(...)` | Creates a perspective camera node | [Scene Graph](../11-3d-graphics/scene-graph.md) |
| `scene3d.NewOrthoCamera(...)` | Creates an orthographic camera node | [Scene Graph](../11-3d-graphics/scene-graph.md) |
| `scene3d.DrawScene(scene)` | Traverses the scene graph and renders | [Triangle Rendering](../11-3d-graphics/triangle-rendering.md) |
| `gfx.DisplayList` | GPU command buffer for triangle rendering | [Display Lists](../11-3d-graphics/display-lists.md) |

## Low-Level

| Feature | Description | Docs |
|---------|-------------|------|
| `dma.CartToRDRAM(offset, dst)` | DMA transfer from cartridge ROM to RDRAM | [DMA Transfers](../12-low-level/dma-transfers.md) |
| `dma.SRAMRead` / `SRAMWrite` | Direct SRAM access via DMA | [DMA Transfers](../12-low-level/dma-transfers.md) |
| `dma.NewPool(base, size)` | Memory pool with head/tail allocation | [Memory Pools](../12-low-level/memory-pools.md) |
| `dma.NewSegmentTable()` | Segment address translation table | [DMA Transfers](../12-low-level/dma-transfers.md) |
| `rspq.NewQueue()` | RSP task queue for submitting microcode tasks | [RSP Task Queue](../12-low-level/rsp-task-queue.md) |
| `rspq.Load(microcode)` | Loads RSP microcode | [RSP Task Queue](../12-low-level/rsp-task-queue.md) |
| `n64os.NewMessageQueue(size)` | OS-level message queue for event handling | [N64 OS Primitives](../12-low-level/n64-os-primitives.md) |
| `n64os.NewScheduler(events, fn)` | Task scheduler for graphics and audio | [N64 OS Primitives](../12-low-level/n64-os-primitives.md) |
| `n64os.NewEventRouter()` | Routes hardware events to message queues | [N64 OS Primitives](../12-low-level/n64-os-primitives.md) |
| `n64os.NewTimer(queue, msg, interval)` | OS-level timer with message delivery | [N64 OS Primitives](../12-low-level/n64-os-primitives.md) |
