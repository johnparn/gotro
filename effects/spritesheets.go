package effects

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const SPRITE_HEIGHT int = 32
const SPRITE_WIDTH int = 32

func Spritesheet(renderer *sdl.Renderer, windowWidth int, windowHeight int) {
	// rand.Seed(time.Now().UnixNano())

	// Load the spritesheet
	spriteSheet, errSpriteSheet := img.Load("32x32-fm.png")

	// SDL_Surface *spritesheet = IMG_Load("BrogueFont5.png");
	if errSpriteSheet != nil {
		sdl.LogCritical(1, "Unable to load bitmap font: %s", sdl.GetError().Error())
		return
	} else {
		fmt.Println("Loaded font")
	}

	// Set the background color (black) to transparent
	format := spriteSheet.Format
	errSetColorKey := spriteSheet.SetColorKey(true, sdl.MapRGBA(format, 0, 0, 0, 255))
	if errSetColorKey != nil {
		sdl.LogCritical(1, "Unable to set color key on the spritesheet: %s",
			sdl.GetError().Error())
	} else {
		fmt.Println("Set color key")
	}

	// Convert the spritesheet to a texture
	tex_sheet, errTexSheet := renderer.CreateTextureFromSurface(spriteSheet)
	if errTexSheet != nil {
		sdl.LogCritical(1, "Unable to convert spritesheet to a texture: %s", sdl.GetError().Error())
		return
	} else {
		fmt.Println("Convert surface to texture")
	}

	tex_sheet.SetBlendMode(sdl.BLENDMODE_ADD)

	buffer, errBuffer := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, int32(windowWidth), int32(windowHeight))
	if errBuffer != nil {
		sdl.LogCritical(1, "Unable to create texture %s", sdl.GetError().Error())
	} else {
		fmt.Println("Create texture")
	}

	var ticks_elapsed uint64

	old_ticks := sdl.GetTicks64()
	renderer.SetRenderTarget(buffer)

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	if ticks_elapsed > 1000 {
		//     ticks_elapsed = 0;
		for i := 0; i < windowWidth/SPRITE_WIDTH+1; i++ {
			for j := 0; j < windowWidth/SPRITE_HEIGHT+1; j++ {
				glyph_x := rand.Intn(60) % 10
				glyph_y := 3 + rand.Intn(6)%6

				// Clipping rects for source (cut from spritesheet)
				// and dest (buffer location)
				dest_rect := sdl.Rect{X: int32(i * 18), Y: int32(j * 28), W: int32(18), H: int32(28)}
				src_rect := sdl.Rect{X: int32(glyph_x * 18), Y: int32(glyph_y * 28), W: int32(18), H: int32(28)}

				// Random color for background
				// SDL_SetRenderDrawColor(renderer, bg.r, bg.g, bg.b, 255);
				// SDL_RenderFillRect(renderer, &dest_rect);

				// Tint the sprite to be blitted
				tex_sheet.SetColorMod(0, 0, 0)
				renderer.Copy(tex_sheet, &src_rect, &dest_rect)
			}
		}
		// Set the renderer back to the screen, clear it, and flip the buffer
		renderer.SetRenderTarget(nil)
		renderer.Copy(buffer, nil, nil)
	}

	ticks_elapsed += sdl.GetTicks64() - old_ticks
}
