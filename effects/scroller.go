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
	LastFrameMs  uint64
	// Sin          [360]float64
	CharHeight         int  // Height of inidividual char on sprite map
	InitialCharHeight  int  // Initial value of CharHeight for reset
	CharWidth          int  // Width of inidividual char on sprite map
	CharGrow           bool // If char is growing or not
	Flipped            bool // Whether scroll message is flipped vertically or not
	RotateChars        bool // Rotate chars or not
	RotationFadingOut  bool // Whether rotation is collapsing before being disabled
	InitialRotateChars bool // Initial value of RotateChars for reset
	Message            []byte
	WaveEnabled        bool
	InitialWaveEnabled bool // Initial value of WaveEnabled for reset
	WaveMode           string
	WaveAmplitude      int
	WaveFrequency      float64
	WavePhase          float64
	WaveBaseOffset     int
	WaveStartTime      uint64
	WaveSpeed          float64
	WaveFadeInTime     uint64 // Time in milliseconds for wave to fade in
	WaveFadeOutTime    uint64 // Time in milliseconds for wave to fade out
	WaveFadeOutStart   uint64 // When to start fading out (0 means no fade out)
	RotationDelay      int    // Delay between character rotations in pixels
	RotationCycle      int    // Total pixels for a full rotation cycle
}

var scrollerSettings = scrollerSettingsStruct{}

const defaultScrollerPixelsPerSecond = 240.0

func nextScrollerOffset(currentOffset float64, frameMilliseconds float64, paused bool) float64 {
	if paused {
		return currentOffset
	}

	return currentOffset + (defaultScrollerPixelsPerSecond * frameMilliseconds / 1000.0)
}

func calculateScrollerWaveY(xPos int, baseY int, currentTime uint64, settings scrollerSettingsStruct) int {
	// Allow fade-out to complete even if wave was disabled
	if !settings.WaveEnabled && settings.WaveFadeOutStart == 0 {
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

	// Calculate fade envelope (0.0 to 1.0)
	fadeEnvelope := 1.0
	timeSinceStart := currentTime - settings.WaveStartTime

	// Fade in (only if not fading out)
	if settings.WaveFadeOutStart == 0 && settings.WaveFadeInTime > 0 && timeSinceStart < settings.WaveFadeInTime {
		fadeEnvelope = float64(timeSinceStart) / float64(settings.WaveFadeInTime)
	}

	// Fade out
	if settings.WaveFadeOutStart > 0 && currentTime >= settings.WaveFadeOutStart {
		timeSinceFadeOut := currentTime - settings.WaveFadeOutStart
		if timeSinceFadeOut < settings.WaveFadeOutTime {
			fadeEnvelope = 1.0 - (float64(timeSinceFadeOut) / float64(settings.WaveFadeOutTime))
		} else {
			// Fade out complete
			fadeEnvelope = 0.0
		}
	}

	return int(math.Round(float64(baseY+settings.WaveBaseOffset) + wave*fadeEnvelope*float64(settings.WaveAmplitude)))
}

// calculateCharHeight calculates the height for a character based on its position in the scroll
func calculateCharHeight(charPos int, settings scrollerSettingsStruct) int {
	if !settings.RotateChars {
		return settings.CharHeight
	}

	// Calculate the delay offset for this character based on its position
	// Each character starts rotating after the previous one by RotationDelay pixels
	adjustedPos := charPos - settings.RotationDelay
	if adjustedPos < 0 {
		adjustedPos = 0
	}

	// Position within the rotation cycle (0 to RotationCycle)
	posInCycle := adjustedPos % settings.RotationCycle

	// First half: grow from 1 to CHAR_HEIGHT
	// Second half: shrink from CHAR_HEIGHT to 1
	if posInCycle < settings.RotationCycle/2 {
		// Growing phase
		return 1 + posInCycle
	} else {
		// Shrinking phase
		return CHAR_HEIGHT - (posInCycle - settings.RotationCycle/2)
	}
}

func setScrollerRotateChars(enabled bool) {
	if enabled {
		scrollerSettings.RotationFadingOut = false
		scrollerSettings.RotateChars = true
		return
	}
	if scrollerSettings.RotateChars {
		// Keep the rotation state active until the collapse has rendered.
		scrollerSettings.RotationFadingOut = true
		return
	}
	scrollerSettings.RotateChars = enabled
}

func updateScrollerRotation(settings *scrollerSettingsStruct) {
	if settings.RotationFadingOut {
		if !settings.CharGrow && settings.CharHeight > 1 {
			settings.CharHeight--
			return
		}

		if settings.CharHeight < settings.InitialCharHeight {
			settings.Flipped = false
			settings.CharGrow = true
			settings.CharHeight++
			return
		}

		settings.Flipped = false
		settings.RotationFadingOut = false
		settings.RotateChars = false
		return
	}
	if !settings.RotateChars {
		return
	}

	// Char direction grow/shrink
	if settings.CharHeight >= settings.InitialCharHeight {
		settings.CharGrow = false
	} else if settings.CharHeight <= 1 {
		settings.CharGrow = true
		settings.Flipped = !settings.Flipped
	}

	// Make char grow by changing CharHeight that is used for the destination rect
	if settings.CharGrow {
		settings.CharHeight++
	} else {
		settings.CharHeight--
	}
}

// SetScrollerWaveConfig enables or disables wave-based vertical movement for each character.
// fadeInMs: time in milliseconds for wave amplitude to fade in (0 for instant)
// fadeOutMs: time in milliseconds for wave amplitude to fade out (0 for no fade out)
func SetScrollerWaveConfig(enabled bool, mode string, amplitude int, frequency float64, phase float64, baseOffset int, fadeInMs uint64, fadeOutMs uint64) {
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
	scrollerSettings.WaveFadeInTime = fadeInMs
	scrollerSettings.WaveFadeOutTime = fadeOutMs
	scrollerSettings.WaveFadeOutStart = 0 // Only set when you want to trigger fade out
}

// StartScrollerWaveFadeOut triggers the wave to start fading out
func StartScrollerWaveFadeOut() {
	scrollerSettings.WaveFadeOutStart = sdl.GetTicks64()
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
	scrollerSettings.LastFrameMs = uint64(sdl.GetTicks64())
	// scrollerSettings.Sin = sin_t
	scrollerSettings.CharWidth = CHAR_WIDTH
	scrollerSettings.InitialCharHeight = CHAR_HEIGHT
	scrollerSettings.CharHeight = CHAR_HEIGHT
	scrollerSettings.CharGrow = false
	scrollerSettings.Flipped = false
	scrollerSettings.RotationFadingOut = false
	scrollerSettings.WaveEnabled = false
	scrollerSettings.InitialWaveEnabled = false
	scrollerSettings.RotateChars = rotateChars
	scrollerSettings.InitialRotateChars = rotateChars
	scrollerSettings.WaveMode = "sine"
	scrollerSettings.WaveAmplitude = 64
	scrollerSettings.WaveFrequency = 8
	scrollerSettings.WavePhase = 0.0
	scrollerSettings.WaveBaseOffset = 0
	scrollerSettings.WaveStartTime = sdl.GetTicks64()
	scrollerSettings.WaveSpeed = 0.0075
	scrollerSettings.WaveFadeInTime = 0
	scrollerSettings.WaveFadeOutTime = 0
	scrollerSettings.WaveFadeOutStart = 0
	scrollerSettings.RotationDelay = CHAR_WIDTH      // Delay between characters
	scrollerSettings.RotationCycle = CHAR_HEIGHT * 2 // Full cycle: grow and shrink

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

// Unpause
type ScrollerUnpause struct {
	RotateChars bool
	WaveEnabled bool
}

// ScrollerPause describes a pause relative to the message scroll offset.
type ScrollerPause struct {
	Offset      int
	Duration    uint64
	RotateChars bool
	WaveEnabled bool
	Unpause     ScrollerUnpause
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
		setScrollerRotateChars(pauseSchedule[nextPause].Unpause.RotateChars)
		if scrollerSettings.WaveEnabled && !pauseSchedule[nextPause].Unpause.WaveEnabled {
			// Wave is being disabled, reset fade-out time for proper fade-out
			scrollerSettings.WaveFadeOutStart = currTime
		}
		if pauseSchedule[nextPause].Unpause.WaveEnabled {
			scrollerSettings.WaveFadeOutStart = 0
			scrollerSettings.WaveStartTime = currTime
		}
		scrollerSettings.WaveEnabled = pauseSchedule[nextPause].Unpause.WaveEnabled
		nextPause++
	}
	if !doPause && nextPause < len(pauseSchedule) && scrollerSettings.D >= pauseSchedule[nextPause].Offset {
		PauseScroller(pauseSchedule[nextPause].Duration)
		setScrollerRotateChars(pauseSchedule[nextPause].RotateChars)
		if pauseSchedule[nextPause].WaveEnabled && !scrollerSettings.WaveEnabled {
			// Wave is being enabled, reset start time for proper fade-in
			scrollerSettings.WaveStartTime = currTime
			scrollerSettings.WaveFadeOutStart = 0
		}
		scrollerSettings.WaveEnabled = pauseSchedule[nextPause].WaveEnabled
	}

	frameMilliseconds := 0.0
	if scrollerSettings.LastFrameMs != 0 {
		frameMilliseconds = float64(currTime - scrollerSettings.LastFrameMs)
	}
	scrollerSettings.LastFrameMs = currTime
	if !doPause {
		scrollerSettings.D = int(math.Round(nextScrollerOffset(float64(scrollerSettings.D), frameMilliseconds, false)))
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
		if (scrollerSettings.RotateChars || scrollerSettings.RotationFadingOut) && scrollerSettings.Flipped {
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
		// Reset wave and rotation to initial values
		scrollerSettings.WaveEnabled = scrollerSettings.InitialWaveEnabled
		setScrollerRotateChars(scrollerSettings.InitialRotateChars)
		scrollerSettings.CharHeight = scrollerSettings.InitialCharHeight
		scrollerSettings.Flipped = false
		scrollerSettings.RotationFadingOut = false
		scrollerSettings.WaveFadeOutStart = 0 // Clear any active fade-out
	}

	// Offset for message
	scrollerSettings.CharOffset = scrollerSettings.WindowWidth - scrollerSettings.D

	// Update rotate settings, including the collapse when rotation is disabled.
	updateScrollerRotation(&scrollerSettings)
}
