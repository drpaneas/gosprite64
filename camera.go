package gosprite64

import "github.com/drpaneas/gosprite64/internal/rendergeom"

type Camera struct {
	X, Y          int
	Width, Height int
}

func newDefaultCamera() *Camera {
	bounds := rendergeom.LogicalBounds()
	return &Camera{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}
}
