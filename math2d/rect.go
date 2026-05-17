package math2d

// Rect is an axis-aligned rectangle defined by top-left corner and dimensions.
// Boundary convention is half-open: [X, X+W) x [Y, Y+H).
type Rect struct {
	X, Y, W, H float32
}

// RectFromCenter creates a Rect centered on the given point.
func RectFromCenter(center Vec2, w, h float32) Rect {
	return Rect{
		X: center.X - w/2,
		Y: center.Y - h/2,
		W: w,
		H: h,
	}
}

func (r Rect) Right() float32  { return r.X + r.W }
func (r Rect) Bottom() float32 { return r.Y + r.H }

func (r Rect) Center() Vec2 {
	return Vec2{X: r.X + r.W/2, Y: r.Y + r.H/2}
}

func (r Rect) ContainsPoint(p Vec2) bool {
	return r.W > 0 && r.H > 0 &&
		p.X >= r.X && p.X < r.Right() &&
		p.Y >= r.Y && p.Y < r.Bottom()
}

func (r Rect) ContainsRect(other Rect) bool {
	return other.W > 0 && other.H > 0 &&
		other.X >= r.X && other.Right() <= r.Right() &&
		other.Y >= r.Y && other.Bottom() <= r.Bottom()
}

func (r Rect) Overlaps(other Rect) bool {
	if r.W <= 0 || r.H <= 0 || other.W <= 0 || other.H <= 0 {
		return false
	}
	return r.X < other.Right() && r.Right() > other.X &&
		r.Y < other.Bottom() && r.Bottom() > other.Y
}

func (r Rect) Intersection(other Rect) (Rect, bool) {
	if !r.Overlaps(other) {
		return Rect{}, false
	}
	x := max32(r.X, other.X)
	y := max32(r.Y, other.Y)
	right := min32(r.Right(), other.Right())
	bottom := min32(r.Bottom(), other.Bottom())
	return Rect{X: x, Y: y, W: right - x, H: bottom - y}, true
}

// Expand returns a new Rect grown by amount on all sides.
func (r Rect) Expand(amount float32) Rect {
	return Rect{
		X: r.X - amount,
		Y: r.Y - amount,
		W: r.W + amount*2,
		H: r.H + amount*2,
	}
}

func min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
