package gosprite64

type CompressionMode uint8

const (
	CompressionDefault CompressionMode = iota
	CompressionDisabled
)

type AssetPolicy struct {
	Persistent  bool
	Compression CompressionMode
}

type SceneOptions struct {
	AssetPolicy AssetPolicy
}

type RuntimeStats struct {
	SheetRAMBytes int
	MapRAMBytes   int
	CachedChunks  int
	VisibleTiles  int
	SheetCount    int
	LayerCount    int
	UploadCount   int
}
