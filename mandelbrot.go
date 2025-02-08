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
	"sync"

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

// New processInput: computes each pixel directly in parallel.
func processInput(xmin float64, ymin float64, xmax float64, ymax float64, width int, height int) (io.WriterTo, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	white, _ := colorful.Hex("#ffffff")
	// Precompute step sizes.
	dx := (xmax - xmin) / float64(width)
	dy := (ymax - ymin) / float64(height)

	// Set background to white.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}

	var wg sync.WaitGroup
	// Process each row concurrently.
	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				c := complex(xmin+float64(x)*dx, ymin+float64(y)*dy)
				stability := isStable(c, 35)
				if cmplx.Abs(stability) > 2 || cmplx.IsNaN(stability) {
					continue // remains white from background.
				}
				magnitude := cmplx.Abs(stability)
				hue := 1.0 - math.Min(magnitude/2.0, 1.0)
				colord := colorful.Hsl(hue*360, 0.5, 0.5)
				img.Set(x, y, colord)
			}
		}(y)
	}
	wg.Wait()

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

// Check if a complex number is stable after a given number of iterations
func isStable(c complex128, numIterations int) complex128 {
	z := complex(0, 0)
	for i := 0; i < numIterations; i++ {
		z = z*z + c
	}
	return z
}
