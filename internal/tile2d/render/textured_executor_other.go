//go:build !n64

package render

func (e TexturedExecutor) EnsurePrepared(tile PreparedTile) bool {
	return false
}

func (e TexturedExecutor) BlitRun(x, y, tileWidth int, count int) {}

func (e TexturedExecutor) BlitTile(x, y int) {}
