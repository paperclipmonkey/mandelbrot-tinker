// Package main implements a Mandelbrot set visualization generator.
//
// The package provides functionality to:
// - Convert slippy map coordinates to Mandelbrot set coordinates
// - Generate complex number matrices
// - Process and visualize the Mandelbrot set
// - Create color-coded PNG images based on set membership
//
// The implementation uses the Go standard library's image package for image creation
// and the go-colorful package for color manipulation and gradient generation.

package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"math"
	"math/cmplx"

	"github.com/lucasb-eyer/go-colorful"
)

// Convert from slippy map z/x/y to mandelbrot coordinates in the -2 to 2 range
func slippyToMandelbrot(z, x, y int) (float64, float64, float64, float64) {
	n := 1 << uint(z) // 2^z
	xmin := float64(x)/float64(n)*4 - 2
	xmax := float64(x+1)/float64(n)*4 - 2
	ymin := float64(y)/float64(n)*4 - 2
	ymax := float64(y+1)/float64(n)*4 - 2
	return xmin, ymin, xmax, ymax
}

// Process the input parameters and generate a Mandelbrot set image
func processInput(xmin float64, ymin float64, xmax float64, ymax float64, width int, height int) (io.WriterTo, error) {
	// log.Printf("xmin: %f, ymin: %f, xmax: %f, ymax: %f, width: %d, height: %d", xmin, ymin, xmax, ymax, width, height)
	c := complexMatrix(xmin, xmax, ymin, ymax, 256)

	// Flatten the matrix into a slice
	flattened := make([]complex128, 0, len(c)*len(c[0]))
	for _, row := range c {
		flattened = append(flattened, row...)
	}
	members := flattened

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Set background to white
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color, _ := colorful.Hex("#ffffff")
			img.Set(x, y, color)
		}
	}

	for _, member := range members {
		// Map the complex number to pixel coordinates
		px := int((real(member) - xmin) / (xmax - xmin) * float64(width))
		py := int((imag(member) - ymin) / (ymax - ymin) * float64(height))

		// Set the pixel color
		if px >= 0 && px < width && py >= 0 && py < height {
			stability := isStable(member, 35)
			if cmplx.Abs(stability) > 2 || cmplx.IsNaN(stability) {
				color, _ := colorful.Hex("#ffffff")
				img.Set(px, py, color)
				continue
			}
			// Map the stability value to a color gradient
			magnitude := cmplx.Abs(stability)
			hue := 1.0 - math.Min(magnitude/2.0, 1.0) // Convert magnitude to a hue value between 0 and 1
			colord := colorful.Hsl(hue*360, 0.5, 0.5)
			img.Set(px, py, colord)
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

// Generate a 2d matrix of complex numbers
func complexMatrix(xmin, xmax, ymin, ymax float64, pixelDensity int) [][]complex128 {
	re := linspace(xmin, xmax, pixelDensity)
	im := linspace(ymin, ymax, pixelDensity)
	matrix := make([][]complex128, pixelDensity)
	for i := range matrix {
		matrix[i] = make([]complex128, pixelDensity)
		for j := range matrix[i] {
			matrix[i][j] = complex(re[j], im[i])
		}
	}
	return matrix
}

// Generate a linearly spaced slice of floats across an arbitrary range
func linspace(start, end float64, num int) []float64 {
	step := (end - start) / float64(num)
	result := make([]float64, num)
	for i := 0; i < num; i++ {
		result[i] = start + float64(i)*step
	}
	return result
}

// Check if a complex number is stable after a given number of iterations
func isStable(c complex128, numIterations int) complex128 {
	z := complex(0, 0)
	for i := 0; i < numIterations; i++ {
		z = z*z + c
	}
	return z
}
