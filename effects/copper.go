package effects

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	copperLineCount  = 9
	copperLineHeight = 9
)

type copperSettings struct {
	renderer  *sdl.Renderer
	width     int
	height    int
	startedAt time.Time
}

var copper = copperSettings{}

// InitializeCopper prepares the moving copper-line effect for the given renderer size.
func InitializeCopper(renderer *sdl.Renderer, windowWidth int32, windowHeight int32) {
	copper.renderer = renderer
	copper.width = int(windowWidth)
	copper.height = int(windowHeight)
	copper.startedAt = time.Now()
}

// RunCopper draws copper-colored lines whose vertical positions follow Lissajous motion.
func RunCopper() {
	if copper.renderer == nil || copper.width <= 0 || copper.height <= 0 {
		return
	}

	elapsed := time.Since(copper.startedAt).Seconds()
	centerY := float64(copper.height) / 2
	amplitude := float64(copper.height) * 0.25
	for line := 0; line < copperLineCount; line++ {
		phase := 2 * float64(line) / copperLineCount
		yMotion := math.Sin(elapsed*0.4+phase) * math.Sin(elapsed*1.1+phase*0.5)
		yPosition := centerY + amplitude*yMotion
		y := int32(math.Round(yPosition))

		colorPhase := elapsed*0.4 + phase
		red := uint8(80 + 10*(0.2+0.5*math.Sin(colorPhase)))
		green := uint8(20 + 20*(0.2+0.5*math.Sin(colorPhase+math.Pi/3)))
		blue := uint8(180 + 30*(0.2+0.5*math.Sin(colorPhase+2*math.Pi/3)))

		verticalProgress := 1 - math.Abs(yPosition-centerY)/amplitude
		if verticalProgress < 0 {
			verticalProgress = 0
		}
		edgeColor := sdl.Color{
			R: scaleCopperColor(red, 0.35+0.65*verticalProgress),
			G: scaleCopperColor(green, 0.35+0.65*verticalProgress),
			B: scaleCopperColor(blue, 0.35+0.65*verticalProgress),
			A: 255,
		}
		// fmt.Printf("Copper height: %d", copperLineHeight)
		drawCopper(y, sdl.Color{R: 255, G: 255, B: 255, A: 255}, edgeColor, copperLineHeight+(line*2))
	}
}

func drawCopper(y int32, centerColor sdl.Color, edgeColor sdl.Color, height int) {
	if height <= 0 {
		return
	}

	center := float64(height-1) / 2
	maxDistance := center
	for offset := 0; offset < height; offset++ {
		distance := 0.0
		if maxDistance > 0 {
			distance = math.Abs(float64(offset)-center) / maxDistance
		}
		color := interpolateCopperColor(centerColor, edgeColor, distance)

		copper.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		copper.renderer.FillRect(&sdl.Rect{
			X: 0,
			Y: y + int32(offset),
			W: int32(copper.width),
			H: 1,
		})
	}
}

func interpolateCopperColor(start sdl.Color, end sdl.Color, amount float64) sdl.Color {
	interpolate := func(from uint8, to uint8) uint8 {
		return uint8(math.Round(float64(from) + (float64(to)-float64(from))*amount))
	}

	return sdl.Color{
		R: interpolate(start.R, end.R),
		G: interpolate(start.G, end.G),
		B: interpolate(start.B, end.B),
		A: interpolate(start.A, end.A),
	}
}

func scaleCopperColor(channel uint8, shade float64) uint8 {
	return uint8(float64(channel) * shade)
}
