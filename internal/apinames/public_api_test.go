package apinames

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCoreGameplayNaming(t *testing.T) {
	gameLoop := mustReadRepoFile(t, "gameloop.go")

	requireContains(t, gameLoop, "type Game interface {")
	requireContains(t, gameLoop, "func Run(g Game)")
	requireNotContains(t, gameLoop, "type Gamelooper interface")
	requireNotContains(t, gameLoop, "func Run(g Gamelooper)")
}

func TestExamplesUseQualifiedImports(t *testing.T) {
	clearscreen := mustReadRepoFile(t, "examples/clearscreen/main.go")
	requireContains(t, clearscreen, `"github.com/drpaneas/gosprite64"`)
	requireContains(t, clearscreen, "gosprite64.Run(&Game{})")
	requireNotContains(t, clearscreen, `import . "github.com/drpaneas/gosprite64"`)

	calibration := mustReadRepoFile(t, "examples/calibration/main.go")
	requireContains(t, calibration, `"github.com/drpaneas/gosprite64"`)
	requireContains(t, calibration, "gosprite64.Run(&Game{})")
	requireNotContains(t, calibration, `import . "github.com/drpaneas/gosprite64"`)
}

func TestHelloWorldGuideUsesQualifiedImports(t *testing.T) {
	guide := mustReadRepoFile(t, "docs/hello_world.md")
	requireContains(t, guide, `import "github.com/drpaneas/gosprite64"`)
	requireContains(t, guide, "gosprite64.Run(&Game{})")
	requireNotContains(t, guide, `import . "github.com/drpaneas/gosprite64"`)
	requireNotContains(t, guide, "Gamelooper")
}

func TestAudioGuideUsesQualifiedImports(t *testing.T) {
	guide := mustReadRepoFile(t, "docs/audio.md")
	requireNotContains(t, guide, `. "github.com/drpaneas/gosprite64"`)
	requireContains(t, guide, "gosprite64.PlaySoundEffect(")
}

func TestInputNaming(t *testing.T) {
	controls := mustReadRepoFile(t, "controls.go")

	for _, snippet := range []string{
		"ButtonA",
		"ButtonB",
		"ButtonZ",
		"ButtonStart",
		"ButtonDPadUp",
		"ButtonDPadDown",
		"ButtonDPadLeft",
		"ButtonDPadRight",
		"func IsButtonDown(",
		"func IsButtonJustPressed(",
		"func StickPosition(",
	} {
		requireContains(t, controls, snippet)
	}

	for _, snippet := range []string{
		"BtnA",
		"BtnB",
		"Btn(",
		"Btnp(",
		"GetStick(",
		"RIGHT =",
		"LEFT  =",
		"UP    =",
		"DOWN  =",
		"X     =",
		"O     =",
	} {
		requireNotContains(t, controls, snippet)
	}
}

func TestDrawHelperNaming(t *testing.T) {
	shapes := mustReadRepoFile(t, "shapes.go")
	requireContains(t, shapes, "func DrawRect(")
	requireContains(t, shapes, "func DrawLine(")
	requireContains(t, shapes, "func DrawImage(src image.Image, x, y int)")
	requireContains(t, shapes, "func DrawWorldImage(src image.Image, worldX, worldY int, cam *Camera)")
	requireContains(t, shapes, "func FillRect(")
	requireNotContains(t, shapes, "func Line(")
	requireNotContains(t, shapes, "func Rectfill(")

	gameLoop := mustReadRepoFile(t, "gameloop.go")
	requireNotContains(t, gameLoop, "func DrawRectFill(")

	text := mustReadRepoFile(t, "text.go")
	requireContains(t, text, "func DrawText(")
	requireNotContains(t, text, "func Print(")
}

func TestDocsTeachOnlyCanonicalDrawNames(t *testing.T) {
	for _, path := range []string{
		"README.md",
		"docs/introduction.md",
		"docs/hello_world.md",
		"docs/square_pixels.md",
	} {
		content := mustReadRepoFile(t, path)
		requireNotContains(t, content, "`Rectfill`")
		requireNotContains(t, content, "`DrawRectFill`")
		requireNotContains(t, content, "`Line`")
		requireNotContains(t, content, "`Print`")
	}
}

func TestMathShorthandHelpersAreGone(t *testing.T) {
	requireMissingRepoFile(t, "math.go")

	pong := mustReadRepoFile(t, "examples/pong/main.go")
	requireContains(t, pong, "rand.IntN(2)")
	requireNotContains(t, pong, "Flr(")
	requireNotContains(t, pong, "Rnd(")

	spaceInvaders := mustReadRepoFile(t, "examples/space_invaders/main.go")
	requireNotContains(t, spaceInvaders, "Btn(")
	requireNotContains(t, spaceInvaders, "Btnp(")
	requireNotContains(t, spaceInvaders, "Rnd(")
}

func TestAudioNaming(t *testing.T) {
	audioSource := mustReadRepoFile(t, "audio.go")
	for _, snippet := range []string{
		"type AudioBundle struct {",
		"func RegisterAudioBundle(bundle AudioBundle)",
		"func PlaySoundEffect(",
		"func PlayMusic(",
		"func StopMusic()",
		"func SetSoundEffectVolume(",
		"func SetMusicVolume(",
	} {
		requireContains(t, audioSource, snippet)
	}

	for _, snippet := range []string{
		"RegisterAudioV1",
		"RegisterSFXNameResolver",
		"PlayEffect(",
		"PlayTrack(",
		"StopTrack(",
		"SetEffectVolume(",
		"SetTrackVolume(",
	} {
		requireNotContains(t, audioSource, snippet)
	}
}

func TestTileEngineNaming(t *testing.T) {
	for _, path := range []string{"bundle.go", "scene.go", "camera.go", "sheet.go", "map.go", "sprite.go", "animation.go", "asset_policy.go"} {
		content := mustReadRepoFile(t, path)
		requireNotContains(t, content, "TMEM")
		requireNotContains(t, content, "DMA")
	}

	scene := mustReadRepoFile(t, "scene.go")
	requireContains(t, scene, "func LoadScene(bundle *Bundle) (*Scene, error)")
	requireContains(t, scene, "func (s *Scene) SheetCount() int")
	requireContains(t, scene, "func (s *Scene) Sheet(index int) *Sheet")
	requireContains(t, scene, "func (s *Scene) SheetByID(id uint16) *Sheet")
	requireContains(t, scene, "func (s *Scene) AnimationCount() int")
	requireContains(t, scene, "func (s *Scene) Animation(index int) *AnimationSet")
	requireContains(t, scene, "func (s *Scene) AnimationByName(name string) *AnimationSet")
	requireContains(t, scene, "func (s *Scene) LayerSheet(layer int) (*Sheet, bool)")
	requireContains(t, scene, "func (s *Scene) LayerAssets(layer int) (MapLayerInfo, *Sheet, bool)")
	requireContains(t, scene, "func (s *Scene) LayerSheetInfo(layer int) (SheetInfo, bool)")
	requireContains(t, scene, "func (s *Scene) Draw(cam *Camera)")

	m := mustReadRepoFile(t, "map.go")
	for _, snippet := range []string{
		"type MapLayerInfo struct {",
		"func (m *Map) Width() int",
		"func (m *Map) Height() int",
		"func (m *Map) TileWidth() int",
		"func (m *Map) TileHeight() int",
		"func (m *Map) LayerCount() int",
		"func (m *Map) LayerInfo(layer int) (MapLayerInfo, bool)",
		"func (m *Map) LayerSheetID(layer int) (uint16, bool)",
		"func (m *Map) TileAt(layer, x, y int) (uint16, bool)",
		"func (m *Map) PixelWidth() int",
		"func (m *Map) PixelHeight() int",
	} {
		requireContains(t, m, snippet)
	}

	sheet := mustReadRepoFile(t, "sheet.go")
	for _, snippet := range []string{
		"type SheetInfo struct {",
		"func (s *Sheet) Info() SheetInfo",
		"func (s *Sheet) Tile(tileID uint16) image.Image",
	} {
		requireContains(t, sheet, snippet)
	}

	animation := mustReadRepoFile(t, "animation.go")
	for _, snippet := range []string{
		"type AnimationSet struct {",
		"type AnimationClip struct {",
		"func (a *AnimationSet) Name() string",
		"func (a *AnimationSet) Clips() []AnimationClip",
		"func (a *AnimationSet) Clip(name string) (AnimationClip, bool)",
	} {
		requireContains(t, animation, snippet)
	}
}

func TestTileEngineCameraAndSceneHooks(t *testing.T) {
	camera := mustReadRepoFile(t, "camera.go")
	requireContains(t, camera, "type Camera struct {")
	requireContains(t, camera, "X, Y          int")
	requireContains(t, camera, "Width, Height int")

	scene := mustReadRepoFile(t, "scene.go")
	requireContains(t, scene, "defaultCamera *Camera")
	requireContains(t, scene, "preparer      *sceneRenderPreparer")
	requireContains(t, scene, "renderer      *tilerender.Renderer")
	requireContains(t, scene, "bridge        *sceneRenderBridge")
	requireContains(t, scene, "func (s *Scene) configureRenderer()")
	requireContains(t, scene, "scene.preparer = newSceneRenderPreparer(scene)")
	requireContains(t, scene, "scene.bridge = newSceneRenderBridge()")
	requireContains(t, scene, "scene.configureRenderer()")
	requireContains(t, scene, "func (s *Scene) Draw(cam *Camera)")
	requireContains(t, scene, "cam = s.defaultCamera")
	requireContains(t, scene, "s.lastDrawStats = s.renderer.DrawPreparedScene(")
}

func TestTileEnginePolicyControls(t *testing.T) {
	assetPolicy := mustReadRepoFile(t, "asset_policy.go")
	requireContains(t, assetPolicy, "type RuntimeStats struct {")
	requireContains(t, assetPolicy, "VisibleTiles  int")
	requireContains(t, assetPolicy, "SheetCount    int")
	requireContains(t, assetPolicy, "LayerCount    int")
	requireContains(t, assetPolicy, "UploadCount   int")
	requireNotContains(t, assetPolicy, "type CompressionMode")
	requireNotContains(t, assetPolicy, "type AssetPolicy struct")
	requireNotContains(t, assetPolicy, "type SceneOptions struct")
}

func TestTileEngineRuntimeStatsSurface(t *testing.T) {
	scene := mustReadRepoFile(t, "scene.go")
	requireContains(t, scene, "func (s *Scene) Stats() RuntimeStats")

	assetPolicy := mustReadRepoFile(t, "asset_policy.go")
	requireContains(t, assetPolicy, "type RuntimeStats struct {")
	requireContains(t, assetPolicy, "VisibleTiles  int")
	requireContains(t, assetPolicy, "SheetCount    int")
	requireContains(t, assetPolicy, "LayerCount    int")
	requireContains(t, assetPolicy, "UploadCount   int")
}

func TestTileEngineAssetRegistrationNaming(t *testing.T) {
	cartfsSource := mustReadRepoFile(t, "cartfs.go")
	requireContains(t, cartfsSource, "func RegisterAssetFS(f cartfs.FS)")

	example := mustReadRepoFile(t, "examples/tilemap/assets_embed.go")
	requireContains(t, example, "//go:embed assets/*")
	requireContains(t, example, "gosprite64.RegisterAssetFS(assetFS)")
}

func TestLegacyNamesAreGoneFromDocsAndExamples(t *testing.T) {
	legacy := []string{
		"Gamelooper",
		"Btn(",
		"Btnp(",
		"Flr(",
		"Rnd(",
		"Rectfill",
		"DrawRectFill",
		"`Line`",
		".Line(",
		"Print(",
		"RegisterAudioV1",
		"RegisterSFXNameResolver",
		"PlayEffect(",
		"PlayTrack(",
		"StopTrack(",
		"SetEffectVolume(",
		"SetTrackVolume(",
	}

	for _, path := range []string{
		"README.md",
		"docs/introduction.md",
		"docs/hello_world.md",
		"docs/square_pixels.md",
		"docs/audio.md",
		"examples/clearscreen/main.go",
		"examples/calibration/main.go",
		"examples/pong/main.go",
		"examples/space_invaders/main.go",
	} {
		content := mustReadRepoFile(t, path)
		for _, snippet := range legacy {
			requireNotContains(t, content, snippet)
		}
	}
}

func mustReadRepoFile(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(repoFilePath(t, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(data)
}

func repoFilePath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", name)
}

func requireMissingRepoFile(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Stat(repoFilePath(t, name)); err == nil {
		t.Fatalf("expected %s to be missing", name)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", name, err)
	}
}

func TestTilemapExampleUsesCanonicalAPIs(t *testing.T) {
	tilemap := mustReadRepoFile(t, "examples/tilemap/main.go")
	requireContains(t, tilemap, `"github.com/drpaneas/gosprite64"`)
	requireContains(t, tilemap, "gosprite64.OpenBundle(")
	requireContains(t, tilemap, "gosprite64.LoadScene(")
	requireContains(t, tilemap, "gosprite64.Camera{")
	requireContains(t, tilemap, "scene.Draw(")
	requireContains(t, tilemap, "gosprite64.ClearScreen()")
	requireContains(t, tilemap, "gosprite64.DrawText(")
	requireNotContains(t, tilemap, `import . "github.com/drpaneas/gosprite64"`)
}

func TestPublicTile2DFilesDoNotImportInternalPackages(t *testing.T) {
	publicFiles := []string{
		"bundle.go",
		"camera.go",
		"asset_policy.go",
		"sprite.go",
	}
	for _, path := range publicFiles {
		content := mustReadRepoFile(t, path)
		for _, banned := range []string{
			`"github.com/drpaneas/gosprite64/internal/tile2d/render"`,
			`"github.com/drpaneas/gosprite64/internal/tile2d/visibility"`,
			`"github.com/drpaneas/gosprite64/internal/tile2d/residency"`,
			`"github.com/drpaneas/gosprite64/internal/tile2d/stats"`,
		} {
			requireNotContains(t, content, banned)
		}
	}
}

func TestSpriteSheetAPI(t *testing.T) {
	ss := mustReadRepoFile(t, "sprite_sheet.go")
	requireContains(t, ss, "type SpriteSheet struct {")
	requireContains(t, ss, "func LoadSpriteSheet(path string) (*SpriteSheet, error)")
	requireContains(t, ss, "func (s *SpriteSheet) FrameCount() int")
	requireContains(t, ss, "func (s *SpriteSheet) FrameWidth() int")
	requireContains(t, ss, "func (s *SpriteSheet) FrameHeight() int")
	requireNotContains(t, ss, "image.Image")
}

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

func requireContains(t *testing.T, content, snippet string) {
	t.Helper()
	if !strings.Contains(content, snippet) {
		t.Fatalf("missing %q", snippet)
	}
}

func requireNotContains(t *testing.T, content, snippet string) {
	t.Helper()
	if strings.Contains(content, snippet) {
		t.Fatalf("unexpected %q", snippet)
	}
}
