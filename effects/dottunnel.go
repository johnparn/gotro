// Dot sphere
// Ported from dottunnel2.cpp: https://github.com/johangardhage/retro-demoeffects

package effects

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	DOTTUNNEL_MAXDEGREES float64 = 256
	DOTTUNNEL_POINTSTEP  float64 = 5
	DOTTUNNEL_NUMCIRCLES float64 = 87
	DOTTUNNEL_ZSTEP      int32   = 5
	DOTTUNNEL_ZMIN       int32   = -240
	DOTTUNNEL_RADIUS     float64 = 50.0
	DOTTUNNEL_DIVD       float64 = 128
	DOTTUNNEL_NUMPOINTS  int     = 256 / 5 // DOTTUNNEL_MAXDEGREES / DOTTUNNEL_POINTSTEP
	DOTTUNNEL_ZMAX       float64 = float64(DOTTUNNEL_NUMCIRCLES) * float64(DOTTUNNEL_ZSTEP)
)

// // Structs
type TunnelSettings struct {
	Renderer     *sdl.Renderer
	WindowHeight int32
	WindowWidth  int32
}

type Point2Df struct {
	x float64
	y float64
}

var Circle [DOTTUNNEL_NUMPOINTS]Point2Df
var counter float64 = 0
var tunnelSettings = TunnelSettings{}
var renderer *sdl.Renderer

func RenderDotTunnel(deltatime float64) {

	counter -= 1
	frame := counter

	var i int32 = 0
	for ; i < int32(DOTTUNNEL_NUMCIRCLES); i++ {
		var xo float64 = math.Cos(((frame*2.0+float64(i)*3.0)*2.0*math.Pi/DOTTUNNEL_MAXDEGREES))*DOTTUNNEL_DIVD/4 + math.Sin(((frame+float64(i)*2)*2.0*math.Pi/DOTTUNNEL_MAXDEGREES))*DOTTUNNEL_DIVD/3
		var yo float64 = math.Cos(((frame*2+float64(i)*2)*2.0*math.Pi/DOTTUNNEL_MAXDEGREES))*DOTTUNNEL_DIVD/5 + math.Sin(((frame*2+float64(i)*3)*2.0*math.Pi/DOTTUNNEL_MAXDEGREES))*DOTTUNNEL_DIVD/4

		var z float64 = float64(DOTTUNNEL_ZMIN + i*DOTTUNNEL_ZSTEP)
		var j int = 0

		for ; j < DOTTUNNEL_NUMPOINTS; j++ {
			// Projection coordinates
			var x = int32(float64(tunnelSettings.WindowWidth/2) + (Circle[j].x*250)/(z-250) + xo)
			var y = int32(float64(tunnelSettings.WindowHeight/2) + (Circle[j].y*250)/(z-250) + yo)

			if x >= 0 && x < tunnelSettings.WindowWidth && y >= 0 && y < tunnelSettings.WindowHeight {
				var color uint8 = uint8(math.Floor((z - float64(DOTTUNNEL_ZMIN)) * (128.0 / DOTTUNNEL_ZMAX)))
				PutPixel(tunnelSettings.Renderer, x, y, color, color, color, 255)
			}
		}
	}
}

func InitializeTunnel(sdl_renderer *sdl.Renderer, windowWidth int32, windowHeight int32) {

	tunnelSettings.Renderer = sdl_renderer
	tunnelSettings.WindowHeight = windowHeight
	tunnelSettings.WindowWidth = windowHeight

	// Init palette
	var paletteIndex uint8 = 0
	for ; paletteIndex < 64; paletteIndex++ {
		SetColor(paletteIndex, paletteIndex*4, paletteIndex*4, paletteIndex*4, 255)
	}

	// Init circle points
	var circleIndex int
	for ; circleIndex < DOTTUNNEL_NUMPOINTS; circleIndex++ {
		Circle[circleIndex].x = DOTTUNNEL_RADIUS * math.Cos(float64(circleIndex)*float64(DOTTUNNEL_POINTSTEP)*float64(2.0)*math.Pi/float64(DOTTUNNEL_MAXDEGREES)) * DOTTUNNEL_DIVD / (DOTTUNNEL_DIVD - 20)
		Circle[circleIndex].y = DOTTUNNEL_RADIUS * math.Sin(float64(circleIndex)*float64(DOTTUNNEL_POINTSTEP)*float64(2.0)*math.Pi/float64(DOTTUNNEL_MAXDEGREES)) * DOTTUNNEL_DIVD / (DOTTUNNEL_DIVD)
	}
}
