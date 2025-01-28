package main

import (
	"flag"
	"fmt"
	"image/color"
	"math/cmplx"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	// Define a string flag with a default value and a short description.
	input := flag.String("input", "", "Input value to be processed")

	// Parse the command-line flags.
	flag.Parse()

	// Check if the input flag was provided.
	if *input == "" {
		fmt.Println("Please provide an input value using the -input flag.")
		return
	}

	processInput()
}

func processInput() (int64, error) {
	c := complexMatrix(-2, 0.5, -1.5, 1.5, 1024)
	members := getMembers(c, 20)

	scatterData := generatePoints(members)

	p := plot.New()

	p.Title.Text = "Mandelbrot Set"

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	s.GlyphStyle.Radius = vg.Points(0.1)

	p.Add(s)

	// Save the plot to a PNG file.
	if err := p.Save(1*1000, 1*1000, "points.png"); err != nil {
		panic(err)
	}
	return 0, nil
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
