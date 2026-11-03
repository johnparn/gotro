package effects

import (
	"math"
	"testing"
	"time"
)

func TestCopperRotationSettlesAtNextHorizontalAngle(t *testing.T) {
	previous := copper
	defer func() {
		copper = previous
	}()

	start := time.Now()
	copper.angle = math.Pi * 0.75
	copper.lastFrame = start
	copper.settling = true

	updateCopperRotation(start.Add(time.Second * 10))
	if copper.angle != math.Pi {
		t.Fatalf("settled angle = %f, want %f", copper.angle, math.Pi)
	}
	if copper.settling {
		t.Fatal("rotation is still settling after reaching horizontal")
	}

	updateCopperRotation(start.Add(time.Second * 11))
	if copper.angle != math.Pi {
		t.Fatalf("angle after settling = %f, want %f", copper.angle, math.Pi)
	}
}
