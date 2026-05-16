package runtimeflow

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestBootstrapOrder(t *testing.T) {
	want := []Stage{
		StageConsole,
		StageVideo,
		StageGameInit,
		StageAudio,
		StageLoop,
	}

	if got := BootstrapOrder(); !slices.Equal(got, want) {
		t.Fatalf("BootstrapOrder() = %v, want %v", got, want)
	}
}

func TestStatusGuards(t *testing.T) {
	tests := []struct {
		name         string
		status       Status
		wantCanDraw  bool
		wantCanQueue bool
	}{
		{
			name:         "zero status blocks draw and audio",
			status:       Status{},
			wantCanDraw:  false,
			wantCanQueue: false,
		},
		{
			name:         "video ready allows draw only",
			status:       Status{VideoReady: true},
			wantCanDraw:  true,
			wantCanQueue: false,
		},
		{
			name:         "audio ready allows queue only",
			status:       Status{AudioReady: true},
			wantCanDraw:  false,
			wantCanQueue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.CanDraw(); got != tt.wantCanDraw {
				t.Fatalf("Status.CanDraw() = %v, want %v", got, tt.wantCanDraw)
			}

			if got := tt.status.CanQueueAudio(); got != tt.wantCanQueue {
				t.Fatalf("Status.CanQueueAudio() = %v, want %v", got, tt.wantCanQueue)
			}
		})
	}
}

func TestRuntimeBootstrapIsExplicit(t *testing.T) {
	logsSource := mustReadRepoFile(t, "logs.go")
	if strings.Contains(logsSource, "func init()") {
		t.Fatal("logs.go must not hide bootstrap in package init")
	}
	if !strings.Contains(logsSource, "func setupConsole()") {
		t.Fatal("logs.go must expose setupConsole")
	}

	runtimeSource := mustReadRepoFile(t, "runtime.go")
	for _, snippet := range []string{
		"type runtimeState struct {",
		"video *videoState",
		"audio *audioState",
		"tile  *tileRuntime",
		"var activeRuntime *runtimeState",
		"type tileRuntime struct {",
		"func newTileRuntime() *tileRuntime",
		"func newRuntimeState() *runtimeState",
		"func activateRuntime(rt *runtimeState)",
		"func currentRuntime() *runtimeState",
		"func (rt *runtimeState) currentVideo() *videoState",
		"func (rt *runtimeState) currentAudio() *audioState",
		"func (rt *runtimeState) currentTile() *tileRuntime",
	} {
		if !strings.Contains(runtimeSource, snippet) {
			t.Fatalf("runtime.go must contain %q", snippet)
		}
	}

	screenSource := mustReadRepoFile(t, "screen.go")
	for _, snippet := range []string{
		"type videoState struct {",
		"func newVideoState() *videoState",
		"func (rt *runtimeState) initVideo()",
		"func beginDrawing()",
		"func endDrawing()",
		"func ClearScreen()",
		"func ClearScreenWith(c color.Color)",
		"currentVideo()",
	} {
		if !strings.Contains(screenSource, snippet) {
			t.Fatalf("screen.go must contain %q", snippet)
		}
	}
	if strings.Contains(screenSource, "var currentScreen") {
		t.Fatal("screen.go must not keep global screen ownership")
	}

	shapesSource := mustReadRepoFile(t, "shapes.go")
	if !strings.Contains(shapesSource, "currentVideo()") {
		t.Fatal("shapes.go must route drawing through currentVideo()")
	}

	textSource := mustReadRepoFile(t, "text.go")
	if !strings.Contains(textSource, "currentVideo()") {
		t.Fatal("text.go must route drawing through currentVideo()")
	}

	audioSource := mustReadRepoFile(t, "audio.go")
	for _, snippet := range []string{
		"type audioConfig struct {",
		"type audioState struct {",
		"var pendingAudioConfig audioConfig",
		"func newAudioState(cfg audioConfig) *audioState",
		"func (a *audioState) ready() bool",
		"func (rt *runtimeState) initAudio()",
		"func (a *audioState) start()",
		"func (a *audioState) feeder()",
		"currentAudio()",
	} {
		if !strings.Contains(audioSource, snippet) {
			t.Fatalf("audio.go must contain %q", snippet)
		}
	}

	gameLoopSource := mustReadRepoFile(t, "gameloop.go")
	assertOrderedSubstrings(t, gameLoopSource,
		"setupConsole()",
		"rt := newRuntimeState()",
		"rt.initVideo()",
		"activateRuntime(rt)",
		"rdp.RDP.SetScissor(",
		"g.Init()",
		"rt.initAudio()",
	)

	sceneSource := mustReadRepoFile(t, "scene.go")
	if !strings.Contains(sceneSource, `return nil, fmt.Errorf("load scene: bundle has multiple maps")`) {
		t.Fatal("scene.go must fail fast on multiple-map bundles")
	}
}

func TestTilemapExampleUsesCanonicalSceneAPI(t *testing.T) {
	example := mustReadRepoFile(t, "examples/tilemap/main.go")
	for _, snippet := range []string{
		`"github.com/drpaneas/gosprite64"`,
		"bundle, err := gosprite64.OpenBundle(",
		"scene, err := gosprite64.LoadScene(bundle)",
		"scene.Draw(camera)",
		"gosprite64.Run(&Game{})",
	} {
		if !strings.Contains(example, snippet) {
			t.Fatalf("examples/tilemap/main.go must contain %q", snippet)
		}
	}
}

func TestTilemapExampleAutoPansCamera(t *testing.T) {
	example := mustReadRepoFile(t, "examples/tilemap/main.go")
	for _, snippet := range []string{
		"type Game struct {",
		"tick   int",
		"g.tick++",
		"m := scene.Map()",
		"g.camera.X = pingPong(g.tick/2, max(0, m.PixelWidth()-g.camera.Width))",
		"g.camera.Y = pingPong(g.tick/3, max(0, m.PixelHeight()-g.camera.Height))",
		"func pingPong(step, limit int) int",
		"Width: 64, Height: 64",
	} {
		if !strings.Contains(example, snippet) {
			t.Fatalf("examples/tilemap/main.go must contain %q", snippet)
		}
	}
}

func TestTilemapExampleShipsBundleAsset(t *testing.T) {
	if _, err := os.Stat(repoFilePath(t, "examples/tilemap/assets/level.bundle")); err != nil {
		t.Fatalf("examples/tilemap/assets/level.bundle must exist: %v", err)
	}
}

func TestTilemapExampleExercisesMultipleSheets(t *testing.T) {
	example := mustReadRepoFile(t, "examples/tilemap/main.go")
	if !strings.Contains(example, "-sheet assets/tiles.sheet -sheet assets/tiles_overlay.sheet") {
		t.Fatal("examples/tilemap/main.go must bundle multiple sheet assets")
	}

	level := mustReadRepoFile(t, "examples/tilemap/assets-src/level.json")
	for _, snippet := range []string{
		`"layer_count": 2`,
		`"sheet_id": 2`,
	} {
		if !strings.Contains(level, snippet) {
			t.Fatalf("examples/tilemap/assets-src/level.json must contain %q", snippet)
		}
	}
}

func TestTilemapExampleDisplaysRuntimeStats(t *testing.T) {
	example := mustReadRepoFile(t, "examples/tilemap/main.go")
	for _, snippet := range []string{
		`"fmt"`,
		`"idle"`,
		"anim := scene.AnimationByName(\"idle\")",
		"anim.Clip(\"idle\")",
		"frameIdx, frameTile := animationCursor(g.tick, clip)",
		"baseLayer, baseSheet, _ := scene.LayerAssets(0)",
		"overlayLayer, _, _ := scene.LayerAssets(1)",
		"baseSheetInfo, _ := scene.LayerSheetInfo(0)",
		"overlaySheetInfo, _ := scene.LayerSheetInfo(1)",
		"animTile := baseSheet.Tile(frameTile + 1)",
		"gosprite64.DrawWorldImage(animTile, 56, 24, camera)",
		`fmt.Sprintf("an:%s %d/%d f:%d@%d", clip.Name, frameIdx+1, len(clip.Frames), frameTile, clip.FPS)`,
		`fmt.Sprintf("tc:%d/%d", baseSheetInfo.TileCount, overlaySheetInfo.TileCount)`,
		`fmt.Sprintf("sh:%d/%d", baseLayer.SheetID, overlayLayer.SheetID)`,
		`fmt.Sprintf("nz:%d/%d", baseLayer.NonZeroTiles, overlayLayer.NonZeroTiles)`,
		"stats := scene.Stats()",
		`fmt.Sprintf("vis:%d", stats.VisibleTiles)`,
		`fmt.Sprintf("up:%d", stats.UploadCount)`,
		`fmt.Sprintf("ram:%d/%d", stats.SheetRAMBytes, stats.MapRAMBytes)`,
		"gosprite64.DrawText(",
	} {
		if !strings.Contains(example, snippet) {
			t.Fatalf("examples/tilemap/main.go must contain %q", snippet)
		}
	}
}

func mustReadRepoFile(t *testing.T, name string) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	path := filepath.Join(filepath.Dir(file), "..", "..", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}

	return string(content)
}

func repoFilePath(t *testing.T, name string) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Join(filepath.Dir(file), "..", "..", name)
}

func assertOrderedSubstrings(t *testing.T, content string, snippets ...string) {
	t.Helper()

	last := -1
	for _, snippet := range snippets {
		index := strings.Index(content, snippet)
		if index == -1 {
			t.Fatalf("missing %q", snippet)
		}
		if index <= last {
			t.Fatalf("%q must appear after previous bootstrap step", snippet)
		}
		last = index
	}
}
