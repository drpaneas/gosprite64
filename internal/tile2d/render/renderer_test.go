package render

import (
	"image"
	"testing"

	"github.com/clktmr/n64/rcp/texture"
	"github.com/drpaneas/gosprite64/internal/tile2d/visibility"
)

type testExecutionBridge struct {
	runFunc  func(x, y, tileWidth, tileHeight int, run PreparedRun)
	tileFunc func(x, y, width, height int, tile PreparedTile)
}

func (t testExecutionBridge) DrawPreparedRun(x, y, tileWidth, tileHeight int, run PreparedRun) {
	if t.runFunc != nil {
		t.runFunc(x, y, tileWidth, tileHeight, run)
	}
}

func (t testExecutionBridge) DrawPreparedTile(x, y, width, height int, tile PreparedTile) {
	if t.tileFunc != nil {
		t.tileFunc(x, y, width, height, tile)
	}
}

func TestRendererSkipsOffscreenTiles(t *testing.T) {
	r := NewRenderer(RenderHooks{})
	calls := r.DrawTileLayer(testScene(), testCamera())
	if calls.VisibleTiles != 0 {
		t.Fatalf("VisibleTiles = %d, want 0", calls.VisibleTiles)
	}
}

func TestRendererBatchesAdjacentTilesBySheet(t *testing.T) {
	r := NewRenderer(RenderHooks{})
	got := r.DrawTileLayer(adjacentSameSheetScene(), visibility.Camera{Width: 32, Height: 32})
	if got.Uploads > 1 {
		t.Fatalf("Uploads = %d, want <= 1", got.Uploads)
	}
}

func TestRendererInvokesDrawHookForVisibleTiles(t *testing.T) {
	var calls int
	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			tileFunc: func(x, y, width, height int, tile PreparedTile) {
				calls++
			},
		},
	})

	got := r.DrawTileLayer(adjacentSameSheetScene(), visibility.Camera{Width: 32, Height: 32})
	if got.VisibleTiles != 2 {
		t.Fatalf("VisibleTiles = %d, want 2", got.VisibleTiles)
	}
	if calls != 2 {
		t.Fatalf("DrawTile calls = %d, want 2", calls)
	}
}

func TestRendererUsesPreparedTileEntries(t *testing.T) {
	var got PreparedTile
	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			tileFunc: func(x, y, width, height int, tile PreparedTile) {
				got = tile
			},
		},
	})

	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      1,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{{
					CellX: 0,
					CellY: 0,
					Tile:  PreparedTile{TileID: 1, SheetID: 2, Source: img},
				}},
			},
		},
	}

	r.DrawTileLayer(scene, visibility.Camera{Width: 8, Height: 8})
	if got.Source != img {
		t.Fatal("expected prepared source image to be forwarded")
	}
	if got.SheetID != 2 || got.TileID != 1 {
		t.Fatalf("got tile = %+v, want prepared ids", got)
	}
}

func TestRendererUsesCompactPreparedDraws(t *testing.T) {
	var calls int
	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			tileFunc: func(x, y, width, height int, tile PreparedTile) {
				calls++
			},
		},
	})

	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      128,
					Height:     64,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 2, CellY: 1, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: image.NewRGBA(image.Rect(0, 0, 8, 8))}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{X: 16, Y: 8, Width: 32, Height: 32})
	if got.VisibleTiles != 1 {
		t.Fatalf("VisibleTiles = %d, want 1", got.VisibleTiles)
	}
	if calls != 1 {
		t.Fatalf("DrawTile calls = %d, want 1", calls)
	}
}

func TestRendererBatchesAdjacentPreparedDrawsIntoSingleCall(t *testing.T) {
	var runCalls int
	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			runFunc: func(x, y, tileWidth, tileHeight int, run PreparedRun) {
				runCalls++
				if run.Count != 2 {
					t.Fatalf("run.Count = %d, want 2", run.Count)
				}
			},
		},
	})

	src := image.NewRGBA(image.Rect(0, 0, 8, 8))
	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      4,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Runs: []PreparedRun{
					{CellX: 0, CellY: 0, Count: 2, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 32, Height: 8})
	if got.VisibleTiles != 2 {
		t.Fatalf("VisibleTiles = %d, want 2", got.VisibleTiles)
	}
	if runCalls != 1 {
		t.Fatalf("DrawRun calls = %d, want 1 batched call", runCalls)
	}
}

func TestRendererCanBlitPreparedRunDirectly(t *testing.T) {
	var runBlits int
	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			runFunc: func(x, y, tileWidth, tileHeight int, run PreparedRun) {
				runBlits++
				if run.Count != 2 {
					t.Fatalf("run.Count = %d, want 2", run.Count)
				}
			},
		},
	})

	src := image.NewRGBA(image.Rect(0, 0, 8, 8))
	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      4,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Runs: []PreparedRun{
					{CellX: 0, CellY: 0, Count: 2, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 32, Height: 8})
	if got.VisibleTiles != 2 {
		t.Fatalf("VisibleTiles = %d, want 2", got.VisibleTiles)
	}
	if runBlits != 1 {
		t.Fatalf("BlitTexturedRun calls = %d, want 1", runBlits)
	}
}

func TestRendererReusesTexturedRunSetupAcrossCompatibleRuns(t *testing.T) {
	var runBlits int
	src := image.NewRGBA(image.Rect(0, 0, 8, 8))

	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			runFunc: func(x, y, tileWidth, tileHeight int, run PreparedRun) {
				runBlits++
			},
		},
	})

	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      4,
					Height:     2,
					TileWidth:  8,
					TileHeight: 8,
				},
				Runs: []PreparedRun{
					{CellX: 0, CellY: 0, Count: 2, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
					{CellX: 0, CellY: 1, Count: 2, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 32, Height: 16})
	if got.VisibleTiles != 4 {
		t.Fatalf("VisibleTiles = %d, want 4", got.VisibleTiles)
	}
	if runBlits != 2 {
		t.Fatalf("BlitTexturedRun calls = %d, want 2", runBlits)
	}
}

func TestRendererReusesTexturedSetupAcrossCompatibleSingleDraws(t *testing.T) {
	var tileBlits int
	src := image.NewRGBA(image.Rect(0, 0, 8, 8))

	r := NewRenderer(RenderHooks{
		Executor: testExecutionBridge{
			tileFunc: func(x, y, width, height int, tile PreparedTile) {
				tileBlits++
			},
		},
	})

	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      2,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 0, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
					{CellX: 1, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 16, Height: 8})
	if got.VisibleTiles != 2 {
		t.Fatalf("VisibleTiles = %d, want 2", got.VisibleTiles)
	}
	if tileBlits != 2 {
		t.Fatalf("BlitTexturedTile calls = %d, want 2", tileBlits)
	}
}

func TestRendererUploadCountTracksPreparedTextureSourceChanges(t *testing.T) {
	srcA := texture.NewTextureFromImage(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	srcB := texture.NewTextureFromImage(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	r := NewRenderer(RenderHooks{})
	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      2,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 0, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: srcA}},
					{CellX: 1, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: srcB}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 16, Height: 8})
	if got.Uploads != 2 {
		t.Fatalf("Uploads = %d, want 2 for distinct prepared sources", got.Uploads)
	}
}

func TestRendererUploadCountIgnoresFallbackImageDraws(t *testing.T) {
	r := NewRenderer(RenderHooks{})
	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      2,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 0, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: image.NewRGBA(image.Rect(0, 0, 8, 8))}},
					{CellX: 1, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: image.NewRGBA(image.Rect(0, 0, 8, 8))}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 16, Height: 8})
	if got.Uploads != 0 {
		t.Fatalf("Uploads = %d, want 0 for fallback image draws", got.Uploads)
	}
}

func TestRendererUploadCountDoesNotRestartAcrossLayersForSamePreparedSource(t *testing.T) {
	src := texture.NewTextureFromImage(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	r := NewRenderer(RenderHooks{})
	scene := PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      1,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 0, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
				},
			},
			{
				Map: visibility.MapInfo{
					Width:      1,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 0, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 1, Source: src}},
				},
			},
		},
	}

	got := r.DrawTileLayer(scene, visibility.Camera{Width: 8, Height: 8})
	if got.Uploads != 1 {
		t.Fatalf("Uploads = %d, want 1 when prepared source is reused across layers", got.Uploads)
	}
}

func testScene() PreparedScene {
	return PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      4,
					Height:     4,
					TileWidth:  8,
					TileHeight: 8,
				},
			},
		},
	}
}

func testCamera() visibility.Camera {
	return visibility.Camera{X: 128, Y: 128, Width: 32, Height: 32}
}

func adjacentSameSheetScene() PreparedScene {
	return PreparedScene{
		Layers: []PreparedLayer{
			{
				Map: visibility.MapInfo{
					Width:      4,
					Height:     1,
					TileWidth:  8,
					TileHeight: 8,
				},
				Draws: []PreparedDraw{
					{CellX: 0, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 7}},
					{CellX: 1, CellY: 0, Tile: PreparedTile{TileID: 1, SheetID: 7}},
				},
			},
		},
	}
}
