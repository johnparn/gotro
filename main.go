package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/johnparn/gotro/effects"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS          = 60
	windowWidth  = 800
	windowHeight = 600
)

var window *sdl.Window
var renderer *sdl.Renderer

func main() {
	sdlInitVideo()
	sdlInitImage()
	sdlInitAudio()

	window = createWindow()
	renderer = createRenderer()
	var _ = sdl.PollEvent() //MacOS won't draw the window without this line

	playMusic()

	drawSceneScroll()

	drawScenePlasma()
	drawStars2()
	drawScenePlasma2()

	drawSceneSinus(sdl.GetTicks64())
	drawDotTunnel()
	drawDotSphere()
	drawTwister()
	drawBoingBall(192, 102, 0)

	_ = renderer.Destroy()
	_ = window.Destroy()

}

func sleepSeconds(secs time.Duration) {
	time.Sleep(time.Second * secs)
}

func sdlInitVideo() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialise video: %s\n", err)
		os.Exit(1)
	}

	defer sdl.Quit()
}

func sdlInitImage() {
	err := img.Init(img.INIT_PNG)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialise image lib: %s\n", err)
		os.Exit(1)
	}

	defer img.Quit()
}

func sdlInitAudio() {
	errSDLAudioInit := sdl.Init(sdl.INIT_AUDIO)
	if errSDLAudioInit != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialise audio: %s\n", errSDLAudioInit)
		os.Exit(1)
	}

	errOpeningAudioDevice := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096)
	if errOpeningAudioDevice != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to open audio device: %s\n", errOpeningAudioDevice)
		os.Exit(1)
	}

	errSDLMixerInit := mix.Init(mix.INIT_MOD)
	if errSDLMixerInit != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialise mixer: %s\n", errSDLMixerInit)
		os.Exit(1)
	}

}

func createWindow() *sdl.Window {
	window, errCreatingSDLWindow := sdl.CreateWindow("Intro",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		windowWidth,
		windowHeight,
		sdl.WINDOW_MAXIMIZED|sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)

	if errCreatingSDLWindow != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create SDL window: %s\n", errCreatingSDLWindow)
		os.Exit(1)
	}
	return window
}

func createRenderer() *sdl.Renderer {

	var numDrivers, _ = sdl.GetNumRenderDrivers()
	fmt.Println("Render drivers", numDrivers)

	var errCreatingSDLRenderer error
	sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1")
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	if runtime.GOOS == "darwin" {
		sdl.SetHint(sdl.HINT_RENDER_DRIVER, "software")
	} else {
		sdl.SetHint(sdl.HINT_RENDER_DRIVER, "opengl")
	}

	renderer, errCreatingSDLRenderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_TARGETTEXTURE)
	renderer.RenderSetVSync(true)

	for i := 0; i < numDrivers; i++ {
		driverInfo, _ := renderer.GetInfo()
		fmt.Println("Driver name (", i, "): ", driverInfo.Name)
	}

	var info sdl.RendererInfo
	info, _ = renderer.GetInfo()
	fmt.Println("RenderInfo", info.Name, info.RendererInfoData)

	if errCreatingSDLRenderer != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", errCreatingSDLRenderer)
		os.Exit(1)
	}
	return renderer
}

func playMusic() {
	if music, errLoadingMusic := mix.LoadMUS("echoing.mod"); errLoadingMusic != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load music from disk: %s\n", errLoadingMusic)
		os.Exit(1)
	} else {
		_ = music.Play(1)
	}
}

func drawDotTunnel() {

	var running bool = true

	effects.InitializeTunnel(renderer, windowWidth, windowHeight)
	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			case sdl.MOUSEBUTTONDOWN:
				running = false
			}
		}
		effects.RenderDotTunnel(float64(sdl.GetTicks64()))
		updateScreen()
	}
}

func drawDotSphere() {

	var running bool = true

	effects.InitializeDotSphere(renderer, windowWidth, windowHeight)
	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			case sdl.MOUSEBUTTONDOWN:
				running = false
			}
		}
		effects.RenderDotSphere(float64(sdl.GetTicks64()))
		updateScreen()
	}
}

func drawTwister() {

	var running bool = true

	effects.InitTwister(renderer, windowWidth, windowHeight)
	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}

		effects.RunTwister()
		updateScreen()
	}
}

func drawBoingBall(R, G, B uint8) {

	var running bool = true

	// backgroundFill(0, 0, 0)
	renderer.SetDrawColor(R, G, B, 255)
	renderer.DrawRect(&sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight})
	effects.InitBoingball(renderer, windowWidth, windowHeight, 128, 64, 32)
	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}
		effects.BoingBall()
		updateScreen()
	}
}

func drawSceneScroll() {

	effects.InitScroller(renderer, windowWidth, windowHeight, true)
	effects.InitializeStars2(renderer, windowWidth, windowHeight)

	// scrollSettings2 := rotateScroller.InitScroller(renderer, windowWidth, windowHeight, false)
	var running bool = true
	var rgbColors [2]effects.RGBColor
	// var startTick uint64 = sdl.GetTicks64()

	rgbColors[0] = effects.RGBColor{
		R: effects.Color{Value: uint8(255), Increment: -1},
		G: effects.Color{Value: uint8(0), Increment: 1},
		B: effects.Color{Value: uint8(255), Increment: -1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rgbColors[1] = effects.RGBColor{
		R: effects.Color{Value: uint8(51), Increment: 1},
		G: effects.Color{Value: uint8(153), Increment: 1},
		B: effects.Color{Value: uint8(201), Increment: -1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	var sine [2]effects.Sine
	sine[0] = effects.Sine{
		Amplitude: 0.0,
		Step:      0.0025,
		Threshold: 0.25,
		AmpDown:   true,
		OffsetX:   0,
		OffsetY:   0,
		Frequency: (math.Pi / FPS),
		Angle:     math.Pi,
		RGBColor:  rgbColors[0],
	}

	sine[1] = effects.Sine{
		Amplitude: 0.0,
		Step:      0.0025,
		Threshold: 0.30,
		AmpDown:   true,
		OffsetX:   0,
		OffsetY:   0,
		Frequency: (math.Pi / FPS),
		Angle:     math.Pi,
		RGBColor:  rgbColors[1],
	}

	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}
		// renderer.Clear()
		renderer.SetDrawColor(15, 15, 23, 1)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(windowWidth), H: int32(windowHeight)})
		// rasterBars()
		effects.RunStars2()
		effects.SineWave(renderer, &sine[0], windowWidth, windowHeight)
		effects.SineWave(renderer, &sine[1], windowWidth, windowHeight)
		effects.RunScroller()

		// if sdl.GetTicks64() >= startTick+40000 {
		// 	running = false
		// 	break
		// }
		updateScreen()

	}
}

func drawStars2() {

	//plasma := effects.InitPlasma(renderer, windowWidth, windowHeight)
	effects.InitializeStars2(renderer, windowWidth, windowHeight)

	var running bool = true

	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}

		effects.RunStars2()
		updateScreen()
	}
}

func drawScenePlasma2() {

	//plasma := effects.InitPlasma(renderer, windowWidth, windowHeight)
	effects.InitializePlasma2(renderer, int(windowWidth), int(windowHeight))

	var running bool = true

	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}

		effects.RunPlasma2()

		// Moving plasma
		// plasma = effects.RunPlasma(&plasma)

		updateScreen()
	}
}

func drawScenePlasma() {

	var startTick uint64 = sdl.GetTicks64()

	//plasma := effects.InitPlasma(renderer, windowWidth, windowHeight)
	effects.InitPlasma(renderer, windowWidth, windowHeight)

	var running bool = true

	for running {
		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}

		effects.RunPlasma()

		if sdl.GetTicks64() >= startTick+8000 {
			//running = false
			break
		}
		updateScreen()
	}
}

func drawSceneSinus(startTick uint64) {

	var amp float64 = 0.0
	var frequency float64 = (math.Pi / FPS)
	var angle float64 = math.Pi
	var running bool = true

	// // Define start colors for curves
	var rbgColors [8]effects.RGBColor

	rbgColors[0] = effects.RGBColor{
		R: effects.Color{Value: uint8(255), Increment: -1},
		G: effects.Color{Value: uint8(0), Increment: 1},
		B: effects.Color{Value: uint8(0), Increment: 1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rbgColors[1] = effects.RGBColor{
		R: effects.Color{Value: uint8(0), Increment: 1},
		G: effects.Color{Value: uint8(255), Increment: -1},
		B: effects.Color{Value: uint8(0), Increment: 1},
		A: effects.Color{Value: uint8(0), Increment: 1},
	}

	rbgColors[2] = effects.RGBColor{
		R: effects.Color{Value: uint8(0), Increment: 1},
		G: effects.Color{Value: uint8(0), Increment: 1},
		B: effects.Color{Value: uint8(255), Increment: -1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rbgColors[3] = effects.RGBColor{
		R: effects.Color{Value: uint8(255), Increment: -1},
		G: effects.Color{Value: uint8(255), Increment: -1},
		B: effects.Color{Value: uint8(0), Increment: 1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rbgColors[4] = effects.RGBColor{
		R: effects.Color{Value: uint8(0), Increment: 1},
		G: effects.Color{Value: uint8(255), Increment: -1},
		B: effects.Color{Value: uint8(255), Increment: -1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rbgColors[5] = effects.RGBColor{
		R: effects.Color{Value: uint8(255), Increment: -1},
		G: effects.Color{Value: uint8(0), Increment: 1},
		B: effects.Color{Value: uint8(255), Increment: -1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rbgColors[6] = effects.RGBColor{
		R: effects.Color{Value: uint8(51), Increment: 1},
		G: effects.Color{Value: uint8(153), Increment: 1},
		B: effects.Color{Value: uint8(201), Increment: -1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	rbgColors[7] = effects.RGBColor{
		R: effects.Color{Value: uint8(153), Increment: -1},
		G: effects.Color{Value: uint8(0), Increment: 1},
		B: effects.Color{Value: uint8(201), Increment: 1},
		A: effects.Color{Value: uint8(255), Increment: 1},
	}

	effects.ResetStars(renderer, windowWidth, windowHeight)

	var sine [6]effects.Sine
	sine[0] = effects.Sine{
		Amplitude: amp,
		Step:      0.0025,
		Threshold: 0.25,
		AmpDown:   true,
		OffsetX:   0,
		OffsetY:   0,
		Frequency: frequency,
		Angle:     angle,
		RGBColor:  rbgColors[1],
	}
	sine[1] = effects.Sine{
		Amplitude: amp,
		Step:      0.0025,
		Threshold: 0.30,
		AmpDown:   true,
		OffsetX:   10,
		OffsetY:   0,
		Frequency: frequency,
		Angle:     angle,
		RGBColor:  rbgColors[2],
	}
	sine[2] = effects.Sine{
		Amplitude: amp,
		Step:      0.0025,
		Threshold: 0.15,
		AmpDown:   true,
		OffsetX:   20,
		OffsetY:   0,
		Frequency: frequency,
		Angle:     angle,
		RGBColor:  rbgColors[5],
	}
	sine[3] = effects.Sine{
		Amplitude: amp,
		Step:      0.0020,
		Threshold: 0.25,
		AmpDown:   true,
		OffsetX:   40,
		OffsetY:   0,
		Frequency: frequency,
		Angle:     angle,
		RGBColor:  rbgColors[6],
	}
	sine[4] = effects.Sine{
		Amplitude: amp,
		Step:      0.0018,
		Threshold: 0.27,
		AmpDown:   true,
		OffsetX:   50,
		OffsetY:   0,
		Frequency: frequency,
		Angle:     angle,
		RGBColor:  rbgColors[7],
	}
	sine[5] = effects.Sine{
		Amplitude: amp,
		Step:      0.0017,
		Threshold: 0.32,
		AmpDown:   true,
		OffsetX:   60,
		OffsetY:   0,
		Frequency: frequency,
		Angle:     angle,
		RGBColor:  rbgColors[3],
	}

	effects.InitStars(renderer, windowWidth, windowHeight)

	for running {

		// Quit this scene
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.KEYUP:
				running = false
			}
		}

		// Clear background
		renderer.SetDrawColor(15, 15, 23, 1)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(windowWidth), H: int32(windowHeight)})

		// Stars
		effects.UpdateStars()

		// Sine waves
		ticks := sdl.GetTicks64()

		if ticks >= startTick+1000 {
			effects.SineWave(renderer, &sine[1], windowWidth, windowHeight)
		}
		if ticks >= 3000+startTick {
			effects.SineWave(renderer, &sine[0], windowWidth, windowHeight)
		}
		if ticks >= 4500+startTick {
			effects.SineWave(renderer, &sine[2], windowWidth, windowHeight)
		}
		if ticks >= 6000+startTick {
			effects.SineWave(renderer, &sine[3], windowWidth, windowHeight)
		}
		if ticks >= 7000+startTick {
			effects.SineWave(renderer, &sine[4], windowWidth, windowHeight)
		}
		if ticks >= 7500+startTick {
			effects.SineWave(renderer, &sine[5], windowWidth, windowHeight)
		}

		// Show cube
		if ticks >= 8000+startTick {
			effects.SpinningCube(renderer, windowWidth, windowHeight)
		}

		if sdl.GetTicks64() >= 24000+startTick {
			running = false
			break
		}
		updateScreen()

	}
}

var delayTime uint64 = 1000 / FPS
var startTime uint64 = 0
var renderTime uint64 = 0
var diffTime uint32 = 0

// var fpscap uint64 = 0
// var fpsticks uint64 = 0
// var fpscount int64 = 0

func updateScreen() {

	// Time taken to render frame
	//startTime = sdl.GetTicks64()
	renderer.Present()
	//renderTime = sdl.GetTicks64() - startTime

	// Clear screen
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// Show FPS
	// if fpsticks <= sdl.GetTicks64()-1000.0 {
	// 	fmt.Println("FPS: ", fpscount)
	// 	fpsticks = sdl.GetTicks64()
	// 	fpscount = 0
	// }
	// fpscount++
}
