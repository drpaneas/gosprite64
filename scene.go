package gosprite64

import (
	"fmt"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
	tileloader "github.com/drpaneas/gosprite64/internal/tile2d/loader"
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
	tilestats "github.com/drpaneas/gosprite64/internal/tile2d/stats"
	"github.com/drpaneas/gosprite64/internal/tile2d/visibility"
)

type Scene struct {
	bundle        *Bundle
	gameMap       *Map
	sheets        []*Sheet
	animations    []*AnimationSet
	defaultCamera *Camera
	preparer      *sceneRenderPreparer
	renderer      *tilerender.Renderer
	bridge        *sceneRenderBridge
	renderScene   tilerender.PreparedScene
	lastDrawStats tilerender.DrawStats
	cachedParsed  []format.ParsedSheet
	staticStats   RuntimeStats
}

func LoadScene(bundle *Bundle) (*Scene, error) {
	if bundle == nil {
		return nil, fmt.Errorf("load scene: nil bundle")
	}

	scene := &Scene{bundle: bundle}
	for _, entry := range bundle.manifest.Entries {
		switch entry.Kind {
		case format.BundleKindSheet:
			sheet, err := bundle.loadSheetEntry(entry)
			if err != nil {
				return nil, err
			}
			scene.sheets = append(scene.sheets, sheet)
		case format.BundleKindMap:
			if scene.gameMap != nil {
				return nil, fmt.Errorf("load scene: bundle has multiple maps")
			}
			m, err := bundle.loadMapEntry(entry)
			if err != nil {
				return nil, err
			}
			scene.gameMap = m
		case format.BundleKindAnim:
			anim, err := bundle.loadAnimEntry(entry)
			if err != nil {
				return nil, err
			}
			scene.animations = append(scene.animations, anim)
		}
	}

	scene.defaultCamera = newDefaultCamera()
	scene.preparer = newSceneRenderPreparer(scene)
	scene.bridge = newSceneRenderBridge()
	if rt := currentRuntime(); rt != nil {
		if tile := rt.currentTile(); tile != nil && tile.renderer != nil {
			scene.renderer = tile.renderer
		}
	}
	if scene.renderer == nil {
		scene.renderer = tilerender.NewRenderer(tilerender.RenderHooks{})
	}
	if scene.gameMap == nil {
		return nil, fmt.Errorf("load scene: bundle has no map")
	}

	if err := tileloader.ValidateSceneAssets(scene.gameMap.parsed, collectParsedSheets(scene.sheets)); err != nil {
		return nil, err
	}

	scene.configureRenderer()
	scene.renderScene = scene.preparer.buildScene()

	scene.cachedParsed = collectParsedSheets(scene.sheets)
	initSnap := tilestats.FromSceneAssets(scene.gameMap.parsed, scene.cachedParsed, tilerender.DrawStats{})
	scene.staticStats = RuntimeStats{
		SheetRAMBytes: initSnap.SheetRAMBytes,
		MapRAMBytes:   initSnap.MapRAMBytes,
		SheetCount:    initSnap.SheetCount,
		LayerCount:    initSnap.LayerCount,
	}

	return scene, nil
}

func (s *Scene) Map() *Map {
	if s == nil {
		return nil
	}
	return s.gameMap
}

func (s *Scene) SheetCount() int {
	if s == nil {
		return 0
	}
	return len(s.sheets)
}

func (s *Scene) Sheet(index int) *Sheet {
	if s == nil || index < 0 || index >= len(s.sheets) {
		return nil
	}
	return s.sheets[index]
}

func (s *Scene) SheetByID(id uint16) *Sheet {
	if id == 0 {
		return nil
	}
	return s.Sheet(int(id) - 1)
}

func (s *Scene) AnimationCount() int {
	if s == nil {
		return 0
	}
	return len(s.animations)
}

func (s *Scene) Animation(index int) *AnimationSet {
	if s == nil || index < 0 || index >= len(s.animations) {
		return nil
	}
	return s.animations[index]
}

func (s *Scene) AnimationByName(name string) *AnimationSet {
	if s == nil {
		return nil
	}
	for _, anim := range s.animations {
		if anim != nil && anim.name == name {
			return anim
		}
	}
	return nil
}

func (s *Scene) LayerSheet(layer int) (*Sheet, bool) {
	if s == nil {
		return nil, false
	}
	info, ok := s.Map().LayerInfo(layer)
	if !ok {
		return nil, false
	}
	sheet := s.SheetByID(info.SheetID)
	if sheet == nil {
		return nil, false
	}
	return sheet, true
}

func (s *Scene) LayerAssets(layer int) (MapLayerInfo, *Sheet, bool) {
	if s == nil {
		return MapLayerInfo{}, nil, false
	}
	info, ok := s.Map().LayerInfo(layer)
	if !ok {
		return MapLayerInfo{}, nil, false
	}
	sheet := s.SheetByID(info.SheetID)
	if sheet == nil {
		return MapLayerInfo{}, nil, false
	}
	return info, sheet, true
}

func (s *Scene) LayerSheetInfo(layer int) (SheetInfo, bool) {
	sheet, ok := s.LayerSheet(layer)
	if !ok {
		return SheetInfo{}, false
	}
	return sheet.Info(), true
}

func (s *Scene) Draw(cam *Camera) {
	if s == nil || s.renderer == nil || s.gameMap == nil {
		return
	}
	if cam == nil {
		cam = s.defaultCamera
	}
	if cam == nil {
		return
	}
	s.lastDrawStats = s.renderer.DrawPreparedScene(
		s.renderScene,
		visibility.Camera{
			X:      cam.X,
			Y:      cam.Y,
			Width:  cam.Width,
			Height: cam.Height,
		},
	)
}

func (s *Scene) configureRenderer() {
	if s == nil || s.renderer == nil || s.bridge == nil {
		return
	}
	s.renderer.SetHooks(tilerender.RenderHooks{
		Executor: s.bridge,
	})
}

func collectParsedSheets(sheets []*Sheet) []format.ParsedSheet {
	parsed := make([]format.ParsedSheet, 0, len(sheets))
	for _, sheet := range sheets {
		if sheet == nil {
			continue
		}
		parsed = append(parsed, sheet.parsed)
	}
	return parsed
}

func (s *Scene) Stats() RuntimeStats {
	if s == nil || s.gameMap == nil {
		return RuntimeStats{}
	}
	stats := s.staticStats
	stats.VisibleTiles = s.lastDrawStats.VisibleTiles
	stats.UploadCount = s.lastDrawStats.Uploads
	return stats
}
