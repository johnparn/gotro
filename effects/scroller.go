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

// const CHAR_WIDTH int = 32
// const CHAR_HEIGHT int = 32

// var Map [6]string = [6]string{" !\"    '()", "  ,-. 0123", "456789:;  ", " ? abcdefg", "hijklmnopq", "rstuvwxyz "}

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

// the demo is for just for trying out go as a programming language, powered by sdl for handling the graphics and sound.                                                        ")

// var msgs = make([]byte, 0)

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

	// for i := 0; i < 360; i++ {
	// 	sin_t[i] = math.Sin(float64(i) * math.Pi / 180.0)
	// }

	// pattern := regexp.MustCompile(`\|p\d+\|`)
	// matches := pattern.FindAllIndex(msg, -1)
	// for _, match := range matches {
	// 	match_substr := msg[match[0]:match[1]]

	// 	clean_substr := bytes.ReplaceAll(match_substr, []byte("|"), []byte(""))

	// 	if string(clean_substr[0:1]) == "p" {
	// 		scrollControl = append(scrollControl, commandStruct{command: "pause", value: len(clean_substr[1:len(clean_substr)]), start: match[0], end: match[1], len: len(match_substr)})
	// 		fmt.Println(scrollControl)
	// 	}
	// }
	// msg = pattern.ReplaceAll(msg, []byte(""))

	// scrollerSettings.startTime = sdl.GetTicks64() + 2400
	// scrollerSettings.endTime = scrollerSettings.startTime + 5000

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

	// var currTime uint64 = sdl.GetTicks64()

	// if doPause == true && currTime >= pauseEndTime {
	// 	unpauseScroller()
	// }

	// // // fmt.Println("Time", scrollerSettings.startTime, scrollerSettings.endTime, currTime)
	// if doPause == false && scrollerSettings.CharOffset == (scrollerSettings.WindowWidth/2)-(scrollerSettings.CharOffset/2)-CHAR_WIDTH {
	// 	fmt.Println(currTime, pauseEndTime, scrollerSettings.CharOffset)
	// 	pauseScroller(currTime, 2000)
	// }

	for i := range msg {

		// if scrollerSettings.CharOffset == i*CHAR_WIDTH && i == 10 {
		// 	fmt.Println("Command position", scrollerSettings.CharOffset/CHAR_WIDTH)
		// 	doPause = true
		// }
		// ctrlChar := getControlCharIndex(string(msg[i]))
		// if ctrlChar >= 0 {
		// 	pauseScroller += sdl.GetTicks64() + 2000
		// 	fmt.Println("CHAR", string(msg[i]), ctrlChar)
		// }
		v := getChar(string(msg[i]))

		// if string(msg[i]) == "|" {
		// 	n := strings.Index(msg[i:], "|")

		// }

		// for _, v := range scrollControl {
		// 	if v.start >= scrollerSettings.CharOffset/CHAR_WIDTH {
		// 		if v.command == "pause" {
		// 			scrollerSettings.D = 0
		// 			// doPause = true
		// 			pauseUntil = sdl.GetTicks64() + 3000
		// 		}
		// 	}
		// }

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
		// dest.Y = int32(math.Sin(float64((math.Pi/60)/3))*0.0*(float64(scrollerSettings.WindowHeight/2)*0.5) + float64(scrollerSettings.WindowHeight/2))
		dest.W = int32(scrollerSettings.CharWidth)
		dest.H = int32(scrollerSettings.CharHeight)

		// SDL_RendererFlip flip = SDL_FLIP_HORIZONTAL | SDL_FLIP_VERTICAL;
		// SDL_RenderCopyEx(renderer, texture, &srcrect, &dstrect, angle, &center, flip);

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
		// var initTime = sdl.GetTicks64()
		scrollerSettings.D = 0
		// Change start anc end time according to new ticks
		// scrollerSettings.startTime += initTime
		// scrollerSettings.endTime += initTime
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
