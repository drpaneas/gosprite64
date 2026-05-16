package gosprite64

import (
	"image"
	"image/color"

	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

// DrawRect draws the outline of a rectangle using DrawLine.
// x1, y1: Top-left corner (0-287, 0-215)
// x2, y2: Bottom-right corner (inclusive)
// c: The color of the rectangle outline
func DrawRect(x1, y1, x2, y2 int, c color.Color) {
	// Ensure x1 <= x2 and y1 <= y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Draw top line
	DrawLine(x1, y1, x2, y1, c)
	// Draw bottom line
	DrawLine(x1, y2, x2, y2, c)
	// Draw left and right side lines (offset by 1 to avoid double-drawing corners)
	if y2-y1 >= 1 {
		DrawLine(x1, y1+1, x1, y2-1, c)
		DrawLine(x2, y1+1, x2, y2-1, c)
	}
}

// DrawLine draws a line from (x1,y1) to (x2,y2) using FillRect.
// The line will be 1 pixel thick.
// Coordinates are expressed in the logical 288x216 canvas.
func DrawLine(x1, y1, x2, y2 int, c color.Color) {
	// For horizontal lines
	if y1 == y2 {
		if x1 > x2 {
			x1, x2 = x2, x1 // Ensure x1 <= x2
		}
		FillRect(x1, y1, x2, y1, c)
		return
	}

	// For vertical lines
	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1 // Ensure y1 <= y2
		}
		FillRect(x1, y1, x1, y2, c)
		return
	}

	// For diagonal lines, we'll draw a series of 1x1 rectangles
	// This is a simple implementation using Bresenham's algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		FillRect(x1, y1, x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// FillRect draws a filled rectangle on screen.
// x1, y1: Top-left corner (0-287, 0-215)
// x2, y2: Bottom-right corner (inclusive)
// c: The color to fill the rectangle with
func FillRect(x1, y1, x2, y2 int, c color.Color) {
	framebufferRect, ok := rendergeom.MapRectInclusive(image.Rectangle{
		Min: image.Pt(x1, y1),
		Max: image.Pt(x2, y2),
	})
	if !ok {
		return
	}

	drawFramebufferRect(
		framebufferRect.Min.X,
		framebufferRect.Min.Y,
		framebufferRect.Max.X,
		framebufferRect.Max.Y,
		c,
	)
}

func DrawImage(src image.Image, x, y int) {
	drawLogicalImage(src, x, y)
}

func DrawWorldImage(src image.Image, worldX, worldY int, cam *Camera) {
	if cam == nil {
		DrawImage(src, worldX, worldY)
		return
	}
	drawLogicalImage(src, worldX-cam.X, worldY-cam.Y)
}

func drawFramebufferRect(x1, y1, x2, y2 int, c color.Color) {
	video := currentVideo()
	if video == nil || video.Framebuffer == nil {
		return
	}

	// Ensure x1 <= x2 and y1 <= y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	framebufferBounds := rendergeom.FramebufferBounds()

	// Clip to screen bounds
	if x1 < framebufferBounds.Min.X {
		x1 = framebufferBounds.Min.X
	}
	if y1 < framebufferBounds.Min.Y {
		y1 = framebufferBounds.Min.Y
	}
	if x2 >= framebufferBounds.Max.X {
		x2 = framebufferBounds.Max.X - 1
	}
	if y2 >= framebufferBounds.Max.Y {
		y2 = framebufferBounds.Max.Y - 1
	}
	if x1 > x2 || y1 > y2 {
		return
	}

	// Create and draw the rectangle
	// Note: image.Rect is half-open interval [Min, Max)
	rect := image.Rect(x1, y1, x2+1, y2+1)
	img := video.uniform(c)
	n64draw.Src.Draw(video.Framebuffer, rect, img, image.Point{})
}

func drawLogicalImage(src image.Image, x, y int) {
	video := currentVideo()
	if video == nil || video.Framebuffer == nil || src == nil {
		return
	}

	srcBounds := src.Bounds()
	logicalDst := image.Rect(x, y, x+srcBounds.Dx(), y+srcBounds.Dy())
	clipped := logicalDst.Intersect(rendergeom.LogicalBounds())
	if clipped.Empty() {
		return
	}

	framebufferRect, ok := rendergeom.MapRectInclusive(image.Rectangle{
		Min: clipped.Min,
		Max: clipped.Max.Sub(image.Pt(1, 1)),
	})
	if !ok {
		return
	}

	srcPt := image.Pt(
		srcBounds.Min.X+(clipped.Min.X-logicalDst.Min.X),
		srcBounds.Min.Y+(clipped.Min.Y-logicalDst.Min.Y),
	)
	n64draw.Src.Draw(
		video.Framebuffer,
		image.Rect(framebufferRect.Min.X, framebufferRect.Min.Y, framebufferRect.Max.X+1, framebufferRect.Max.Y+1),
		src,
		srcPt,
	)
}
