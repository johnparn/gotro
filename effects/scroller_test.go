package effects

import (
	"math"
	"testing"
)

func TestCalculateScrollerWaveYUsesSine(t *testing.T) {
	settings := scrollerSettingsStruct{
		WaveEnabled:    true,
		WaveMode:       "sine",
		WaveAmplitude:  10,
		WaveFrequency:  0.1,
		WavePhase:      0.5,
		WaveBaseOffset: 2,
	}

	got := calculateScrollerWaveY(20, 100, 0, settings)
	want := int(math.Round(float64(100+2) + float64(settings.WaveAmplitude)*math.Sin(float64(20)*settings.WaveFrequency+settings.WavePhase)))

	if got != want {
		t.Fatalf("calculateScrollerWaveY() = %d, want %d", got, want)
	}
}

func TestCalculateScrollerWaveYUsesCosine(t *testing.T) {
	settings := scrollerSettingsStruct{
		WaveEnabled:    true,
		WaveMode:       "cosine",
		WaveAmplitude:  12,
		WaveFrequency:  0.05,
		WavePhase:      0.25,
		WaveBaseOffset: -3,
	}

	got := calculateScrollerWaveY(40, 90, 0, settings)
	want := int(math.Round(float64(90-3) + float64(settings.WaveAmplitude)*math.Cos(float64(40)*settings.WaveFrequency+settings.WavePhase)))

	if got != want {
		t.Fatalf("calculateScrollerWaveY() = %d, want %d", got, want)
	}
}

func TestCalculateScrollerWaveYContinuesWhenPaused(t *testing.T) {
	settings := scrollerSettingsStruct{
		WaveEnabled:    true,
		WaveMode:       "sine",
		WaveAmplitude:  8,
		WaveFrequency:  0.1,
		WavePhase:      0.0,
		WaveBaseOffset: 0,
	}

	settings.WaveStartTime = 0
	settings.WaveSpeed = 0.01

	first := calculateScrollerWaveY(10, 100, 100, settings)
	second := calculateScrollerWaveY(10, 100, 200, settings)

	if first == second {
		t.Fatalf("expected paused wave to continue changing over time, got %d and %d", first, second)
	}
}

func TestNextScrollerOffsetUsesTimeDelta(t *testing.T) {
	const nominalFrameMs = 16.6666666667
	got := nextScrollerOffset(0, nominalFrameMs, false)
	want := 4.0

	if math.Abs(got-want) > 0.1 {
		t.Fatalf("nextScrollerOffset() = %.3f, want about %.3f", got, want)
	}
}

func TestUpdateScrollerRotationFadeOutResetsOrientation(t *testing.T) {
	settings := scrollerSettingsStruct{
		CharHeight:        3,
		InitialCharHeight: 7,
		Flipped:           true,
		RotateChars:       true,
		RotationFadingOut: true,
	}

	updateScrollerRotation(&settings)
	if settings.CharHeight != 2 || !settings.Flipped || !settings.RotationFadingOut {
		t.Fatalf("rotation fade step = height %d, flipped %t, fading %t; want 2, true, true", settings.CharHeight, settings.Flipped, settings.RotationFadingOut)
	}

	updateScrollerRotation(&settings)
	if settings.CharHeight != 1 || !settings.Flipped || !settings.RotationFadingOut {
		t.Fatalf("rotation fade penultimate step = height %d, flipped %t, fading %t; want 1, true, true", settings.CharHeight, settings.Flipped, settings.RotationFadingOut)
	}

	updateScrollerRotation(&settings)
	if settings.CharHeight != 2 || settings.Flipped || !settings.RotationFadingOut {
		t.Fatalf("rotation fade recovery = height %d, flipped %t, fading %t; want 2, false, true", settings.CharHeight, settings.Flipped, settings.RotationFadingOut)
	}

	for settings.CharHeight < settings.InitialCharHeight {
		updateScrollerRotation(&settings)
	}
	if settings.CharHeight != settings.InitialCharHeight || settings.Flipped || !settings.RotationFadingOut {
		t.Fatalf("rotation fade final frame = height %d, flipped %t, fading %t; want %d, false, true", settings.CharHeight, settings.Flipped, settings.RotationFadingOut, settings.InitialCharHeight)
	}

	updateScrollerRotation(&settings)
	if settings.CharHeight != settings.InitialCharHeight || settings.Flipped || settings.RotationFadingOut || settings.RotateChars {
		t.Fatalf("rotation fade completion = height %d, flipped %t, fading %t, rotating %t; want %d, false, false, false", settings.CharHeight, settings.Flipped, settings.RotationFadingOut, settings.RotateChars, settings.InitialCharHeight)
	}
}
