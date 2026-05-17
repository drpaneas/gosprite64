package scene3d

import "github.com/drpaneas/gosprite64/math3d"

// CameraData holds perspective/orthographic projection parameters.
type CameraData struct {
	FOV    float32 // field of view in degrees (perspective)
	Aspect float32 // width/height ratio
	Near   float32
	Far    float32
	Ortho  bool // if true, use orthographic projection
	// Ortho extents (only used when Ortho is true)
	Left, Right, Bottom, Top float32
}

// NewPerspectiveCamera creates a camera node with perspective projection.
func NewPerspectiveCamera(name string, fov, aspect, near, far float32) *Node {
	return &Node{
		Type:    NodeCamera,
		Name:    name,
		Visible: true,
		Scale:   math3d.Vec3{X: 1, Y: 1, Z: 1},
		Camera: &CameraData{
			FOV:    fov,
			Aspect: aspect,
			Near:   near,
			Far:    far,
		},
	}
}

// NewOrthoCamera creates a camera node with orthographic projection.
func NewOrthoCamera(name string, left, right, bottom, top, near, far float32) *Node {
	return &Node{
		Type:    NodeCamera,
		Name:    name,
		Visible: true,
		Scale:   math3d.Vec3{X: 1, Y: 1, Z: 1},
		Camera: &CameraData{
			Ortho:  true,
			Left:   left,
			Right:  right,
			Bottom: bottom,
			Top:    top,
			Near:   near,
			Far:    far,
		},
	}
}

// ProjectionMatrix returns the projection matrix for this camera.
func (c *CameraData) ProjectionMatrix() (math3d.Mat4, uint16) {
	if c.Ortho {
		return math3d.Ortho(c.Left, c.Right, c.Bottom, c.Top, c.Near, c.Far, 1.0), 0
	}
	return math3d.Perspective(c.FOV, c.Aspect, c.Near, c.Far, 1.0)
}
