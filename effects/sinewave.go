package effects

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Sine struct {
	AmpDown   bool
	Step      float64
	Amplitude float64
	Threshold float64
	Frequency float64
	Angle     float64
	RGBColor  RGBColor
	OffsetX   int
	OffsetY   int
}

func SineWave(renderer *sdl.Renderer, sine *Sine, windowWidth int, windowHeight int) {

	drawSine(renderer, sine, windowWidth, windowHeight)
	SetAmplitude(sine)
	sine.Angle += sine.Frequency / 2

}

func SetAmplitude(sine *Sine) {
	if sine.Amplitude >= sine.Threshold {
		sine.AmpDown = true
	} else if sine.Amplitude <= -sine.Threshold {
		sine.AmpDown = false
	}

	if sine.AmpDown {
		sine.Amplitude -= sine.Step

	} else if !sine.AmpDown {
		sine.Amplitude += sine.Step
	}
}

func drawSine(renderer *sdl.Renderer, sine *Sine, windowWidth int, windowHeight int) {

	// Copy reference values as we don't want to change the reference here
	angle := sine.Angle
	rColor := sine.RGBColor.R
	gColor := sine.RGBColor.G
	bColor := sine.RGBColor.B
	alpha := sine.RGBColor.A

	var i int32
	// renderer.SetDrawBlendMode(sdl.BLENDMODE_ADD)

	for i = 0; i < int32(windowWidth); i++ {
		y := math.Sin(angle)*sine.Amplitude*(float64(windowHeight/2)*0.5) + float64(sine.OffsetX) + float64(windowHeight/2)

		rColor = GetColor(rColor)
		gColor = GetColor(gColor)
		bColor = GetColor(bColor)

		renderer.SetDrawColor(rColor.Value, gColor.Value, bColor.Value, alpha.Value)
		renderer.DrawPoint(i, int32(y))

		angle -= (sine.Frequency / 3)
	}
}
