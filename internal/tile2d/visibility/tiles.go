package visibility

type Camera struct {
	X, Y          int
	Width, Height int
}

type MapInfo struct {
	Width, Height         int
	TileWidth, TileHeight int
}

type VisibleBounds struct {
	MinX, MinY int
	MaxX, MaxY int
}

func VisibleCellBounds(cam Camera, m MapInfo) VisibleBounds {
	if m.Width <= 0 || m.Height <= 0 || m.TileWidth <= 0 || m.TileHeight <= 0 || cam.Width <= 0 || cam.Height <= 0 {
		return VisibleBounds{}
	}

	minX := max(0, cam.X/m.TileWidth)
	minY := max(0, cam.Y/m.TileHeight)

	maxX := min(m.Width, (cam.X+cam.Width+m.TileWidth-1)/m.TileWidth)
	maxY := min(m.Height, (cam.Y+cam.Height+m.TileHeight-1)/m.TileHeight)

	if minX > m.Width {
		minX = m.Width
	}
	if minY > m.Height {
		minY = m.Height
	}
	if maxX < minX {
		maxX = minX
	}
	if maxY < minY {
		maxY = minY
	}

	return VisibleBounds{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}
}
