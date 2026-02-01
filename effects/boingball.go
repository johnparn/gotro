package effects

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type boingballSettings struct {
	Renderer     *sdl.Renderer // sdl renderer instance
	WindowWidth  int32         // Width of window
	WindowHeight int32         // Height of window
	XPos         int32         // Position X of sprite
	YPos         int32         // Position Y of sprite
	XIncrement   int32
	YIncrement   int32
	BGColors     RGBColor // Background color
	Step         int
	Sprite       *sdl.Surface
	StartTick    uint64
}

// type Option func(p boingballSettings) boingballSettings

var bbSettings boingballSettings

const ObjectWidth int32 = 128
const ObjectHeight int32 = 128

func drawSprite(x, y int32, bgColors RGBColor) {
	src := sdl.Rect{W: 455, H: 456}
	dst := sdl.Rect{X: x, Y: y, W: ObjectWidth, H: ObjectHeight}
	sprite, _ := img.Load("boingball.png")
	bbSettings.Sprite = sprite
	texture, _ := bbSettings.Renderer.CreateTextureFromSurface(sprite)
	// _ = renderer.SetDrawColor(R, G, B, 255)
	_ = bbSettings.Renderer.Clear()
	_ = bbSettings.Renderer.Copy(texture, &src, &dst)
}

func InitBoingball(renderer *sdl.Renderer, windowWidth int32, windowHeight int32, R, G, B uint8) {

	var bgColors = RGBColor{
		R: Color{Value: uint8(R), Increment: -1},
		G: Color{Value: uint8(G), Increment: 1},
		B: Color{Value: uint8(B), Increment: -1},
	}

	bbSettings = boingballSettings{}
	bbSettings.Renderer = renderer
	bbSettings.WindowWidth = (windowWidth)
	bbSettings.WindowHeight = (windowHeight)
	bbSettings.BGColors = bgColors
	bbSettings.XPos = -128
	bbSettings.YPos = -128
	bbSettings.XIncrement = 2
	bbSettings.YIncrement = 2
	bbSettings.StartTick = sdl.GetTicks64()
}

func setBGColor(bgColors RGBColor) RGBColor {
	bgColors.R = GetColor(bgColors.R)
	bgColors.G = GetColor(bgColors.G)
	bgColors.B = GetColor(bgColors.B)
	// bbSettings.Renderer.SetDrawColor(bgColors.R.Value, bgColors.G.Value, bgColors.B.Value, 1)
	// _ = bbSettings.Renderer.Clear()
	return bgColors
}

func BoingBall() {

	bbSettings.Renderer.SetDrawColor(bbSettings.BGColors.R.Value, bbSettings.BGColors.G.Value, bbSettings.BGColors.B.Value, 0)
	bbSettings.Renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: bbSettings.WindowWidth, H: bbSettings.WindowHeight})

	if bbSettings.XPos >= bbSettings.WindowWidth-ObjectWidth {
		bbSettings.XIncrement = -3
	} else if bbSettings.XPos <= 0 {
		bbSettings.XIncrement = 3
	}

	if bbSettings.YPos >= bbSettings.WindowHeight-ObjectHeight {
		bbSettings.YIncrement = -3
	} else if bbSettings.YPos <= 0 {
		bbSettings.YIncrement = 3
	}

	// Render object
	drawSprite(bbSettings.XPos, bbSettings.YPos, bbSettings.BGColors)

	// Update position
	bbSettings.XPos += bbSettings.XIncrement
	bbSettings.YPos += bbSettings.YIncrement

	// Change color
	bbSettings.BGColors = setBGColor(bbSettings.BGColors)

}
