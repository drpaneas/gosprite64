package stats

import (
	"github.com/drpaneas/gosprite64/internal/tile2d/format"
	"github.com/drpaneas/gosprite64/internal/tile2d/residency"
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
)

type Snapshot struct {
	SheetRAMBytes int
	MapRAMBytes   int
	CachedChunks  int
	VisibleTiles  int
	SheetCount    int
	LayerCount    int
	UploadCount   int
}

func FromChunkCache(c interface{ Len() int }) Snapshot {
	if c == nil {
		return Snapshot{}
	}
	return Snapshot{CachedChunks: c.Len()}
}

func FromResidencyStore(s *residency.Store) Snapshot {
	_ = s
	return Snapshot{}
}

func FromSceneAssets(m format.ParsedMap, sheets []format.ParsedSheet, draw tilerender.DrawStats) Snapshot {
	var snap Snapshot
	for _, sheet := range sheets {
		snap.SheetRAMBytes += len(sheet.Pixels)
	}
	snap.SheetCount = len(sheets)
	snap.LayerCount = len(m.Layers)
	for _, layer := range m.Layers {
		snap.MapRAMBytes += len(layer.Cells) * 2
	}
	snap.VisibleTiles = draw.VisibleTiles
	snap.UploadCount = draw.Uploads
	return snap
}
