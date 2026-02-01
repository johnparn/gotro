// Adapted from https://gist.github.com/rkibria/8a94ae7ecbbf8abc06a165d9ebfaa6f2
package effects

import (
	"image/color"

	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func subVec3(a []float64, b []float64) []float64 {
	return []float64{a[0] - b[0], a[1] - b[1], a[2] - b[2], 0}
}

func dotProduct3(a []float64, b []float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func crossProduct3(a []float64, b []float64) []float64 {
	return []float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
		0}
}

func matMatMult(m4_1 []float64, m4_2 []float64) []float64 {
	result := make([]float64, 16)
	result[0] += m4_1[0] * m4_2[0] // row 1 x column 1
	result[0] += m4_1[1] * m4_2[4]
	result[0] += m4_1[2] * m4_2[8]
	result[0] += m4_1[3] * m4_2[12]
	result[1] += m4_1[0] * m4_2[1] // row 1 x column 2
	result[1] += m4_1[1] * m4_2[5]
	result[1] += m4_1[2] * m4_2[9]
	result[1] += m4_1[3] * m4_2[13]
	result[2] += m4_1[0] * m4_2[2]
	result[2] += m4_1[1] * m4_2[6]
	result[2] += m4_1[2] * m4_2[10]
	result[2] += m4_1[3] * m4_2[14]
	result[3] += m4_1[0] * m4_2[3]
	result[3] += m4_1[1] * m4_2[7]
	result[3] += m4_1[2] * m4_2[11]
	result[3] += m4_1[3] * m4_2[15]
	result[4] += m4_1[4] * m4_2[0]
	result[4] += m4_1[5] * m4_2[4]
	result[4] += m4_1[6] * m4_2[8]
	result[4] += m4_1[7] * m4_2[12]
	result[5] += m4_1[4] * m4_2[1]
	result[5] += m4_1[5] * m4_2[5]
	result[5] += m4_1[6] * m4_2[9]
	result[5] += m4_1[7] * m4_2[13]
	result[6] += m4_1[4] * m4_2[2]
	result[6] += m4_1[5] * m4_2[6]
	result[6] += m4_1[6] * m4_2[10]
	result[6] += m4_1[7] * m4_2[14]
	result[7] += m4_1[4] * m4_2[3]
	result[7] += m4_1[5] * m4_2[7]
	result[7] += m4_1[6] * m4_2[11]
	result[7] += m4_1[7] * m4_2[15]
	result[8] += m4_1[8] * m4_2[0]
	result[8] += m4_1[9] * m4_2[4]
	result[8] += m4_1[10] * m4_2[8]
	result[8] += m4_1[11] * m4_2[12]
	result[9] += m4_1[8] * m4_2[1]
	result[9] += m4_1[9] * m4_2[5]
	result[9] += m4_1[10] * m4_2[9]
	result[9] += m4_1[11] * m4_2[13]
	result[10] += m4_1[8] * m4_2[2]
	result[10] += m4_1[9] * m4_2[6]
	result[10] += m4_1[10] * m4_2[10]
	result[10] += m4_1[11] * m4_2[14]
	result[11] += m4_1[8] * m4_2[3]
	result[11] += m4_1[9] * m4_2[7]
	result[11] += m4_1[10] * m4_2[11]
	result[11] += m4_1[11] * m4_2[15]
	result[12] += m4_1[12] * m4_2[0]
	result[12] += m4_1[13] * m4_2[4]
	result[12] += m4_1[14] * m4_2[8]
	result[12] += m4_1[15] * m4_2[12]
	result[13] += m4_1[12] * m4_2[1]
	result[13] += m4_1[13] * m4_2[5]
	result[13] += m4_1[14] * m4_2[9]
	result[13] += m4_1[15] * m4_2[13]
	result[14] += m4_1[12] * m4_2[2]
	result[14] += m4_1[13] * m4_2[6]
	result[14] += m4_1[14] * m4_2[10]
	result[14] += m4_1[15] * m4_2[14]
	result[15] += m4_1[12] * m4_2[3]
	result[15] += m4_1[13] * m4_2[7]
	result[15] += m4_1[14] * m4_2[11]
	result[15] += m4_1[15] * m4_2[15]
	return result
}

func vecMatMult(v []float64, m []float64) []float64 {
	return []float64{
		m[0]*v[0] + m[1]*v[1] + m[2]*v[2] + m[3]*v[3],
		m[4]*v[0] + m[5]*v[1] + m[6]*v[2] + m[7]*v[3],
		m[8]*v[0] + m[9]*v[1] + m[10]*v[2] + m[11]*v[3],
		m[12]*v[0] + m[13]*v[1] + m[14]*v[2] + m[15]*v[3],
	}
}

func getTranslMat4(dx float64, dy float64, dz float64) []float64 {
	return []float64{
		1.0, 0.0, 0.0, dx,
		0.0, 1.0, 0.0, dy,
		0.0, 0.0, 1.0, dz,
		0.0, 0.0, 0.0, 1.0,
	}
}

func getRotXMat4(phi float64) []float64 {
	cos_phi := math.Cos(phi)
	sin_phi := math.Sin(phi)
	return []float64{
		1.0, 0.0, 0.0, 0.0,
		0.0, cos_phi, -sin_phi, 0.0,
		0.0, sin_phi, cos_phi, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func getRotYMat4(phi float64) []float64 {
	cos_phi := math.Cos(phi)
	sin_phi := math.Sin(phi)
	return []float64{
		cos_phi, 0.0, sin_phi, 0.0,
		0.0, 1.0, 0.0, 0.0,
		-sin_phi, 0.0, cos_phi, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func getRotZMat4(phi float64) []float64 {
	cos_phi := math.Cos(phi)
	sin_phi := math.Sin(phi)
	return []float64{
		cos_phi, -sin_phi, 0.0, 0.0,
		sin_phi, cos_phi, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func projectVerts(m []float64, srcV [][]float64) [][]float64 {
	var result [][]float64
	for _, v3 := range srcV {
		v4 := []float64{v3[0], v3[1], v3[2], 1}
		pvec := vecMatMult(v4, m)
		result = append(result, pvec)
	}
	return result
}

func perspDiv(srcV4 [][]float64) [][]float64 {
	var result [][]float64
	for _, vec4 := range srcV4 {
		result = append(result, []float64{vec4[0] / vec4[2], vec4[1] / vec4[2]})
	}
	return result
}

func cullBackfaces(viewPoint []float64, srcIdcs [][]uint, worldVerts [][]float64) []int {
	var idcs []int
	for i, tri := range srcIdcs {
		v0 := worldVerts[tri[0]]
		v1 := worldVerts[tri[1]]
		v2 := worldVerts[tri[2]]
		viewPointToTriVec := subVec3(v0, viewPoint)
		s10 := subVec3(v1, v0)
		s20 := subVec3(v2, v0)
		normal := crossProduct3(s10, s20)
		dp := dotProduct3(viewPointToTriVec, normal)
		if dp > 0 {
			idcs = append(idcs, i)
		}
	}
	return idcs
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// Bresenham implementation from https://github.com/eaburns/quart/blob/master/geom/2d_draw.go#L164
func drawEdge(renderer *sdl.Renderer, w, h int, c color.Color, p0, p1 []float64) {
	o_x := w / 2
	o_y := h / 2
	ar := float64(w) / float64(h)

	x0 := o_x + int(p0[0]*float64(o_x))
	y0 := o_y - int(p0[1]*float64(o_y)*ar)
	x1 := o_x + int(p1[0]*float64(o_x))
	y1 := o_y - int(p1[1]*float64(o_y)*ar)

	// Bresenham's alg: http://en.wikipedia.org/wiki/Bresenham's_line_algorithm
	steep := abs(y0-y1) > abs(x0-x1)
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := abs(y1 - y0)
	err := dx / 2
	y := y0

	ystep := -1
	if y0 < y1 {
		ystep = 1
	}

	for x := x0; x <= x1; x++ {

		if steep {
			renderer.SetDrawColor(0, 255, 0, 1)
			renderer.DrawPoint(int32(y), int32(x))
			// img.Set(y, x, c)
		} else {
			renderer.SetDrawColor(0, 255, 0, 1)
			renderer.DrawPoint(int32(x), int32(y))
			// img.Set(x, y, c)
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
}

func drawObject(renderer *sdl.Renderer, w, h int, c color.Color,
	drawIdxList []int, meshTris [][]uint, perspVerts [][]float64) {

	for _, triIdx := range drawIdxList {
		tri := meshTris[triIdx]
		i0 := tri[0]
		i1 := tri[1]
		i2 := tri[2]
		p0 := perspVerts[i0]
		p1 := perspVerts[i1]
		p2 := perspVerts[i2]

		drawEdge(renderer, w, h, c, p0, p1)
		drawEdge(renderer, w, h, c, p0, p2)
		drawEdge(renderer, w, h, c, p1, p2)
	}
}

func SpinningCube(renderer *sdl.Renderer, windowWidth int, windowHeight int) {

	w := windowWidth
	h := windowHeight

	// OpenGL coord system: -z goes into the screen
	// 8 vertices form corners of the unit cube
	var cubeVerts = [][]float64{
		{0.5, 0.5, 0.5},    // front top right     0
		{0.5, -0.5, 0.5},   // front bottom right  1
		{-0.5, -0.5, 0.5},  // front bottom left   2
		{-0.5, 0.5, 0.5},   // front top left      3
		{0.5, 0.5, -0.5},   // back top right      4
		{0.5, -0.5, -0.5},  // back bottom right   5
		{-0.5, -0.5, -0.5}, // back bottom left    6
		{-0.5, 0.5, -0.5},  // back top left       7
	}

	var cubeIdxs = [][]uint{
		{0, 1, 3}, // front face  3 0
		{2, 3, 1}, //             2 1
		{3, 2, 7}, // left face   7 3
		{6, 7, 2}, //             6 2
		{4, 5, 0}, // right face  0 4
		{1, 0, 5}, //             1 5
		{4, 0, 7}, // top face    7 4
		{3, 7, 0}, //             3 0
		{1, 5, 2}, // bottom face 2 1
		{6, 2, 5}, //             6 5
		{7, 6, 4}, // back face   4 7
		{5, 4, 6}, //             5 6
	}

	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0x00}, color.RGBA{0x00, 0xff, 0x00, 0xff},
	}

	viewPoint := []float64{0, 0, 0, 0}

	angle := math.Pi / float64(180) * float64(sdl.GetTicks64()/10)
	mx := getRotXMat4(angle)
	my := getRotYMat4(angle)
	mz := getRotZMat4(angle)
	mt := getTranslMat4(0, 0, -2)
	m1 := matMatMult(my, mx)
	m2 := matMatMult(mz, m1)
	m3 := matMatMult(mt, m2)
	worldVerts := projectVerts(m3, cubeVerts[:])
	drawIdxList := cullBackfaces(viewPoint, cubeIdxs, worldVerts)
	perspVerts := perspDiv(worldVerts)
	drawObject(renderer, w, h, palette[1], drawIdxList, cubeIdxs, perspVerts)
}
