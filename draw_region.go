package gosprite64

// DrawRegion defines a sub-rectangle of the screen that all drawing
// is clipped and offset to. Used for split-screen multiplayer.
type DrawRegion struct {
	X, Y, W, H int
}

// Active reports whether this region restricts drawing.
// A zero-value DrawRegion is inactive (full screen).
func (r DrawRegion) Active() bool {
	return r.W > 0 && r.H > 0
}

// Offset translates a local coordinate into screen space.
func (r DrawRegion) Offset(x, y int) (int, int) {
	return x + r.X, y + r.Y
}

// Clip offsets and clips a rectangle to the region bounds.
// Returns the screen-space coordinates and false if entirely outside.
func (r DrawRegion) Clip(x1, y1, x2, y2 int) (int, int, int, int, bool) {
	x1 += r.X
	y1 += r.Y
	x2 += r.X
	y2 += r.Y

	if r.Active() {
		if x1 < r.X {
			x1 = r.X
		}
		if y1 < r.Y {
			y1 = r.Y
		}
		rr := r.X + r.W
		rb := r.Y + r.H
		if x2 > rr {
			x2 = rr
		}
		if y2 > rb {
			y2 = rb
		}
		if x1 >= x2 || y1 >= y2 {
			return 0, 0, 0, 0, false
		}
	}
	return x1, y1, x2, y2, true
}

// ContainsPoint checks whether a local coordinate is within the region.
func (r DrawRegion) ContainsPoint(x, y int) bool {
	if !r.Active() {
		return true
	}
	sx, sy := r.Offset(x, y)
	return sx >= r.X && sx < r.X+r.W && sy >= r.Y && sy < r.Y+r.H
}

var drawRegionStack []DrawRegion

func currentDrawRegion() DrawRegion {
	if len(drawRegionStack) == 0 {
		return DrawRegion{}
	}
	return drawRegionStack[len(drawRegionStack)-1]
}

// SetDrawRegion restricts all subsequent drawing to the given screen-space
// rectangle. Coordinates passed to drawing functions become local to the
// region's top-left corner. Calls can be nested.
func SetDrawRegion(x, y, w, h int) {
	r := DrawRegion{X: x, Y: y, W: w, H: h}
	drawRegionStack = append(drawRegionStack, r)
	applyScissor(r)
}

// ResetDrawRegion removes the most recent draw region, restoring the
// previous one (or full screen if none remain).
func ResetDrawRegion() {
	if len(drawRegionStack) > 0 {
		drawRegionStack = drawRegionStack[:len(drawRegionStack)-1]
	}
	if len(drawRegionStack) > 0 {
		applyScissor(drawRegionStack[len(drawRegionStack)-1])
	} else {
		clearScissor()
	}
}
