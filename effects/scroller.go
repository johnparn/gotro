package effects

import (
	"fmt"
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
	CharHeight  int  // Height of inidividual char on sprite map
	CharWidth   int  // Width of inidividual char on sprite map
	CharGrow    bool // If char is growing or not
	Flipped     bool // Whether scroll message is flipped vertically or not
	RotateChars bool // Rotate chars or not
}

var scrollerSettings = scrollerSettingsStruct{}

type commandStruct struct {
	start   int
	end     int
	len     int
	command string
	value   int
}

var msg []byte = []byte("welcome!             this is a text for testing purposes only               ")

// Setup scroller settings
func InitScroller(renderer *sdl.Renderer, windowWidth int, windowHeight int, rotateChars bool) {

	// msgs = append([]byte("welcome!               "))
	// msgs = append([]byte("this is a text for testing purposes only               "))
	// msgs = append([]byte("the demo is for just for trying out go as a programming language, powered by sdl for handling the graphics and sound.                                                        "))

	scrollerSettings.Renderer = renderer
	scrollerSettings.RotateChars = rotateChars

	_, errTexture := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, int32(windowWidth), int32(windowHeight))
	if errTexture != nil {
		fmt.Println("Could not create texture")
		os.Exit(1)
	}

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

	scrollerSettings.CharWidth = CHAR_WIDTH
	scrollerSettings.CharHeight = CHAR_HEIGHT
	scrollerSettings.CharGrow = false
	scrollerSettings.Flipped = false
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

func getControlCharIndex(character string) int {
	var charIndex int = -1
	for i := 0; i < len(ControlCharsMap)-1; i++ {
		if character == string(ControlCharsMap[i]) {
			charIndex = i
			break
		}
	}
	return charIndex
}

var doPause bool = false
var pauseEndTime uint64

func pauseScroller(currTime uint64, duration uint64) {
	fmt.Println("Pause", currTime, duration)
	doPause = true
	pauseEndTime = currTime + duration
}

func unpauseScroller() {
	doPause = false
}

// Sinewave
// var amp float64 = 0.0
// var frequency float64 = (math.Pi / FPS)
// var angle float64 = math.Pi
// var running bool = true
// sine.Angle += sine.Frequency / 2

// var angle := sine.Angle
// for i = 0; i < int32(windowWidth); i++ {
// 	y := math.Sin(angle)*sine.Amplitude*(float64(windowHeight/2)*0.5) + float64(sine.OffsetX) + float64(windowHeight/2)
// 	renderer.SetDrawColor(rColor.Value, gColor.Value, bColor.Value, alpha.Value)
// 	renderer.DrawPoint(i, int32(y))
// 	angle -= (sine.Frequency / 3)
// }

func RunScroller() {
	var dest, src sdl.Rect

	for i := range msg {

		v := getChar(string(msg[i]))

		// In case there is no matching character then add space
		if v.X < 0 && v.Y < 0 {
			scrollerSettings.CharOffset += CHAR_WIDTH
			continue
		}
		// Add space if it is a space
		if string(msg[i]) == " " {
			scrollerSettings.CharOffset += CHAR_WIDTH
			continue
		}

		src.X = int32(v.X)
		src.Y = int32(v.Y)
		src.W = int32(CHAR_WIDTH)
		src.H = int32(CHAR_HEIGHT)

		// Size of destination rect
		dest.X = int32(scrollerSettings.CharOffset)
		dest.Y = int32(scrollerSettings.WindowHeight)/2 - int32(CHAR_HEIGHT/2)
		dest.W = int32(scrollerSettings.CharWidth)
		dest.H = int32(scrollerSettings.CharHeight)

		// If text is to be "rotated", ie shrink/grow and flipped vertically or not
		if scrollerSettings.RotateChars && scrollerSettings.Flipped {
			dest.Y = int32(scrollerSettings.WindowHeight)/2 - int32(scrollerSettings.CharHeight/2)
			// Flip chars vertically
			scrollerSettings.Renderer.CopyEx(scrollerSettings.CharMap, &src, &dest, 0.0, &sdl.Point{X: int32(scrollerSettings.WindowWidth / 2), Y: int32(scrollerSettings.WindowHeight / 2)}, sdl.FLIP_VERTICAL)

		} else {
			dest.Y = int32(scrollerSettings.WindowHeight)/2 - int32(scrollerSettings.CharHeight/2)
			scrollerSettings.Renderer.Copy(scrollerSettings.CharMap, &src, &dest)
		}

		scrollerSettings.CharOffset += CHAR_WIDTH

		// Do not handle characters outside window width
		if scrollerSettings.CharOffset > scrollerSettings.WindowWidth {
			break
		}

	}

	// At end of message reset
	if scrollerSettings.D > (len(msg)*CHAR_WIDTH + scrollerSettings.WindowWidth) {
		scrollerSettings.D = 0
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
