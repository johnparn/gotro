package effects

// Ported from stars3.cpp: https://github.com/johangardhage/retro-demoeffects

import (
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const N_STARS int = 1000

var speed int = 2
var zmin int = -250
var zmax int = 250

type star struct {
	z int
	x int
	y int
}

var stars2 [N_STARS]star

type stars2Settings struct {
	WindowWidth  int
	WindowHeight int
	Renderer     *sdl.Renderer
	Colors       [256]RGBColor
}

var settings = stars2Settings{}

func RunStars2() {

	for i := 1; i < NUM_STARS; i++ {
		if stars2[i].z < zmin {
			stars2[i].z = rand.Intn(zmax)
		} else {
			stars2[i].z -= speed
		}

		var x int = (settings.WindowWidth / 2) + (stars2[i].x*255.0)/(stars2[i].z+255.0)
		var y int = (settings.WindowHeight / 2) + (stars2[i].y*255.0)/(stars2[i].z+255.0)
		var z int = -stars2[i].z

		if x >= 0 && x < settings.WindowWidth && y >= 0 && y < settings.WindowHeight {
			var colorIndex uint8 = uint8(float64(z+int(math.Abs(float64(zmin)))) * 0.128)
			PutPixel(settings.Renderer, int32(x), int32(y), settings.Colors[colorIndex].R.Value, settings.Colors[colorIndex].G.Value, settings.Colors[colorIndex].B.Value, 255)
		}
	}
}

func (s *stars2Settings) setColorBuffer(colorIndex, r, g, b, a uint8) {
	s.Colors[colorIndex] = RGBColor{R: Color{Value: r}, G: Color{Value: g}, B: Color{Value: b}, A: Color{Value: a}}
}

func InitializeStars2(renderer *sdl.Renderer, windowWidth int, windowHeight int) {
	settings.Renderer = renderer
	settings.WindowWidth = windowWidth
	settings.WindowHeight = windowHeight

	var i uint8
	for i = range 64 {
		settings.setColorBuffer(uint8(i), uint8(i*4), uint8(i*4), uint8(i*4), uint8(255))
	}

	for i := 1; i < NUM_STARS; i++ {
		stars2[i].x = rand.Intn(settings.WindowWidth) - (settings.WindowWidth / 2)
		stars2[i].y = rand.Intn(settings.WindowHeight) - (settings.WindowHeight / 2)
		stars2[i].z = rand.Intn(zmax)
	}
}
