package effects

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const NUM_STARS int = 900

type Star struct {
	X, Y, speed int
}

type StarsSettings struct {
	NumberOfStars int
	Renderer      *sdl.Renderer
	WindowWidth   int
	WindowHeight  int
}

var stars [NUM_STARS]Star

// int getStarColor(int);

var starSettings StarsSettings

func ResetStars(renderer *sdl.Renderer, windowWidth, windowHeight int) {

	renderer.SetDrawColor(15, 15, 23, 1)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(windowWidth), H: int32(windowHeight)})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < NUM_STARS; i++ {
		stars[i].X = rand.Intn(windowWidth) % 20
		stars[i].Y = rand.Intn(windowHeight) % 20
		stars[i].speed = 1 + (rand.Intn(60) % 12)
	}
}

func InitStars(renderer *sdl.Renderer, windowWidth, windowHeight int) {
	starSettings = StarsSettings{}
	starSettings.Renderer = renderer
	starSettings.WindowWidth = windowWidth
	starSettings.WindowHeight = windowHeight
	starSettings.NumberOfStars = NUM_STARS

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < NUM_STARS; i++ {
		stars[i].X -= stars[i].speed

		if stars[i].X < 0 {
			stars[i].X = rand.Intn(starSettings.WindowWidth)
			stars[i].Y = rand.Intn(starSettings.WindowHeight)
			stars[i].speed = 1 + (rand.Intn(60) % 12)
		}
	}

}

func UpdateStars() {
	var rect sdl.Rect

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < NUM_STARS; i++ {
		if stars[i].X < starSettings.WindowWidth {
			rect.X = int32(stars[i].X)
			rect.Y = int32(stars[i].Y)
			rect.W = int32(1)
			rect.H = int32(1)

			var c sdl.Color = getStarColor(stars[i].speed)
			starSettings.Renderer.SetDrawColor(c.R, c.G, c.B, 1)
			starSettings.Renderer.DrawPoint(rect.X, rect.Y)
			stars[i].X -= stars[i].speed

			if stars[i].X < 0 {
				stars[i].X = rand.Intn(starSettings.WindowWidth)
				stars[i].Y = rand.Intn(starSettings.WindowHeight)
				stars[i].speed = 1 + (rand.Intn(60) % 12)
			}
		}
	}
}

func getStarColor(speed int) sdl.Color {

	var color sdl.Color

	switch speed {

	case 1:
		color.R = 51
		color.G = 51
		color.B = 51
	case 2:
	case 3:
		color.R = 153
		color.G = 153
		color.B = 153
	case 4:
	case 5:
	case 6:
		color.R = 204
		color.G = 204
		color.B = 204
	default:
		color.R = 255
		color.G = 255
		color.B = 255
	}
	return color
}
