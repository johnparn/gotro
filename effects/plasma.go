package effects

// Adapted from https://github.com/Aionmagan/Plasma-Effect-SDL/blob/master/src/main.c

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Vec2 struct {
	X float64
	Y float64
}

type PlasmaSettings struct {
	Coor         Vec2
	WindowWidth  int
	WindowHeight int
	Plasma       [][]uint8
	Palette      [256]sdl.Color
	Renderer     *sdl.Renderer
}

var plasmaSettings = PlasmaSettings{}

func getColor(x int, y int, windowWidth int, windowHeight int) uint8 {
	return uint8(128.0 + (128.0 * math.Sin(float64(x)/128.0)) + 128.0 + (128.0 * math.Sin(float64(y)/16.0)) + 128.0 + (128.0 * math.Sin(math.Sqrt(float64((x-windowWidth/2)*(x-windowWidth/2)+(y-windowHeight/2)*(y-windowHeight/2)))/8)) + 128.0 + (128.0 * math.Sin(math.Sqrt(float64(x*x+y*y))/8)))
}

// Pregenerate plasma, palette
func InitPlasma(renderer *sdl.Renderer, windowWidth int, windowHeight int) {

	plasmaSettings.Renderer = renderer
	plasmaSettings.WindowWidth = windowWidth
	plasmaSettings.WindowHeight = windowHeight

	// Pregenerate plasma
	// 2D slices - 1st dimension equals window width, 2nd equals window height
	plasmaEffect := make([][]uint8, windowWidth)

	for i := range plasmaEffect {
		plasmaEffect[i] = make([]uint8, windowHeight)
	}

	for i := range 256 {
		// Palette 2
		r := uint8(128.0 + 128*math.Sin(math.Pi*float64(i)/192.0))
		g := uint8(128.0 + 128*math.Sin(math.Pi*float64(i)/128.0))
		b := uint8(128.0 + 128*math.Sin(math.Pi*float64(i)/64.0))

		// Palette 3
		// 	r := uint8(128.0 + 128*math.Sin(math.Pi*float64(i)/16.0))
		// 	g := uint8(128.0 + 128*math.Cos(math.Pi*float64(i)/96.0))
		// 	b := 32 // uint8(128.0 + 128*math.Cos(math.Pi*float64(i)/192.0))

		plasmaSettings.Palette[i] = sdl.Color{R: r, G: g, B: b, A: 255}
	}

	// Create plasma
	for x := range plasmaSettings.WindowWidth {
		for y := range plasmaSettings.WindowHeight {
			color := getColor(x, y, plasmaSettings.WindowWidth, plasmaSettings.WindowHeight)
			plasmaEffect[x][y] = color
		}
	}
	plasmaSettings.Plasma = plasmaEffect
}

func RunPlasma() {

	paletteShift := int(sdl.GetTicks64() / 8)
	//draw every pixel, with the shifted palette color
	for y := 0; y < plasmaSettings.WindowHeight; y++ {
		for x := 0; x < plasmaSettings.WindowWidth; x++ {
			color := plasmaSettings.Palette[(int(plasmaSettings.Plasma[x][y])+paletteShift)%256]
			PutPixel(plasmaSettings.Renderer, int32(x), int32(y), color.R, color.G, color.B, color.A)
		}
	}
}
