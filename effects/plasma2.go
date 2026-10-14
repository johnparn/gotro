package effects

// Ported from plasma.cpp: https://github.com/johangardhage/retro-demoeffects

import (
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const PLASMA_FRAMES int = 720
const PLASMA_COLORS uint8 = 255

type Plasma2Settings struct {
	WindowWidth    int
	WindowHeight   int
	Renderer       *sdl.Renderer
	Texture        *sdl.Texture
	PixelBuffer    []byte
	XValues        []float64
	SineValues     float64
	SineTable      []float64
	Colors         [256]RGBColor
	FrameRateLimit int
	LastFrameTime  time.Time
}

var framecounter float64
var plasma2Settings = Plasma2Settings{}
var frame int

// SetPlasma2FrameRateLimit sets the maximum render rate; fps <= 0 disables the limit.
func SetPlasma2FrameRateLimit(fps int) {
	plasma2Settings.FrameRateLimit = fps
	plasma2Settings.LastFrameTime = time.Time{}
}

func RunPlasma2() {
	if plasma2Settings.Texture == nil {
		return
	}
	if fps := plasma2Settings.FrameRateLimit; fps > 0 {
		frameDuration := time.Second / time.Duration(fps)
		if !plasma2Settings.LastFrameTime.IsZero() {
			if remaining := frameDuration - time.Since(plasma2Settings.LastFrameTime); remaining > 0 {
				time.Sleep(remaining)
			}
		}
		plasma2Settings.LastFrameTime = time.Now()
	}

	// Calculate frame
	framecounter += 2
	frame = int(math.Mod(framecounter, float64(PLASMA_FRAMES)))

	width := plasma2Settings.WindowWidth
	height := plasma2Settings.WindowHeight
	for x := 0; x < width; x++ {
		plasma2Settings.XValues[x] = 75 + plasma2Settings.SineTable[x*2+frame/2] +
			plasma2Settings.SineTable[x+frame*2] + plasma2Settings.SineTable[x/2+frame]*2
	}

	for y := 0; y < height; y++ {
		var yc float64 = 75 + plasma2Settings.SineTable[y+frame*2]*2 + plasma2Settings.SineTable[y*2+frame/2] + plasma2Settings.SineTable[y+frame]*2
		rowOffset := y * width * 4
		for x, xc := range plasma2Settings.XValues {
			color := plasma2Settings.Colors[uint8(xc*yc)]
			pixelOffset := rowOffset + x*4
			plasma2Settings.PixelBuffer[pixelOffset] = color.R.Value
			plasma2Settings.PixelBuffer[pixelOffset+1] = color.G.Value
			plasma2Settings.PixelBuffer[pixelOffset+2] = color.B.Value
			plasma2Settings.PixelBuffer[pixelOffset+3] = color.A.Value
		}
	}

	if err := plasma2Settings.Texture.Update(nil, unsafe.Pointer(&plasma2Settings.PixelBuffer[0]), width*4); err != nil {
		fmt.Println("Could not update plasma texture:", err)
		return
	}
	if err := plasma2Settings.Renderer.Copy(plasma2Settings.Texture, nil, nil); err != nil {
		fmt.Println("Could not copy plasma texture:", err)
	}
}

func (p *Plasma2Settings) setColorBuffer(colorIndex, r, g, b, a uint8) {
	p.Colors[colorIndex] = RGBColor{R: Color{Value: r}, G: Color{Value: g}, B: Color{Value: b}, A: Color{Value: a}}
}

func InitializePlasma2(renderer *sdl.Renderer, windowWidth int, windowHeight int) {
	plasma2Settings.Renderer = renderer
	plasma2Settings.WindowWidth = windowWidth
	plasma2Settings.WindowHeight = windowHeight
	plasma2Settings.PixelBuffer = make([]byte, windowWidth*windowHeight*4)
	plasma2Settings.XValues = make([]float64, windowWidth)
	plasma2Settings.SineValues = float64(windowWidth + PLASMA_FRAMES*2)

	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(windowWidth), int32(windowHeight))
	if err != nil {
		fmt.Println("Could not create plasma texture:", err)
		return
	}
	plasma2Settings.Texture = texture

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
