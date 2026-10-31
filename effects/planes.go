package effects

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const planeDistance = 320.0

// Point3Df matches the original C++ effect's coordinate type.
type Point3Df struct {
	X float32
	Y float32
	Z float32
}

type PlaneSettings struct {
	Renderer      *sdl.Renderer
	WindowWidth   int32
	WindowHeight  int32
	Texture       []RGBColor
	FrameTexture  *sdl.Texture
	PixelBuffer   []byte
	TextureWidth  int
	TextureHeight int
}

var planeSettings = PlaneSettings{}
var planeFrame float64

const maxPlaneDeltaSeconds = 0.25

func advancePlaneFrame(deltaSeconds float64) float64 {
	if deltaSeconds <= 0 {
		deltaSeconds = 1.0 / 60.0
	}
	if deltaSeconds > maxPlaneDeltaSeconds {
		deltaSeconds = maxPlaneDeltaSeconds
	}
	planeFrame += deltaSeconds * 20.0
	return planeFrame
}

type planeProjection struct {
	ax float32
	ay float32
	az float32
	bx float32
	by float32
	bz float32
	cx float32
	cy float32
	cz float32
}

func planeProjectionForFrame(frame float32) planeProjection {
	bp := Point3Df{X: frame, Y: -16, Z: frame}
	up := Point3Df{X: 256, Y: 0, Z: 0}
	vp := Point3Df{X: 0, Y: 0, Z: 256}

	return planeProjection{
		cx: (up.Y*vp.Z - vp.Y*up.Z),
		cy: (vp.X*up.Z - up.X*vp.Z),
		cz: (up.X*vp.Y - vp.X*up.Y) * planeDistance,
		ax: vp.Y*bp.Z - bp.Y*vp.Z,
		ay: bp.X*vp.Z - vp.X*bp.Z,
		az: (vp.X*bp.Y - bp.X*vp.Y) * planeDistance,
		bx: bp.Y*up.Z - up.Y*bp.Z,
		by: up.X*bp.Z - bp.X*up.Z,
		bz: (bp.X*up.Y - up.X*bp.Y) * planeDistance,
	}
}

func (projection planeProjection) scanline(y, width, height int32) (u, v, du, dv int32) {
	row := float32(y - height/2)
	left := -float32(width / 2)
	a := projection.az + projection.ay*row + projection.ax*left
	b := projection.bz + projection.by*row + projection.bx*left
	c := projection.cz + projection.cy*row + projection.cx*left
	invC := float32(1) / 65536
	if c < -65536 || c > 65536 {
		invC = 1 / c
	}

	const fixedScale = float32(1 << 24)
	return int32(a * fixedScale * invC),
		int32(b * fixedScale * invC),
		int32(fixedScale * projection.ax * invC),
		int32(fixedScale * projection.bx * invC)
}

func planeUVToTextureIndex(u, v int32) int {
	return (int((v>>8)&0xFF00) + int((u>>16)&0xFF))
}

func initPlanePalette() {
	for i := 0; i < 256; i++ {
		r := uint8((math.Sin(float64(i)/24.0) + 1.0) * 127.0)
		g := uint8((math.Cos(float64(i)/18.0+0.7) + 1.0) * 127.0)
		b := uint8((math.Sin(float64(i)/30.0+1.3) + 1.0) * 127.0)
		SetColor(uint8(i), r, g, b, 255)
	}
}

func initPlaneTexture() {
	planeSettings.TextureWidth = 256
	planeSettings.TextureHeight = 256
	planeSettings.Texture = make([]RGBColor, planeSettings.TextureWidth*planeSettings.TextureHeight)
	for y := 0; y < planeSettings.TextureHeight; y++ {
		for x := 0; x < planeSettings.TextureWidth; x++ {
			r := uint8(128 + 60*math.Sin(float64(x)/14.0))
			g := uint8(128 + 60*math.Cos(float64(y)/16.0))
			b := uint8(128 + 60*math.Sin(float64(x+y)/18.0))
			planeSettings.Texture[y*planeSettings.TextureWidth+x] = RGBColor{
				R: Color{Value: r},
				G: Color{Value: g},
				B: Color{Value: b},
				A: Color{Value: 255},
			}
		}
	}
}

func loadPlaneTextureFromFile(fileName string) error {
	surface, err := img.Load(fileName)
	if err != nil {
		return err
	}
	defer surface.Free()

	width := int(surface.W)
	height := int(surface.H)
	planeSettings.TextureWidth = width
	planeSettings.TextureHeight = height
	planeSettings.Texture = make([]RGBColor, width*height)

	palette := surface.Format.Palette
	if palette == nil {
		rgbaSurface, err := surface.ConvertFormat(uint32(sdl.PIXELFORMAT_RGBA32), 0)
		if err != nil {
			return err
		}
		defer rgbaSurface.Free()
		pixels := rgbaSurface.Pixels()
		pitch := int(rgbaSurface.Pitch)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*pitch + x*4
				planeSettings.Texture[y*width+x] = RGBColor{
					R: Color{Value: pixels[idx]},
					G: Color{Value: pixels[idx+1]},
					B: Color{Value: pixels[idx+2]},
					A: Color{Value: pixels[idx+3]},
				}
			}
		}
		return nil
	}

	pixels := surface.Pixels()
	pitch := int(surface.Pitch)
	paletteColors := unsafe.Slice(palette.Colors, int(palette.Ncolors))
	for y := 0; y < height; y++ {
		rowStart := y * pitch
		for x := 0; x < width; x++ {
			idx := rowStart + x
			paletteIndex := int(pixels[idx])
			if paletteIndex < 0 || paletteIndex >= len(paletteColors) {
				paletteIndex = 0
			}

			color := paletteColors[paletteIndex]
			pixIdx := y*width + x
			planeSettings.Texture[pixIdx] = RGBColor{
				R: Color{Value: color.R},
				G: Color{Value: color.G},
				B: Color{Value: color.B},
				A: Color{Value: color.A},
			}
		}
	}

	return nil
}

func initializePlaneFrameTexture() error {
	texture, err := planeSettings.Renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		planeSettings.WindowWidth,
		planeSettings.WindowHeight,
	)
	if err != nil {
		return err
	}
	planeSettings.FrameTexture = texture
	planeSettings.PixelBuffer = make([]byte, int(planeSettings.WindowWidth*planeSettings.WindowHeight)*4)
	return nil
}

func InitializePlane(renderer *sdl.Renderer, windowWidth int32, windowHeight int32, imageFileName ...string) {
	planeSettings.Renderer = renderer
	planeSettings.WindowWidth = windowWidth
	planeSettings.WindowHeight = windowHeight
	planeSettings.Texture = nil
	planeSettings.TextureWidth = 0
	planeSettings.TextureHeight = 0
	planeSettings.PixelBuffer = nil
	if planeSettings.FrameTexture != nil {
		planeSettings.FrameTexture.Destroy()
		planeSettings.FrameTexture = nil
	}
	initPlanePalette()

	if len(imageFileName) > 0 && imageFileName[0] != "" {
		if err := loadPlaneTextureFromFile(imageFileName[0]); err != nil {
			fmt.Printf("Could not load plane image %q: %v; falling back to procedural texture\n", imageFileName[0], err)
			initPlaneTexture()
		}
	} else {
		initPlaneTexture()
	}

	if err := initializePlaneFrameTexture(); err != nil {
		fmt.Println("Could not create plane frame texture:", err)
	}
}

func renderPlaneFrame() {
	if planeSettings.FrameTexture == nil || len(planeSettings.PixelBuffer) == 0 || len(planeSettings.Texture) == 0 {
		return
	}
	width := planeSettings.WindowWidth
	height := planeSettings.WindowHeight
	if width <= 0 || height <= 0 || planeSettings.TextureWidth <= 0 || planeSettings.TextureHeight <= 0 {
		return
	}

	projection := planeProjectionForFrame(float32(planeFrame))
	textureWidth := planeSettings.TextureWidth
	textureHeight := planeSettings.TextureHeight
	for y := int32(0); y < height; y++ {
		u, v, du, dv := projection.scanline(y, width, height)
		rowOffset := int(y * width * 4)
		for x := int32(0); x < width; x++ {
			textureIndex := planeUVToTextureIndex(u, v)
			texX := (textureIndex & 0xff) % textureWidth
			texY := ((textureIndex >> 8) & 0xff) % textureHeight
			color := planeSettings.Texture[texY*textureWidth+texX]
			pixelOffset := rowOffset + int(x)*4
			planeSettings.PixelBuffer[pixelOffset] = color.R.Value
			planeSettings.PixelBuffer[pixelOffset+1] = color.G.Value
			planeSettings.PixelBuffer[pixelOffset+2] = color.B.Value
			planeSettings.PixelBuffer[pixelOffset+3] = color.A.Value
			u += du
			v += dv
		}
	}

	if err := planeSettings.FrameTexture.Update(nil, unsafe.Pointer(&planeSettings.PixelBuffer[0]), int(width)*4); err != nil {
		fmt.Println("Could not update plane frame texture:", err)
		return
	}
	if err := planeSettings.Renderer.Copy(planeSettings.FrameTexture, nil, nil); err != nil {
		fmt.Println("Could not copy plane frame texture:", err)
	}
}

func RenderPlane(deltatime float64) {
	advancePlaneFrame(deltatime)

	if planeSettings.Renderer != nil {
		planeSettings.Renderer.SetDrawColor(0, 0, 0, 255)
		planeSettings.Renderer.Clear()
	}

	renderPlaneFrame()
}
