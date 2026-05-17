# Book Structure Redesign

## Goal

Reorganize the GoSprite64 documentation from a flat list of pages into a 13-part structured book with sections for every feature in the codebase, so that no capability is hidden and new developers have a clear learning path.

## Context

The current documentation is a flat list of 14 pages in `docs/SUMMARY.md`. The codebase contains at least 50 exported functions and 10 sub-packages, many with zero documentation. Features like save data (EEPROM/SRAM/FlashRAM), 2D collision detection, custom fonts, parallax scrolling, screen transitions, multi-controller support, rumble, sequence audio, 3D scene graph, display lists, DMA transfers, and RSP task queues have no documentation at all.

The Excalibur.js framework demonstrates best-in-class game engine documentation with numbered sections, per-feature reference pages, and multi-step tutorials. This redesign adopts that structure for GoSprite64.

## Decision

Restructure the book into 13 numbered parts, move existing pages to their new locations, and create stub pages for every undocumented feature. This is a structural reorganization only - full content writing, example programs, and web demo integration are separate follow-up sub-projects.

## Structure

The new SUMMARY.md will have this hierarchy:

### Part 1: Welcome
- Why GoSprite64 (new)
- Feature Overview (new)

### Part 2: Getting Started
- Installation (from existing getting_started.md)
- Hello World (from existing hello_world.md)
- Editor Setup (extracted from getting_started.md)

### Part 3: Tutorial - Build a Platformer
- Step 1: Start the Engine (stub)
- Step 2: Draw a Tilemap World (seed from existing first_tile_game.md)
- Step 3: Add a Player Sprite (stub)
- Step 4: Animate the Player (stub)
- Step 5: Move with the D-Pad (stub)
- Step 6: Camera Following (stub)
- Step 7: Add Sound Effects (stub)
- Step 8: Add a Title Screen (stub)
- Step 9: Screen Transitions (stub)
- Step 10: Score Display (stub)
- Step 11: Final Polish (stub)

### Part 4: Core Concepts
- The Game Loop (new)
- The Fixed Canvas (new)
- Square Pixels (from existing square_pixels.md)
- Colors (new)

### Part 5: Graphics
- Drawing Primitives (new)
- Sprites (from existing sprites.md)
- Sprite Sheets (new)
- Animation Player (new)
- Custom Fonts (new)
- Text Alignment (new)
- Parallax Scrolling (new)
- Screen Transitions (new)
- Draw Regions (from existing draw_regions.md)

### Part 6: Input
- D-Pad and Buttons (new)
- Analog Stick (new)
- Multi-Controller Support (new)
- Rumble (new)
- Input Recording and Replay (from existing input_replay.md)

### Part 7: Audio
- Sound Effects and Music (from existing audio.md)
- Sequence Player (new)
- Instrument Banks (new)

### Part 8: Tile Scenes
- Tile2D Pipeline Overview (from existing tile2d.md)
- Tile Sheets and Maps (new)
- Bundles and Loading (new)
- Camera and Scrolling (new)

### Part 9: Game Systems
- State Machine (from existing state_machine.md)
- Timers (from existing timers.md)
- Menus (from existing menus.md)
- Save Data (new)

### Part 10: 2D Math
- Vectors (from existing math2d.md, split)
- Rectangles (new)
- Collision Detection (new)
- Easing Functions (new)
- Grid Utilities (new)
- Random Numbers (new)

### Part 11: 3D Graphics
- 3D Math (new)
- Scene Graph (new)
- Display Lists (new)
- Triangle Rendering (new)

### Part 12: Low-Level
- DMA Transfers (new)
- RSP Task Queue (new)
- N64 OS Primitives (new)
- Memory Pools (new)

### Part 13: Reference
- API Quick Reference (new)
- Performance Notes (new)
- Troubleshooting (new)

## File Layout

mdbook supports subdirectories. Each part gets its own directory:

```
docs/
  SUMMARY.md
  01-welcome/
    why-gosprite64.md
    feature-overview.md
  02-getting-started/
    installation.md
    hello-world.md
    editor-setup.md
  03-tutorial/
    01-start-the-engine.md
    02-draw-a-tilemap.md
    03-add-a-player-sprite.md
    04-animate-the-player.md
    05-move-with-dpad.md
    06-camera-following.md
    07-add-sound-effects.md
    08-add-title-screen.md
    09-screen-transitions.md
    10-score-display.md
    11-final-polish.md
  04-core-concepts/
    game-loop.md
    fixed-canvas.md
    square-pixels.md
    colors.md
  05-graphics/
    drawing-primitives.md
    sprites.md
    sprite-sheets.md
    animation-player.md
    custom-fonts.md
    text-alignment.md
    parallax.md
    transitions.md
    draw-regions.md
  06-input/
    buttons-and-dpad.md
    analog-stick.md
    multi-controller.md
    rumble.md
    input-replay.md
  07-audio/
    sfx-and-music.md
    sequence-player.md
    instrument-banks.md
  08-tile-scenes/
    pipeline-overview.md
    tile-sheets-and-maps.md
    bundles-and-loading.md
    camera-and-scrolling.md
  09-game-systems/
    state-machine.md
    timers.md
    menus.md
    save-data.md
  10-math/
    vectors.md
    rectangles.md
    collision-detection.md
    easing-functions.md
    grid-utilities.md
    random-numbers.md
  11-3d-graphics/
    3d-math.md
    scene-graph.md
    display-lists.md
    triangle-rendering.md
  12-low-level/
    dma-transfers.md
    rsp-task-queue.md
    n64-os-primitives.md
    memory-pools.md
  13-reference/
    api-quick-reference.md
    performance-notes.md
    troubleshooting.md
```

## Migration Rules

- Existing pages keep their content intact and are moved to new paths
- Each stub page contains exactly: a title, a one-sentence description of the feature, and "This page is under construction."
- The tutorial section (Part 3) gets stubs with step titles but no tutorial content yet
- `first_tile_game.md` content moves into `03-tutorial/02-draw-a-tilemap.md` as the seed for the multi-step tutorial
- Old top-level doc files are removed after their content is moved

## Verification

The spec is complete when:

- `mdbook build` succeeds with the new SUMMARY.md
- All existing content is accessible under its new path
- Every exported feature in the codebase has at least a stub page in the book
- No broken internal links
- The book renders with a clear multi-level navigation sidebar

## Non-Goals

This sub-project does not include:

- Writing full content for stub pages (separate sub-project)
- Creating example programs (separate sub-project)
- Integrating web-based N64 emulator demos (separate sub-project)
- Writing the multi-step platformer tutorial content (separate sub-project)

## Rollout

1. Create the directory structure under `docs/`
2. Write the new `SUMMARY.md` with all parts and pages
3. Move existing doc content to new paths
4. Create stub pages for all undocumented features
5. Remove old top-level doc files
6. Verify `mdbook build` succeeds
7. Commit
