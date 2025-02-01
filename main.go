package main

import (
	"image/color"
	"io"
	"log"
	"math/cmplx"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

func processInput(xmin float64, ymin float64, xmax float64, ymax float64, width int, height int) (io.WriterTo, error) {
	log.Printf("xmin: %f, ymin: %f, xmax: %f, ymax: %f, width: %d, height: %d", xmin, ymin, xmax, ymax, width, height)
	c := complexMatrix(xmin, xmax, ymin, ymax, max(width, height, 100))
	members := getMembers(c, 20)

	scatterData := generatePoints(members)

	p := plot.New()

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	s.GlyphStyle.Radius = vg.Points(0.1)

	p.Add(s)

	return p.WriterTo(vg.Length(width), vg.Length(height), "png")
}

func complexMatrix(xmin, xmax, ymin, ymax float64, pixelDensity int) [][]complex128 {
	re := linspace(xmin, xmax, int((xmax-xmin)*float64(pixelDensity)))
	im := linspace(ymin, ymax, int((ymax-ymin)*float64(pixelDensity)))
	matrix := make([][]complex128, len(im))
	for i := range matrix {
		matrix[i] = make([]complex128, len(re))
		for j := range matrix[i] {
			matrix[i][j] = complex(re[j], im[i])
		}
	}
	return matrix
}

func linspace(start, end float64, num int) []float64 {
	step := (end - start) / float64(num-1)
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

func generatePoints(ns []complex128) plotter.XYs {
	pts := make(plotter.XYs, len(ns))
	for i := range pts {
		pts[i].X = real(ns[i])
		pts[i].Y = imag(ns[i])
	}
	return pts
}
