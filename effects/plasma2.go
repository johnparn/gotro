package effects

// Ported from plasma.cpp: https://github.com/johangardhage/retro-demoeffects

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const PLASMA_FRAMES int = 720
const PLASMA_COLORS uint8 = 255

type Plasma2Settings struct {
	WindowWidth  int
	WindowHeight int
	Renderer     *sdl.Renderer
	SineValues   float64
	SineTable    []float64
	Colors       [256]RGBColor
}

var framecounter float64
var plasma2Settings = Plasma2Settings{}
var frame int

func RunPlasma2() {
	// Calculate frame
	framecounter += 9
	frame = int(math.Mod(framecounter, float64(PLASMA_FRAMES)))

	for y := 0; y <= plasma2Settings.WindowHeight; y++ {
		var yc float64 = 75 + plasma2Settings.SineTable[y+frame*2]*2 + plasma2Settings.SineTable[y*2+frame/2] + plasma2Settings.SineTable[y+frame]*2
		// var x int = 0
		for x := 0; x <= plasma2Settings.WindowWidth; x++ {
			var xc float64 = 75 + plasma2Settings.SineTable[x*2+frame/2] + plasma2Settings.SineTable[x+frame*2] + plasma2Settings.SineTable[x/2+frame]*2

			color := uint8(xc * yc)
			PutPixel(plasma2Settings.Renderer, int32(x), int32(y), plasma2Settings.Colors[color].R.Value,
				plasma2Settings.Colors[color].G.Value,
				plasma2Settings.Colors[color].B.Value,
				plasma2Settings.Colors[color].A.Value)
		}
	}
}

func (p *Plasma2Settings) setColorBuffer(colorIndex, r, g, b, a uint8) {
	p.Colors[colorIndex] = RGBColor{R: Color{Value: r}, G: Color{Value: g}, B: Color{Value: b}, A: Color{Value: a}}
}

func InitializePlasma2(renderer *sdl.Renderer, windowWidth int, windowHeight int) {
	plasma2Settings.Renderer = renderer
	plasma2Settings.WindowWidth = windowWidth
	plasma2Settings.WindowHeight = windowHeight
	plasma2Settings.SineValues = float64(windowWidth + PLASMA_FRAMES*2)

	// Init palette
	var r uint8 = 0
	var g uint8 = 0
	var b uint8 = 0

	var i uint8 = 0
	for ; i < PLASMA_COLORS; i++ {
		plasma2Settings.setColorBuffer(i, 0, 0, 0, 127)
	}
	for i = 0; i < 42; i++ {
		plasma2Settings.setColorBuffer(i, r*4, g*4, b*4, 127)
		r++
	}
	for i = 42; i < 84; i++ {
		plasma2Settings.setColorBuffer(i, r*4, g*4, b*4, 127)
		g++
	}
	for i = 84; i < 126; i++ {
		plasma2Settings.setColorBuffer(i, r*4, g*4, b*4, 127)
		b++
	}
	for i = 126; i < 168; i++ {
		plasma2Settings.setColorBuffer(i, r*4, g*4, b*4, 127)
		r--
	}

	for i = 168; i < 210; i++ {
		plasma2Settings.setColorBuffer(i, r*4, g*4, b*4, 127)
		g--
	}
	for i = 210; i < 252; i++ {
		plasma2Settings.setColorBuffer(i, r*4, g*4, b*4, 127)
		b--
	}

	// Init sine table
	var j float64 = 0
	for ; j < float64(plasma2Settings.SineValues); j++ {
		plasma2Settings.SineTable = append(plasma2Settings.SineTable, math.Cos(j*math.Pi/180))
	}

}
