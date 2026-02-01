// Dot sphere
// Ported from dotball.cpp: https://github.com/johangardhage/retro-demoeffects

package effects

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const RADIUS float64 = 150.0
const MAXPOINTS int32 = 2048
const POINTSTEP float64 = 0.1
const ZMIN float64 = 20
const ZMAX float64 = 160

// Count number of points
var NumPoints int32 = 0

// var WindowWidth int32
// var WindowHeight int32
var Renderer *sdl.Renderer
var palette [256]RGBColor

type SphereSettings struct {
	Renderer     *sdl.Renderer
	WindowHeight int32
	WindowWidth  int32
}

var sphereSettings = SphereSettings{}

type Vertex struct {
	// Original coordinates
	x float64
	y float64
	z float64
	// Rotated coordinates
	rx float64
	ry float64
	rz float64
	// Screen coordinates
	sx float64
	sy float64
}

var Dot [MAXPOINTS]Vertex

func RotateVertex(vertex *Vertex, ax float64, ay float64, az float64) {
	// Rotate around x axis
	vertex.ry = vertex.y*math.Cos(ax) - vertex.z*math.Sin(ax)
	vertex.rz = vertex.y*math.Sin(ax) + vertex.z*math.Cos(ax)

	// Rotate around y axis
	vertex.rx = vertex.x*math.Cos(ay) + vertex.rz*math.Sin(ay)
	vertex.rz = vertex.x*-math.Sin(ay) + vertex.rz*math.Cos(ay)

	// Rotate around z axis
	tmpx := vertex.rx*math.Cos(az) - vertex.ry*math.Sin(az)
	vertex.ry = vertex.rx*math.Sin(az) + vertex.ry*math.Cos(az)
	vertex.rx = tmpx
}

func ProjectVertex(vertex *Vertex, scale float64) {
	x := sphereSettings.WindowWidth / 2
	y := sphereSettings.WindowHeight / 2
	eye := 250.0

	vertex.sx = float64(x) + (scale*vertex.rx*eye)/(scale*vertex.rz+eye)
	vertex.sy = float64(y) + (scale*vertex.ry*eye)/(scale*vertex.rz+eye)
}

var ax float64 = 1.0
var ay float64 = 1.0
var az float64 = 1.0

func RenderDotSphere(deltatime float64) {

	var distance float64 = 1
	ax += 0.01
	ay += 0.01
	az += 0.01

	var i int32 = 0

	for ; i < NumPoints; i++ {
		RotateVertex(&Dot[i], ax, ay, az)
		ProjectVertex(&Dot[i], distance)

		x := Dot[i].sx
		y := Dot[i].sy
		z := -math.Floor(Dot[i].rz)

		if x >= 0 && x <= float64(sphereSettings.WindowWidth) && y >= 0 && y <= float64(sphereSettings.WindowHeight) && z > ZMIN && z < ZMAX {
			color := uint8(math.Floor((z + math.Abs(ZMIN)) * (64.0 / (math.Abs(ZMIN) + ZMAX))))
			PutPixel(sphereSettings.Renderer, int32(x), int32(y), palette[color].R.Value, palette[color].G.Value, palette[color].B.Value, palette[color].A.Value)
			// Renderer.SetDrawColor(palette[color].R.Value, palette[color].G.Value, palette[color].B.Value, palette[color].A.Value)
			// Renderer.DrawPoint(int32(x), int32(y))
			// gfx.FilledCircleRGBA(Renderer, int32(x), int32(y), 2, palette[color].R.Value, palette[color].G.Value, palette[color].B.Value, palette[color].A.Value)
		}

	}
}

func InitializeDotSphere(renderer *sdl.Renderer, windowWidth int32, windowHeight int32) {

	sphereSettings.WindowWidth = windowWidth
	sphereSettings.WindowHeight = windowHeight
	sphereSettings.Renderer = renderer

	// Init palette
	var i uint8 = 0
	for ; i < 64; i++ {
		SetColor(i, i*4, i*4, i*4, i*4)
	}

	// Generate ball using dots
	for alpha := 2.0 * math.Pi; alpha > 0; alpha -= POINTSTEP {
		for beta := math.Pi; beta > 0; beta -= POINTSTEP {
			Dot[NumPoints].x = RADIUS * math.Cos(alpha) * math.Sin(beta)
			Dot[NumPoints].y = RADIUS * math.Cos(beta)
			Dot[NumPoints].z = RADIUS * math.Sin(alpha) * math.Sin(beta)
			NumPoints++
			if NumPoints > MAXPOINTS {
				return
			}
		}
	}
}
