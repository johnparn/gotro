package effects

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	copperLineCount     = 9
	copperLineHeight    = 9
	copperRotationSpeed = 0.8
)

type CopperColor struct {
	Center sdl.Color
	Edge   sdl.Color
}

type CopperPalette [copperLineCount]CopperColor

type copperSettings struct {
	renderer  *sdl.Renderer
	width     int
	height    int
	startedAt time.Time
	colors    CopperPalette
	angle     float64
	lastFrame time.Time
	rotating  bool
	settling  bool
}

var copper = copperSettings{}

// InitializeCopper prepares the moving copper-line effect with a color pair for each line.
func InitializeCopper(renderer *sdl.Renderer, windowWidth int32, windowHeight int32, colors CopperPalette) {
	copper.renderer = renderer
	copper.width = int(windowWidth)
	copper.height = int(windowHeight)
	copper.startedAt = time.Now()
	copper.colors = colors
	copper.angle = 0
	copper.lastFrame = copper.startedAt
	copper.rotating = false
	copper.settling = false
}

// SetCopperRotation enables bar rotation or lets the current rotation finish horizontally.
func SetCopperRotation(enabled bool) {
	copper.rotating = enabled
	if enabled {
		copper.settling = false
		return
	}

	nearestHorizontal := math.Round(copper.angle / math.Pi)
	if math.Abs(copper.angle-nearestHorizontal*math.Pi) < 1e-9 {
		copper.angle = nearestHorizontal * math.Pi
		copper.settling = false
		return
	}
	copper.settling = true
}

// RunCopper draws animated copper bars and returns any SDL rendering error.
func RunCopper() error {
	if copper.renderer == nil || copper.width <= 0 || copper.height <= 0 {
		return nil
	}

	updateCopperRotation(time.Now())
	elapsed := time.Since(copper.startedAt).Seconds()
	centerY := float64(copper.height) / 2
	amplitude := float64(copper.height) * 0.25
	for line := 0; line < copperLineCount; line++ {
		phase := 2 * float64(line) / copperLineCount
		yMotion := math.Sin(elapsed*0.4+phase) * math.Sin(elapsed*1.1+phase*0.5)
		yPosition := centerY + amplitude*yMotion
		y := int32(math.Round(yPosition))

		colors := copper.colors[line]
		if err := drawCopper(y, colors.Center, colors.Edge, copperLineHeight+(line*2), copper.angle); err != nil {
			return fmt.Errorf("draw copper bar %d: %w", line, err)
		}
	}
	return nil
}

func updateCopperRotation(now time.Time) {
	delta := now.Sub(copper.lastFrame).Seconds()
	copper.lastFrame = now
	if !copper.rotating && !copper.settling {
		return
	}

	nextAngle := copper.angle + delta*copperRotationSpeed
	if copper.settling {
		targetAngle := math.Ceil(copper.angle/math.Pi) * math.Pi
		if nextAngle >= targetAngle {
			copper.angle = targetAngle
			copper.settling = false
			return
		}
	}
	copper.angle = nextAngle
}

func drawCopper(y int32, centerColor sdl.Color, edgeColor sdl.Color, height int, angle float64) error {
	if height <= 0 {
		return nil
	}

	center := float64(height-1) / 2
	centerX := float64(copper.width) / 2
	centerY := float64(y) + center
	cosAngle, sinAngle := math.Cos(angle), math.Sin(angle)
	halfWidth := float64(copper.width) / 2
	maxDistance := center
	for offset := 0; offset < height; offset++ {
		distance := 0.0
		if maxDistance > 0 {
			distance = math.Abs(float64(offset)-center) / maxDistance
		}
		color := interpolateCopperColor(centerColor, edgeColor, distance)

		if err := copper.renderer.SetDrawColor(color.R, color.G, color.B, color.A); err != nil {
			return fmt.Errorf("set copper draw color: %w", err)
		}

		localY := float64(offset) - center
		x1 := centerX - halfWidth*cosAngle - localY*sinAngle
		y1 := centerY - halfWidth*sinAngle + localY*cosAngle
		x2 := centerX + halfWidth*cosAngle - localY*sinAngle
		y2 := centerY + halfWidth*sinAngle + localY*cosAngle
		if err := copper.renderer.DrawLine(
			int32(math.Round(x1)),
			int32(math.Round(y1)),
			int32(math.Round(x2)),
			int32(math.Round(y2)),
		); err != nil {
			return fmt.Errorf("draw copper line: %w", err)
		}
	}
	return nil
}

func interpolateCopperColor(start sdl.Color, end sdl.Color, amount float64) sdl.Color {
	interpolate := func(from uint8, to uint8) uint8 {
		return uint8(math.Round(float64(from) + (float64(to)-float64(from))*amount))
	}

	return sdl.Color{
		R: interpolate(start.R, end.R),
		G: interpolate(start.G, end.G),
		B: interpolate(start.B, end.B),
		A: interpolate(start.A, end.A),
	}
}
