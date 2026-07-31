package effects

//
// Inspired by rototunnel.cpp by Johan Gardhage

import (
	"fmt"
	"math"
	"os"
	"unsafe"

	//"slices"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type RotoTunnelSettings struct {
	WindowWidth   int
	WindowHeight  int
	Renderer      *sdl.Renderer
	Texture       *sdl.Texture
	distanceTable [][]int
	angleTable    [][]int
	imageName     string
	imageColor    []RGBColor
	imageWidth    int
	imageHeight   int
	frame         float64
	lastTick      uint64
	dx            float64
	dy            float64
	cameraSpeed   float64
	tunnelSpeed   float64
	pixelBuffer   []byte
	pitch         int
}

var rotoTunnelSettings = RotoTunnelSettings{
	imageName:   "flowers.bmp",
	frame:       0,
	cameraSpeed: 0.005,
	tunnelSpeed: 0.002,
	lastTick:    0,
}

func RunRotoTunnel() {
	now := sdl.GetTicks64()
	var delta float64
	if rotoTunnelSettings.lastTick == 0 {
		delta = 1.0 / 60.0
	} else {
		delta = float64(now-rotoTunnelSettings.lastTick) / 1000.0
	}
	rotoTunnelSettings.lastTick = now
	// Advance in logical frame steps so the motion stays steady without speeding up over time.
	frameStep := delta * 60.0
	rotoTunnelSettings.frame += frameStep

	// Match the original effect’s motion more closely with configurable drift and camera sway.
	sx := float64(rotoTunnelSettings.imageWidth) * rotoTunnelSettings.tunnelSpeed * rotoTunnelSettings.frame
	sy := float64(rotoTunnelSettings.imageHeight) * 0.0045 * rotoTunnelSettings.frame // Rotating
	// sy := float64(rotoTunnelSettings.imageHeight) * 0.01 * delta // Slow drift based on delta time

	cameraPhase := rotoTunnelSettings.frame * rotoTunnelSettings.cameraSpeed
	rotoTunnelSettings.dx = float64(rotoTunnelSettings.WindowWidth/2) + float64(rotoTunnelSettings.WindowWidth/4)*math.Cos(cameraPhase)
	rotoTunnelSettings.dy = float64(rotoTunnelSettings.WindowHeight/2) + float64(rotoTunnelSettings.WindowHeight/4)*math.Sin(cameraPhase*3.3)

	// rotoTunnelSettings.dx = float64(rotoTunnelSettings.WindowWidth/2) + float64(rotoTunnelSettings.WindowWidth/4)*math.Cos(rotoTunnelSettings.frame*rotoTunnelSettings.cameraSpeed)
	// rotoTunnelSettings.dy = float64(rotoTunnelSettings.WindowHeight/2) + float64(rotoTunnelSettings.WindowHeight/4)*math.Sin(rotoTunnelSettings.frame*rotoTunnelSettings.cameraSpeed)
	dx := rotoTunnelSettings.dx
	dy := rotoTunnelSettings.dy

	width := rotoTunnelSettings.WindowWidth
	height := rotoTunnelSettings.WindowHeight
	stride := rotoTunnelSettings.pitch
	buffer := rotoTunnelSettings.pixelBuffer
	imageWidth := rotoTunnelSettings.imageWidth
	imageHeight := rotoTunnelSettings.imageHeight
	imageColors := rotoTunnelSettings.imageColor
	distanceTable := rotoTunnelSettings.distanceTable
	angleTable := rotoTunnelSettings.angleTable

	wrap := func(value float64, size int) int {
		wrapped := int(math.Mod(value, float64(size)))
		if wrapped < 0 {
			wrapped += size
		}
		return wrapped
	}

	for y := range height {
		rowBase := y * stride
		for x := 0; x < width; x++ {
			baseX := x + int(dx)
			baseY := y + int(dy)
			tx := wrap(float64(distanceTable[baseX][baseY])+sx, imageWidth)
			ty := wrap(float64(angleTable[baseX][baseY])+sy, imageHeight)

			color := imageColors[ty*imageWidth+tx]
			idx := rowBase + x*4
			buffer[idx] = color.R.Value
			buffer[idx+1] = color.G.Value
			buffer[idx+2] = color.B.Value
			buffer[idx+3] = color.A.Value
		}
	}

	if err := rotoTunnelSettings.Texture.Update(nil, unsafe.Pointer(&buffer[0]), stride); err != nil {
		fmt.Println("Could not update roto tunnel texture:", err)
	}
	if err := rotoTunnelSettings.Renderer.Copy(rotoTunnelSettings.Texture, nil, nil); err != nil {
		fmt.Println("Could not copy roto tunnel texture:", err)
	}
}

func InitializeRotoTunnel(renderer *sdl.Renderer, windowWidth int, windowHeight int, imageName string) {
	rotoTunnelSettings.Renderer = renderer
	rotoTunnelSettings.WindowWidth = windowWidth
	rotoTunnelSettings.WindowHeight = windowHeight
	if imageName != "" {
		rotoTunnelSettings.imageName = imageName
	}
	rotoTunnelSettings.angleTable = Make2D[int](rotoTunnelSettings.WindowWidth*2, rotoTunnelSettings.WindowHeight*2)
	rotoTunnelSettings.distanceTable = Make2D[int](rotoTunnelSettings.WindowWidth*2, rotoTunnelSettings.WindowHeight*2)

	texture, errTexture := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, int32(windowWidth), int32(windowHeight))
	if errTexture != nil {
		fmt.Println("Could not create roto tunnel texture:", errTexture)
		os.Exit(1)
	}
	rotoTunnelSettings.Texture = texture
	rotoTunnelSettings.pitch = windowWidth * 4
	rotoTunnelSettings.pixelBuffer = make([]byte, windowWidth*windowHeight*4)

	//load image into surface
	image, errTmpImg := img.Load(rotoTunnelSettings.imageName)
	if errTmpImg != nil {
		fmt.Println("Could not load image")
		os.Exit(1)
	}

	// Set color as transparent
	keycolor := sdl.MapRGB(image.Format, 0, 0, 0)
	image.SetColorKey(true, keycolor)

	image, errCharMap := image.ConvertFormat(sdl.PIXELFORMAT_RGBA8888, 0)
	if errCharMap != nil {
		fmt.Println("Could not convert  image format")
		os.Exit(1)
	}

	pixels := image.Pixels()
	fmt.Printf("%d bytes per pixel, in total %d pixels", len(pixels)/image.PixelNum(), image.PixelNum())

	// Get image width and height.
	rotoTunnelSettings.imageWidth = int(image.W)
	rotoTunnelSettings.imageHeight = int(image.H)

	var ratio = 128 * rotoTunnelSettings.imageHeight

	rotoTunnelSettings.imageColor = make([]RGBColor, rotoTunnelSettings.imageWidth*rotoTunnelSettings.imageHeight+rotoTunnelSettings.imageWidth) // + imageWidth to avoid index out of range error

	pitch := int(image.Pitch)
	for y := 0; y < rotoTunnelSettings.imageHeight; y++ {
		rowStart := y * pitch
		for x := 0; x < rotoTunnelSettings.imageWidth; x++ {
			idx := rowStart + x*4
			pixIdx := y*rotoTunnelSettings.imageWidth + x
			rotoTunnelSettings.imageColor[pixIdx] = RGBColor{
				R: Color{Value: pixels[idx]},
				G: Color{Value: pixels[idx+1]},
				B: Color{Value: pixels[idx+2]},
				A: Color{Value: pixels[idx+3]},
			}
		}
	}

	var y int = 0
	for ; y < rotoTunnelSettings.WindowHeight*2; y++ {
		var x int = 0

		for ; x < rotoTunnelSettings.WindowWidth*2; x++ {
			angle := int(math.Atan2(float64(x-rotoTunnelSettings.WindowWidth), float64(y-rotoTunnelSettings.WindowHeight)) * (float64(rotoTunnelSettings.imageWidth) / math.Pi))
			distance := int(float64(ratio) / math.Sqrt(float64(x-rotoTunnelSettings.WindowWidth)*float64(x-rotoTunnelSettings.WindowWidth)+float64(y-rotoTunnelSettings.WindowHeight)*float64(y-rotoTunnelSettings.WindowHeight)))

			rotoTunnelSettings.angleTable[x][y] = angle
			rotoTunnelSettings.distanceTable[x][y] = distance
		}
	}
}
