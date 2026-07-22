package effects

import (
	_ "image/jpeg"

	"github.com/veandco/go-sdl2/sdl"
)

func PutPixel(renderer *sdl.Renderer, x int32, y int32, r uint8, g uint8, b uint8, a uint8) {
	renderer.SetDrawColor(r, g, b, a)
	renderer.DrawPoint(x, y)
}

func Make2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

// func ReadImage() {

// 	if reader, err := os.Open(filepath.Join(dir_to_scan, imgFile.Name())); err == nil {
// 		defer reader.Close()
// 		im, _, err := image.DecodeConfig(reader)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
// 			continue
// 		}
// 		fmt.Printf("%s %d %d\n", imgFile.Name(), im.Width, im.Height)
// 	} else {
// 		fmt.Println("Impossible to open the file:", err)
// 	}

// }
