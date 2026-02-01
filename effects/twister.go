package effects

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

/*
 *  Adapted from
 *  Twister effect
 *  2019, Giulio Zausa
 *  gcc twister.cc -o twister -l stdc++ -l SDL2 -std=c++17 -Ofast
 *  https://github.com/giulioz/DemoEffects/blob/master/twister.cc
 */

const TWISTER_PI_2 float64 = math.Pi / 2
const TWISTER_SCALE float64 = 100.0

type ColorBuf [800 * 600 * 4]uint32

type TwisterSettings struct {
	Renderer     *sdl.Renderer
	WindowWidth  int32 // Width of window
	WindowHeight int32 // Height of window
	Offset       float64
}

var twisterSettings = TwisterSettings{}
var colorBuffer ColorBuf

var p float64 = 0.0
var deltaTime float64 = 0.0
var Y0 int32
var Y1 int32
var Y2 int32
var Y3 int32
var x int32 = 0
var divider float64 = 2000

func convCoord(x int32, y int32) int32 {
	return (y * twisterSettings.WindowWidth) + x
}

func blend(color1 uint32, color2 uint32, alpha uint32) uint32 {
	rb := color1 & 0xff00ff00
	b := color1 & 0x0000ff00
	g := color1 & 0x00ff0000

	rb += ((color2 & 0xff00ff) - rb) * alpha >> 24
	b += ((color2 & 0x00ff00) - b) * alpha >> 24
	g += ((color2 & 0xff00ff) - g) * alpha >> 24

	//return rb | g | b

	return (rb & 0xff00ff) | (g & 0xff00ff) | (b & 0x00ff00)
}

func InitTwister(sdl_renderer *sdl.Renderer, windowWidth int32, windowHeight int32) {
	twisterSettings.Renderer = sdl_renderer
	twisterSettings.WindowWidth = windowWidth
	twisterSettings.WindowHeight = windowHeight
	twisterSettings.Offset = float64(windowHeight / 2)
}

func vline(buf *ColorBuf, x int32, y0 float64, y1 float64, color uint32) {

	for y := y0; y < y1; y += 0.5 {
		i := convCoord(x, int32(y))
		d := (math.Abs(0.5-math.Pow((y-y0)/float64(y1-y0), 3))/2 + 0.75)
		buf[i] = blend(color, buf[i], uint32(d*255))

		var r uint8 = uint8(buf[i] >> 0)
		var g uint8 = uint8(buf[i] >> 8)
		var b uint8 = uint8(buf[i] >> 16)
		var a uint8 = uint8(buf[i] >> 24)

		PutPixel(twisterSettings.Renderer, x, int32(y), r, g, b, a)
	}

}

func RunTwister() {

	p = deltaTime / divider

	// Render each pixel
	for x = range twisterSettings.WindowWidth {

		Y0 := ((math.Sin(p) * TWISTER_SCALE) + twisterSettings.Offset)
		Y1 := ((math.Sin(p+(TWISTER_PI_2)) * TWISTER_SCALE) + twisterSettings.Offset)
		Y2 := ((math.Sin(p+math.Pi) * TWISTER_SCALE) + twisterSettings.Offset)
		Y3 := ((math.Sin(p+TWISTER_PI_2*3) * TWISTER_SCALE) + twisterSettings.Offset)

		vline(&colorBuffer, x, Y3, Y0, 0x00ffff00)
		vline(&colorBuffer, x, Y2, Y3, 0xffff0000)
		vline(&colorBuffer, x, Y1, Y2, 0xff00ff00)
		vline(&colorBuffer, x, Y0, Y1, 0xff000000)

		p += 0.005 * math.Sin(deltaTime/divider) * 3
	}

	deltaTime += 102
}
