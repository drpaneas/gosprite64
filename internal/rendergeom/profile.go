package rendergeom

import "image"

const (
	logicalWidth      = 288
	logicalHeight     = 216
	framebufferWidth  = 320
	framebufferHeight = 240
	originX           = 16
	originY           = 12
)

var (
	logicalBounds     = image.Rect(0, 0, logicalWidth, logicalHeight)
	framebufferBounds = image.Rect(0, 0, framebufferWidth, framebufferHeight)
	origin            = image.Pt(originX, originY)
)

// LogicalBounds returns the public 288x216 logical authoring canvas.
func LogicalBounds() image.Rectangle {
	return logicalBounds
}

// FramebufferBounds returns the internal 320x240 framebuffer bounds.
func FramebufferBounds() image.Rectangle {
	return framebufferBounds
}

// Origin returns the logical canvas origin inside the framebuffer.
func Origin() image.Point {
	return origin
}

// MapPoint maps a logical point into framebuffer space.
func MapPoint(p image.Point) (image.Point, bool) {
	if !p.In(logicalBounds) {
		return image.Point{}, false
	}

	return p.Add(origin), true
}

// MapRectInclusive clips a logical rectangle and maps it into framebuffer space.
// The input and output rectangles use inclusive bottom-right semantics for Max.
func MapRectInclusive(r image.Rectangle) (image.Rectangle, bool) {
	minX, minY, maxX, maxY, ok := clipRectInclusive(r)
	if !ok {
		return image.Rectangle{}, false
	}

	return image.Rectangle{
		Min: image.Pt(minX+origin.X, minY+origin.Y),
		Max: image.Pt(maxX+origin.X, maxY+origin.Y),
	}, true
}

// CenteredRect returns size centered inside bounds.
func CenteredRect(bounds image.Rectangle, size image.Point) image.Rectangle {
	min := image.Pt(
		bounds.Min.X+(bounds.Dx()-size.X)/2,
		bounds.Min.Y+(bounds.Dy()-size.Y)/2,
	)
	return image.Rectangle{Min: min, Max: min.Add(size)}
}

func clipRectInclusive(r image.Rectangle) (minX, minY, maxX, maxY int, ok bool) {
	minX, minY = r.Min.X, r.Min.Y
	maxX, maxY = r.Max.X, r.Max.Y
	if minX > maxX || minY > maxY {
		return 0, 0, 0, 0, false
	}

	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX >= logicalWidth {
		maxX = logicalWidth - 1
	}
	if maxY >= logicalHeight {
		maxY = logicalHeight - 1
	}
	if minX > maxX || minY > maxY {
		return 0, 0, 0, 0, false
	}

	return minX, minY, maxX, maxY, true
}
