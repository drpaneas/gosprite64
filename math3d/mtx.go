package math3d

import "math"

// Mat4 is a 4x4 matrix in row-major order using float32.
// This is the working format; use ToN64Mtx to convert to hardware format.
type Mat4 [4][4]float32

// Identity returns a 4x4 identity matrix.
func Identity() Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Mul multiplies two 4x4 matrices: result = a * b.
func (a Mat4) Mul(b Mat4) Mat4 {
	var out Mat4
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			out[r][c] = a[r][0]*b[0][c] + a[r][1]*b[1][c] + a[r][2]*b[2][c] + a[r][3]*b[3][c]
		}
	}
	return out
}

// MulVec4 transforms a Vec4 by the matrix.
func (m Mat4) MulVec4(v Vec4) Vec4 {
	return Vec4{
		X: m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		Y: m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		Z: m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		W: m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
}

// Perspective builds a perspective projection matrix matching libultra's
// guPerspectiveF. fovy is in degrees, aspect = width/height.
// Returns the matrix and the perspNorm value needed by the RSP.
func Perspective(fovy, aspect, near, far, scale float32) (Mat4, uint16) {
	m := Identity()
	fovyRad := fovy * math.Pi / 180.0
	yscale := float32(math.Cos(float64(fovyRad/2)) / math.Sin(float64(fovyRad/2)))
	m[0][0] = yscale / aspect
	m[1][1] = yscale
	m[2][2] = (near + far) / (near - far)
	m[2][3] = -1
	m[3][2] = 2 * near * far / (near - far)
	m[3][3] = 0
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			m[r][c] *= scale
		}
	}
	var perspNorm uint16
	if near+far <= 2.0 {
		perspNorm = 65535
	} else {
		pn := float64(1<<17) / float64(near+far)
		if pn <= 0 {
			perspNorm = 1
		} else {
			perspNorm = uint16(pn)
		}
	}
	return m, perspNorm
}

// Ortho builds an orthographic projection matrix matching libultra's guOrthoF.
func Ortho(left, right, bottom, top, near, far, scale float32) Mat4 {
	m := Identity()
	m[0][0] = 2 / (right - left)
	m[1][1] = 2 / (top - bottom)
	m[2][2] = -2 / (far - near)
	m[3][0] = -(right + left) / (right - left)
	m[3][1] = -(top + bottom) / (top - bottom)
	m[3][2] = -(far + near) / (far - near)
	m[3][3] = 1
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			m[r][c] *= scale
		}
	}
	return m
}

// LookAt builds a view matrix matching libultra's guLookAtReflectF
// (without the LookAt reflect light struct output).
func LookAt(eyeX, eyeY, eyeZ, atX, atY, atZ, upX, upY, upZ float32) Mat4 {
	m := Identity()

	xLook := atX - eyeX
	yLook := atY - eyeY
	zLook := atZ - eyeZ

	invLen := float32(-1.0 / math.Sqrt(float64(xLook*xLook+yLook*yLook+zLook*zLook)))
	xLook *= invLen
	yLook *= invLen
	zLook *= invLen

	xRight := upY*zLook - upZ*yLook
	yRight := upZ*xLook - upX*zLook
	zRight := upX*yLook - upY*xLook
	invLen = float32(1.0 / math.Sqrt(float64(xRight*xRight+yRight*yRight+zRight*zRight)))
	xRight *= invLen
	yRight *= invLen
	zRight *= invLen

	upX = yLook*zRight - zLook*yRight
	upY = zLook*xRight - xLook*zRight
	upZ = xLook*yRight - yLook*xRight
	invLen = float32(1.0 / math.Sqrt(float64(upX*upX+upY*upY+upZ*upZ)))
	upX *= invLen
	upY *= invLen
	upZ *= invLen

	m[0][0] = xRight
	m[1][0] = yRight
	m[2][0] = zRight
	m[3][0] = -(eyeX*xRight + eyeY*yRight + eyeZ*zRight)

	m[0][1] = upX
	m[1][1] = upY
	m[2][1] = upZ
	m[3][1] = -(eyeX*upX + eyeY*upY + eyeZ*upZ)

	m[0][2] = xLook
	m[1][2] = yLook
	m[2][2] = zLook
	m[3][2] = -(eyeX*xLook + eyeY*yLook + eyeZ*zLook)

	m[0][3] = 0
	m[1][3] = 0
	m[2][3] = 0
	m[3][3] = 1
	return m
}

// Translate builds a translation matrix.
func Translate(x, y, z float32) Mat4 {
	m := Identity()
	m[3][0] = x
	m[3][1] = y
	m[3][2] = z
	return m
}

// Scale builds a scale matrix.
func Scale(x, y, z float32) Mat4 {
	m := Identity()
	m[0][0] = x
	m[1][1] = y
	m[2][2] = z
	return m
}

// Rotate builds a rotation matrix. angle is in degrees, (x,y,z) is the axis.
// Matches libultra's guRotateF.
func Rotate(angle, x, y, z float32) Mat4 {
	v := Vec3{x, y, z}.Normalize()
	x, y, z = v.X, v.Y, v.Z

	a := angle * math.Pi / 180.0
	sinA := float32(math.Sin(float64(a)))
	cosA := float32(math.Cos(float64(a)))
	t := 1 - cosA

	ab := x * y * t
	bc := y * z * t
	ca := z * x * t

	m := Identity()

	xx := x * x
	m[0][0] = xx + cosA*(1-xx)
	m[2][1] = bc - x*sinA
	m[1][2] = bc + x*sinA

	yy := y * y
	m[1][1] = yy + cosA*(1-yy)
	m[2][0] = ca + y*sinA
	m[0][2] = ca - y*sinA

	zz := z * z
	m[2][2] = zz + cosA*(1-zz)
	m[1][0] = ab - z*sinA
	m[0][1] = ab + z*sinA

	return m
}

// N64Mtx is the N64 hardware fixed-point matrix format (s15.16).
// First 8 words are integer portions, last 8 words are fractional portions.
// Total size: 64 bytes (16 x uint32, but stored as [4][4]int32 per the SDK).
type N64Mtx [16]uint32

// ToN64Mtx converts a float32 Mat4 to the N64 fixed-point Mtx format.
// This matches libultra's guMtxF2L exactly.
func (mf Mat4) ToN64Mtx() N64Mtx {
	var mtx N64Mtx
	for r := 0; r < 4; r++ {
		for c := 0; c < 2; c++ {
			tmp1 := int32(mf[r][2*c] * 65536.0)
			tmp2 := int32(mf[r][2*c+1] * 65536.0)
			u1 := uint32(tmp1)
			u2 := uint32(tmp2)
			mtx[r*2+c] = (u1 & 0xFFFF0000) | ((u2 >> 16) & 0xFFFF)
			mtx[8+r*2+c] = ((u1 << 16) & 0xFFFF0000) | (u2 & 0xFFFF)
		}
	}
	return mtx
}

// FromN64Mtx converts an N64 fixed-point Mtx back to float32 Mat4.
// This matches libultra's guMtxL2F.
func FromN64Mtx(mtx N64Mtx) Mat4 {
	var mf Mat4
	for r := 0; r < 4; r++ {
		for c := 0; c < 2; c++ {
			intPart := mtx[r*2+c]
			fracPart := mtx[8+r*2+c]
			tmp1 := int32((intPart & 0xFFFF0000) | ((fracPart >> 16) & 0xFFFF))
			tmp2 := int32(((intPart << 16) & 0xFFFF0000) | (fracPart & 0xFFFF))
			mf[r][c*2] = float32(tmp1) / 65536.0
			mf[r][c*2+1] = float32(tmp2) / 65536.0
		}
	}
	return mf
}
