# Why GoSprite64?

![Gopher](../images/logo.png)

**GoSprite64** is your portal to building retro-fueled 2D games for the **Nintendo 64**, using the modern power of **Go**. With clean APIs, minimal setup (Linux, Mac, Windows are all supported), and a rebellious retro soul, it lets you bring your pixel dreams to life - on real N64 hardware.

## What's GoSprite64?

GoSprite64 is a Go library for making 2D games that run natively on the Nintendo 64. It wraps low-level N64 quirks in a modern API inspired by modern game engines, so you can focus on your game logic - not the hardware headaches.

Designed for developers who love Go and grew up on cartridges, GoSprite64 makes retro dev surprisingly fun and productive.

## Feature Highlights

GoSprite64 gives you a complete toolkit for N64 game development:

- **Graphics** - [Drawing primitives](../05-graphics/drawing-primitives.md), [sprites](../05-graphics/sprites.md) with flip/scale/rotate, [sprite sheets](../05-graphics/sprite-sheets.md), [animation player](../05-graphics/animation-player.md), [custom bitmap fonts](../05-graphics/custom-fonts.md) with [alignment](../05-graphics/text-alignment.md), [parallax scrolling](../05-graphics/parallax.md), [screen transitions](../05-graphics/transitions.md), and [draw regions](../05-graphics/draw-regions.md) for split-screen
- **Tile Scenes** - Full [Tile2D pipeline](../08-tile-scenes/pipeline-overview.md) for authoring and rendering tile-based worlds with [tile sheets and maps](../08-tile-scenes/tile-sheets-and-maps.md), [bundles](../08-tile-scenes/bundles-and-loading.md), and [camera scrolling](../08-tile-scenes/camera-and-scrolling.md)
- **Audio** - VADPCM-compressed [sound effects and music](../07-audio/sfx-and-music.md), [sequence player](../07-audio/sequence-player.md) for MIDI-like playback, and [instrument banks](../07-audio/instrument-banks.md)
- **Input** - [D-pad and buttons](../06-input/buttons-and-dpad.md), [analog stick](../06-input/analog-stick.md) with deadzone, [multi-controller](../06-input/multi-controller.md) support for up to 4 players, [rumble pak](../06-input/rumble.md) control, and [input recording/replay](../06-input/input-replay.md)
- **Game Systems** - [State machine](../09-game-systems/state-machine.md) with push/pop/switch, [timers](../09-game-systems/timers.md) (one-shot and repeating), [menus](../09-game-systems/menus.md) with D-pad navigation, and [save data](../09-game-systems/save-data.md) (EEPROM, SRAM, FlashRAM)
- **2D Math** - [Vectors](../10-math/vectors.md), [rectangles](../10-math/rectangles.md), [AABB collision detection](../10-math/collision-detection.md) with sweep and resolution, [easing functions](../10-math/easing-functions.md), [grid utilities](../10-math/grid-utilities.md), and [deterministic random numbers](../10-math/random-numbers.md)
- **3D Graphics** - [3D math](../11-3d-graphics/3d-math.md) (Mat4, Vec3, perspective/ortho projections), [scene graph](../11-3d-graphics/scene-graph.md), [display lists](../11-3d-graphics/display-lists.md), and [triangle rendering](../11-3d-graphics/triangle-rendering.md)

## Fixed Resolution

GoSprite64 exposes one official fixed resolution and drawing space: **288x216** logical pixels.

That is the public rendering contract for gameplay code. The runtime centers the canvas inside the framebuffer and presents it with square pixels, while public drawing APIs such as `FillRect`, `DrawRect`, `DrawLine`, and `DrawText` all operate in that same logical space.

If you build and run `examples/calibration`, you should see this reference frame:

![Calibration scene showing the fixed 288x216 logical canvas](../images/fixed-resolution-calibration.png)

For a deep dive into why pixels on the N64 are not always square and how GoSprite64 solves this, read the [Square Pixels](../04-core-concepts/square-pixels.md) chapter.

## Why Go?

Go is a clean, fast, pragmatic and efficient language. By using Go for Nintendo 64 development, GoSprite64 opens the door for cloud developers to create retro-style games with confidence and speed. The library bridges modern programming concepts with the raw power of a classic console.

## What's in This Book?

This book introduces you to GoSprite64, guiding you through everything from setup to building full 2D games.

You'll learn how to:

- Build and flash N64 ROMs
- Draw and move sprites
- Design and scroll tilemaps
- Handle input from controllers
- Play sound effects and music
- Build game screens with state machines and menus

Whether you're nostalgic for the era or just curious about console programming, this book aims to get you productive with GoSprite64 as fast as possible.

## Who Is This Book For?

This book is for:

- Developers interested in retro console programming
- Go programmers curious about low-level game development
- Hobbyists or indie devs looking to make something fun for the Nintendo 64

Some experience with Go is recommended. If you're brand new to game development or Go, consider starting with a simpler platform or tutorial first.

## Helpful Links

- [GoSprite64 GitHub](https://github.com/drpaneas/gosprite64) - main development repo
- [GoSprite64 Website](https://gosprite64.dev) - official docs and examples
- [GoSprite64 Discussions](https://github.com/drpaneas/gosprite64/discussions) - get help, share ideas, or show off your projects
- [clktmr/n64](https://github.com/clktmr/n64) - low-level Go SDK for N64
- [Embedded-Go](https://github.com/embeddedgo/go) - support for N64's architecture
- [Awesome N64 Dev](https://github.com/shiftclock/awesome-n64dev) - great collection of tools, docs, and inspiration
