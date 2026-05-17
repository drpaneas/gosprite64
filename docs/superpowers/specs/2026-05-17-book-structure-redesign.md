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
- Installation (from existing getting_started.md - keeps everything except the Cursor/VS Code section)
- Hello World (from existing hello_world.md)
- Editor Setup (extracted from getting_started.md - the "Cursor / VS Code" section starting at the `## Cursor / VS Code` heading through the end of the `.vscode/settings.json` instructions)

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
- Vectors (from existing math2d.md, the "Vec2" section)
- Rectangles (from existing math2d.md, the "Rect" section)
- Collision Detection (new - covers math2d.AABBOverlap, AABBPenetration, AABBResolve, AABBSweep, Collider, Layer)
- Easing Functions (from existing math2d.md, the "Easing and interpolation" section)
- Grid Utilities (new - covers math2d.Grid, Run, GridCell)
- Random Numbers (from existing math2d.md, the "Rand" section)

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

mdbook supports subdirectories. The `book.toml` already has `src = "docs"` which means mdbook reads from the `docs/` directory. Subdirectories inside `docs/` are fully supported and require no config changes.

Each part gets its own directory:

```
docs/
  SUMMARY.md
  images/                         <-- shared image directory
    logo.png                      (moved from docs/logo.png)
    fixed-resolution-calibration.png  (moved from docs/)
    par-comparison.png            (moved from docs/)
    canvas-layout.png             (moved from docs/)
    rendering-pipeline.png        (moved from docs/)
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
  superpowers/                    <-- excluded from book, kept as-is
    specs/
    plans/
```

## Image Migration

All images currently in `docs/` are moved to a shared `docs/images/` directory. Image references in moved pages are updated to use relative paths from their new location:

| Image | Current location | New location |
|-------|-----------------|--------------|
| `logo.png` | `docs/logo.png` | `docs/images/logo.png` |
| `fixed-resolution-calibration.png` | `docs/` | `docs/images/` |
| `par-comparison.png` | `docs/` | `docs/images/` |
| `canvas-layout.png` | `docs/` | `docs/images/` |
| `rendering-pipeline.png` | `docs/` | `docs/images/` |

Pages that reference these images must update their paths. For example, `square_pixels.md` moves to `04-core-concepts/square-pixels.md`, so its reference `![...](par-comparison.png)` becomes `![...](../images/par-comparison.png)`.

The `introduction.md` reference `![Gopher](../logo.png)` becomes `![Gopher](../images/logo.png)` after the page moves to `01-welcome/`.

## Content Extraction Rules

### getting_started.md split

- **installation.md** gets: everything from `# Getting Started` through the end of step 5 ("Build all examples"), the `n64.env` explanation, the resolution/canvas paragraph, the calibration screenshot, the Windows section, and the Linux Fallback section.
- **editor-setup.md** gets: the `## Cursor / VS Code` section including all `.vscode/settings.json` instructions, the `go.alternateTools` note, and the "restart language server" instruction.

### math2d.md split

| New page | Content from math2d.md |
|----------|----------------------|
| `vectors.md` | The `## Vec2` section (lines 11-96) |
| `rectangles.md` | The `## Rect` section (lines 98-157) |
| `collision-detection.md` | New stub - covers `math2d.AABBOverlap`, `AABBPenetration`, `AABBResolve`, `AABBSweep`, `Collider`, `Layer` |
| `easing-functions.md` | The `## Easing and interpolation` section (lines 224-313) |
| `grid-utilities.md` | New stub - covers `math2d.Grid[T]`, `Run`, `GridCell`, `ScanRow`, `ScanCol` |
| `random-numbers.md` | The `## Rand` section (lines 159-222) |

Each split page keeps the package import instruction and the introductory sentence from math2d.md.

### first_tile_game.md migration

The full content moves to `03-tutorial/02-draw-a-tilemap.md`. References to `examples/simplegame` and `examples/tilemap` stay intact since those paths are repository-relative, not doc-relative.

## Excluded Directories

The `docs/superpowers/` directory (containing specs and plans) is not part of the mdbook build. It is not listed in SUMMARY.md and mdbook ignores unlisted files. This directory stays in place and is not moved.

## Stub Page Format

Each stub page contains:

```markdown
# [Feature Name]

[One sentence describing what this feature does and why you would use it.]

> This page is under construction. The feature is available in the codebase but documentation has not been written yet.
```

The blockquote format renders as a visible callout in mdbook without looking unprofessional.

## Verification

The spec is complete when:

- `mdbook build` succeeds with the new SUMMARY.md and produces no warnings about missing files
- All existing content is accessible under its new path with no broken image references
- Every exported feature in the codebase has at least a stub page in the book
- The book renders with a clear multi-level navigation sidebar

## Non-Goals

This sub-project does not include:

- Writing full content for stub pages (separate sub-project)
- Creating example programs (separate sub-project)
- Integrating web-based N64 emulator demos (separate sub-project)
- Writing the multi-step platformer tutorial content (separate sub-project)

## Rollout

1. Create all subdirectories under `docs/`
2. Create `docs/images/` and move all image files there
3. Copy existing doc files to their new paths, updating image references
4. Create stub pages for all undocumented features
5. Split `getting_started.md` into `installation.md` and `editor-setup.md`
6. Split `math2d.md` into the 4 content pages plus 2 stub pages
7. Write the new `SUMMARY.md` pointing at all new paths
8. Verify `mdbook build` succeeds with no warnings
9. Remove old top-level doc files that have been moved
10. Verify `mdbook build` still succeeds after cleanup
11. Commit
