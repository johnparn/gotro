package effects

import (
	"github.com/veandco/go-sdl2/sdl"
)

func PutPixel(renderer *sdl.Renderer, x int32, y int32, r uint8, g uint8, b uint8, a uint8) {
	renderer.SetDrawColor(r, g, b, a)
	renderer.DrawPoint(x, y)
}
