package effects

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const concentricRingLifetime = 3 * time.Second

// ConcentricInterval controls how often a new ring is spawned.
var ConcentricInterval = 500 * time.Millisecond

// ConcentricHeightScale controls ring height: 1 is circular, smaller values flatten the rings.
var ConcentricHeightScale = 1.0

// AdjustConcentricHeightScale changes ring height while keeping it between 0 and 1.
func AdjustConcentricHeightScale(delta float64) {
	scale := math.Max(0, math.Min(1, ConcentricHeightScale))
	ConcentricHeightScale = math.Max(0, math.Min(1, scale+delta))
}

type concentricRing struct {
	startedAt time.Time
}

type concentricSettings struct {
	renderer  *sdl.Renderer
	width     int
	height    int
	lastSpawn time.Time
	rings     []concentricRing
	points    []sdl.Point
	scalingStep float64
	scalingDirection int
}

var concentric = concentricSettings{}

// InitializeConcentric prepares the expanding-ring effect for the given renderer size.
func InitializeConcentric(renderer *sdl.Renderer, windowWidth int, windowHeight int) {
	concentric.renderer = renderer
	concentric.width = windowWidth
	concentric.height = windowHeight
	concentric.lastSpawn = time.Time{}
	concentric.rings = concentric.rings[:0]
	concentric.points = concentric.points[:0]
	concentric.scalingStep = 0.01
	concentric.scalingDirection = 0
}

// RunConcentric draws expanding rings that cycle through colors and fade out.
func RunConcentric() {
	if concentric.renderer == nil || concentric.width <= 0 || concentric.height <= 0 {
		return
	}

	now := time.Now()
	heightScale := math.Max(0, math.Min(1, ConcentricHeightScale))
	interval := ConcentricInterval
	if interval <= 0 {
		interval = time.Millisecond
	}
	if len(concentric.rings) == 0 || now.Sub(concentric.lastSpawn) >= interval {
		concentric.rings = append(concentric.rings, concentricRing{startedAt: now})
		concentric.lastSpawn = now
	}

	previousBlendMode := sdl.BlendMode(sdl.BLENDMODE_NONE)
	if err := concentric.renderer.GetDrawBlendMode(&previousBlendMode); err != nil {
		return
	}
	if err := concentric.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return
	}
	defer concentric.renderer.SetDrawBlendMode(previousBlendMode)

	activeRings := concentric.rings[:0]
	centerX := float64(concentric.width) / 2
	centerY := float64(concentric.height) / 2
	maxRadius := centerX
	
	// if heightScale > 0 {
	// 	maxRadius = math.Hypot(centerX, centerY/heightScale)
	// }


	if (concentric.scalingStep > 0) {	
		if (concentric.scalingDirection == 0) {
			AdjustConcentricHeightScale(-concentric.scalingStep)
		} else {
			AdjustConcentricHeightScale(concentric.scalingStep)
		}
		if (heightScale <= 0) {
			concentric.scalingDirection = 1
		} else if (heightScale >= 1) {
			concentric.scalingDirection = 0
		}
	}	

	for _, ring := range concentric.rings {
		progress := float64(now.Sub(ring.startedAt)) / float64(concentricRingLifetime) / 3
		if progress >= 1 {
			continue
		}
		activeRings = append(activeRings, ring)

		radius := maxRadius * progress
		alpha := uint8(255 * (1 - progress))
		phase := progress * 2 * math.Pi
		red := uint8(127.5 + 127.5*math.Sin(phase))
		green := uint8(127.5 + 127.5*math.Sin(phase+2*math.Pi/3))
		blue := uint8(127.5 + 127.5*math.Sin(phase+4*math.Pi/3))
		concentric.renderer.SetDrawColor(red, green, blue, alpha)

		segmentCount := int(math.Ceil(2 * math.Pi * radius))
		if segmentCount < 24 {
			segmentCount = 24
		}
		concentric.points = concentric.points[:0]
		for segment := 0; segment < segmentCount; segment++ {
			angle := 2 * math.Pi * float64(segment) / float64(segmentCount)
			x := centerX + radius*math.Cos(angle)
			y := centerY + radius*heightScale*math.Sin(angle)
			concentric.points = append(concentric.points, sdl.Point{
				X: int32(math.Round(x)),
				Y: int32(math.Round(y)),
			})
		}
		concentric.renderer.DrawPoints(concentric.points)
	}
	concentric.rings = activeRings
}
