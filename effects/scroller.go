package effects

import (
	"fmt"
	"math"
	"os"

	//"slices"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const CHAR_WIDTH int = 48
const CHAR_HEIGHT int = 50

var Map [7]string = [7]string{"abcdef", "ghijkl", "mnopqr", "stuvwx", "yz0123", "456789", "?!().,"}
var ControlCharsMap []string = []string{"|", "^"}

type scrollerSettingsStruct struct {
	Renderer *sdl.Renderer // sdl renderer instance
	// startTime    uint64
	// endTime      uint64
	CharMap      *sdl.Texture // sdl texture instance
	WindowWidth  int          // Width of window
	WindowHeight int          // Height of window
	D            int          // Scroll speed, ie Pixels to move scroll
	CharOffset   int
	// Sin          [360]float64
	CharHeight     int  // Height of inidividual char on sprite map
	CharWidth      int  // Width of inidividual char on sprite map
	CharGrow       bool // If char is growing or not
	Flipped        bool // Whether scroll message is flipped vertically or not
	RotateChars    bool // Rotate chars or not
	Message        []byte
	WaveEnabled    bool
	WaveMode       string
	WaveAmplitude  int
	WaveFrequency  float64
	WavePhase      float64
	WaveBaseOffset int
	WaveStartTime  uint64
	WaveSpeed      float64
}

var scrollerSettings = scrollerSettingsStruct{}

func calculateScrollerWaveY(xPos int, baseY int, currentTime uint64, settings scrollerSettingsStruct) int {
	if !settings.WaveEnabled {
		return baseY
	}

	elapsed := float64(currentTime-settings.WaveStartTime) * settings.WaveSpeed
	phase := settings.WavePhase + elapsed

	var wave float64
	switch settings.WaveMode {
	case "cosine":
		wave = math.Cos(float64(xPos)*settings.WaveFrequency + phase)
	default:
		wave = math.Sin(float64(xPos)*settings.WaveFrequency + phase)
	}

	return int(math.Round(float64(baseY+settings.WaveBaseOffset) + wave*float64(settings.WaveAmplitude)))
}

// SetScrollerWaveConfig enables or disables wave-based vertical movement for each character.
func SetScrollerWaveConfig(enabled bool, mode string, amplitude int, frequency float64, phase float64, baseOffset int) {
	if mode != "sine" && mode != "cosine" {
		mode = "sine"
	}

	scrollerSettings.WaveEnabled = enabled
	scrollerSettings.WaveMode = mode
	scrollerSettings.WaveAmplitude = amplitude
	scrollerSettings.WaveFrequency = frequency
	scrollerSettings.WavePhase = phase
	scrollerSettings.WaveBaseOffset = baseOffset
	scrollerSettings.WaveStartTime = sdl.GetTicks64()
	scrollerSettings.WaveSpeed = 0.0075
}

// Setup scroller settings
func InitScroller(renderer *sdl.Renderer, windowWidth int, windowHeight int, msg []byte, rotateChars bool) {

	scrollerSettings.Renderer = renderer
	scrollerSettings.RotateChars = rotateChars
	scrollerSettings.Message = msg

	// _, errTexture := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, int32(windowWidth), int32(windowHeight))
	// if errTexture != nil {
	// 	fmt.Println("Could not create texture")
	// 	os.Exit(1)
	// }

	//load image into surface
	tmpImg, errTmpImg := img.Load("260.bmp")
	if errTmpImg != nil {
		fmt.Println("Could not load image")
		os.Exit(1)
	}

	// Set color as transparent
	keycolor := sdl.MapRGB(tmpImg.Format, 0, 0, 0)
	tmpImg.SetColorKey(true, keycolor)

	charMap, errCharMap := tmpImg.ConvertFormat(sdl.PIXELFORMAT_RGBA8888, 0)
	if errCharMap != nil {
		fmt.Println("Could not convert  image format")
		os.Exit(1)
	}

	tmpImg.Free()

	cMap, errCMapTexture := renderer.CreateTextureFromSurface(charMap)
	if errCMapTexture != nil {
		fmt.Println("Could not convert  image format")
		os.Exit(1)
	}

	scrollerSettings.CharMap = cMap
	scrollerSettings.WindowWidth = windowWidth
	scrollerSettings.WindowHeight = windowHeight
	scrollerSettings.D = 1
	scrollerSettings.CharOffset = 0
	// scrollerSettings.Sin = sin_t
	scrollerSettings.CharWidth = CHAR_WIDTH
	scrollerSettings.CharHeight = CHAR_HEIGHT
	scrollerSettings.CharGrow = false
	scrollerSettings.Flipped = false
	scrollerSettings.WaveEnabled = true
	scrollerSettings.WaveMode = "sine"
	scrollerSettings.WaveAmplitude = 64
	scrollerSettings.WaveFrequency = 8
	scrollerSettings.WavePhase = 0.0
	scrollerSettings.WaveBaseOffset = 0
	scrollerSettings.WaveStartTime = sdl.GetTicks64()
	scrollerSettings.WaveSpeed = 0.0075

	doPause = false
	pauseEndTime = 0
	nextPause = 0
}

// func parseText(text string)

type Vector3d struct {
	X int
	Y int
	Z int
}

// Get bitmap font matching char
func getChar(c string) Vector3d {
	var pos Vector3d = Vector3d{X: 0, Y: 0, Z: 0}

	for i := 0; i < len(Map); i++ {

		for j := 0; j < len(Map[i]); j++ {

			if c == string(Map[i][j]) {
				return pos
			}
			pos.X += CHAR_WIDTH
		}

		pos.Y += CHAR_HEIGHT
		pos.X = 0
	}

	pos.X = -1
	pos.Y = -1

	return pos
}

var doPause bool = false
var pauseEndTime uint64
var pauseSchedule []ScrollerPause
var nextPause int

// ScrollerPause describes a pause relative to the message scroll offset.
type ScrollerPause struct {
	Offset   int
	Duration uint64
}

// SetScrollerPauses configures one-shot pauses in ascending offset order.
func SetScrollerPauses(pauses []ScrollerPause) {
	pauseSchedule = append([]ScrollerPause(nil), pauses...)
	nextPause = 0
}

// PauseScroller pauses movement for the specified number of seconds.
func PauseScroller(seconds uint64) {
	doPause = true
	pauseEndTime = sdl.GetTicks64() + seconds*1000
}

func unpauseScroller() {
	doPause = false
}

func RunScroller() {
	var dest, src sdl.Rect

	currTime := sdl.GetTicks64()
	if doPause && currTime >= pauseEndTime {
		unpauseScroller()
	}
	if !doPause && nextPause < len(pauseSchedule) && scrollerSettings.D >= pauseSchedule[nextPause].Offset {
		PauseScroller(pauseSchedule[nextPause].Duration)
		nextPause++
	}

	for i := range scrollerSettings.Message {

		v := getChar(string(scrollerSettings.Message[i]))

		// In case there is no matching character then add space
		if v.X < 0 && v.Y < 0 {
			scrollerSettings.CharOffset += CHAR_WIDTH
			continue
		}
		// Add space if it is a space
		if string(scrollerSettings.Message[i]) == " " {
			scrollerSettings.CharOffset += CHAR_WIDTH
			continue
		}

		src.X = int32(v.X)
		src.Y = int32(v.Y)
		src.W = int32(CHAR_WIDTH)
		src.H = int32(CHAR_HEIGHT)

		// Size of destination rect
		dest.X = int32(scrollerSettings.CharOffset)
		baseY := scrollerSettings.WindowHeight/2 - CHAR_HEIGHT/2
		dest.Y = int32(calculateScrollerWaveY(scrollerSettings.CharOffset, baseY, currTime, scrollerSettings))
		dest.W = int32(scrollerSettings.CharWidth)
		dest.H = int32(scrollerSettings.CharHeight)

		// SDL_RendererFlip flip = SDL_FLIP_HORIZONTAL | SDL_FLIP_VERTICAL;
		// SDL_RenderCopyEx(renderer, texture, &srcrect, &dstrect, angle, &center, flip);

		// If text is to be "rotated", ie shrink/grow and flipped vertically or not
		if scrollerSettings.RotateChars && scrollerSettings.Flipped {
			// Flip chars vertically
			scrollerSettings.Renderer.CopyEx(scrollerSettings.CharMap, &src, &dest, 0.0, &sdl.Point{X: int32(scrollerSettings.WindowWidth / 2), Y: int32(scrollerSettings.WindowHeight / 2)}, sdl.FLIP_VERTICAL)

		} else {
			scrollerSettings.Renderer.Copy(scrollerSettings.CharMap, &src, &dest)
		}

		scrollerSettings.CharOffset += CHAR_WIDTH

		// Do not handle characters outside window width
		if scrollerSettings.CharOffset > scrollerSettings.WindowWidth {
			break
		}

	}

	// At end of message reset
	if scrollerSettings.D > (len(scrollerSettings.Message)*CHAR_WIDTH + scrollerSettings.WindowWidth) {
		// var initTime = sdl.GetTicks64()
		scrollerSettings.D = 0
		nextPause = 0
	}

	// Offset for message
	scrollerSettings.CharOffset = scrollerSettings.WindowWidth - scrollerSettings.D
	if !doPause {
		// Speed of scroll
		scrollerSettings.D += 4
	}

	// Update rotate settings
	if scrollerSettings.RotateChars {
		// Char direction grow/shrink
		if scrollerSettings.CharHeight >= CHAR_HEIGHT {
			scrollerSettings.CharGrow = false
		} else if scrollerSettings.CharHeight <= 1 {
			scrollerSettings.CharGrow = true
			scrollerSettings.Flipped = !scrollerSettings.Flipped
		}

		// Make char grow by changing CharHeight that is used for the destination rect
		if scrollerSettings.CharGrow {
			scrollerSettings.CharHeight += 1
		} else if !scrollerSettings.CharGrow {
			scrollerSettings.CharHeight -= 1
		}
	}
}
