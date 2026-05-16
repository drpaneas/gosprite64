package gosprite64

type RuntimeStats struct {
	SheetRAMBytes int
	MapRAMBytes   int
	CachedChunks  int
	VisibleTiles  int
	SheetCount    int
	LayerCount    int
	UploadCount   int
}
