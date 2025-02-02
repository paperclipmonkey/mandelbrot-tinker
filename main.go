package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"os"
	"strconv"
)

func main() {
	port := 80
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	log.Printf("Server opened on port %d", port)

	webserver(port)
}

// Convert from slippy map z/x/y to mandelbrot coordinates in the -2 to 2 range
func slippyToMandelbrot(z, x, y int) (float64, float64, float64, float64) {
	n := 1 << uint(z)
	xmin := float64(x)/float64(n)*4 - 2
	xmax := float64(x+1)/float64(n)*4 - 2
	ymin := float64(y)/float64(n)*4 - 2
	ymax := float64(y+1)/float64(n)*4 - 2
	return xmin, ymin, xmax, ymax
}

func processInput(xmin float64, ymin float64, xmax float64, ymax float64, width int, height int) (io.WriterTo, error) {
	log.Printf("xmin: %f, ymin: %f, xmax: %f, ymax: %f, width: %d, height: %d", xmin, ymin, xmax, ymax, width, height)
	c := complexMatrix(xmin, xmax, ymin, ymax, 256)
	members := getMembers(c, 20)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for _, member := range members {
		// Map the complex number to pixel coordinates
		px := int((real(member) - xmin) / (xmax - xmin) * float64(width))
		py := int((imag(member) - ymin) / (ymax - ymin) * float64(height))

		// Set the pixel color
		if px >= 0 && px < width && py >= 0 && py < height {
			img.Set(px, py, color.RGBA{R: 255, B: 128, A: 255})
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

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

func linspace(start, end float64, num int) []float64 {
	step := (end - start) / float64(num)
	result := make([]float64, num)
	for i := 0; i < num; i++ {
		result[i] = start + float64(i)*step
	}
	return result
}

func isStable(c complex128, numIterations int) bool {
	z := complex(0, 0)
	for i := 0; i < numIterations; i++ {
		z = z*z + c
	}
	return cmplx.Abs(z) <= 2
}

func getMembers(c [][]complex128, numIterations int) []complex128 {
	var members []complex128
	for _, row := range c {
		for _, val := range row {
			if isStable(val, numIterations) {
				members = append(members, val)
			}
		}
	}
	return members
}
