package effects

import (
	"strconv"
)

type Color struct {
	Value     uint8
	Increment int
}

type RGBColor struct {
	R Color
	G Color
	B Color
	A Color
}

func SetColor(colorIndex uint8, r uint8, g uint8, b uint8, a uint8) {
	rgba := RGBColor{
		R: Color{Value: r},
		G: Color{Value: g},
		B: Color{Value: b},
		A: Color{Value: a},
	}
	palette[colorIndex] = rgba
}

func GetColor(color Color) Color {

	// Set increment if unset
	if color.Increment == 0 {
		if color.Value == 255 {
			color.Increment = -1

		} else {
			color.Increment = 1
		}
	}

	// Change color direction, up or down
	switch color.Value {
	case 255:
		color.Increment = -1

	case 0:
		color.Increment = 1
	}
	color.Value += uint8(color.Increment)

	return color
}

type Hex string

func (h Hex) toRGB() (RGBColor, error) {
	return Hex2RGBA(h)
}

func Hex2RGBA(hex Hex) (RGBColor, error) {
	var rgba RGBColor
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGBColor{}, err
	}

	rgba = RGBColor{
		R: Color{Value: uint8((values >> 24) & 0xFF)},
		G: Color{Value: uint8((values >> 16) & 0xFF)},
		B: Color{Value: uint8((values >> 8) & 0xFF)},
		A: Color{Value: uint8((values >> 0) & 0xFF)},
	}

	return rgba, nil
}
