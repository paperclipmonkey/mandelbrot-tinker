package main

import (
	"testing"
)

func TestSlippyToMandelbrot(t *testing.T) {
	tests := []struct {
		z, x, y                int
		xmin, ymin, xmax, ymax float64
	}{
		{0, 0, 0, -2, -2, 2, 2},
		{1, 0, 0, -2, -2, 0, 0},
		{1, 1, 1, 0, 0, 2, 2},
		{2, 1, 1, -1, -1, 0, 0},
	}

	for _, tt := range tests {
		xmin, ymin, xmax, ymax := slippyToMandelbrot(tt.z, tt.x, tt.y)
		if xmin != tt.xmin || ymin != tt.ymin || xmax != tt.xmax || ymax != tt.ymax {
			t.Errorf("slippyToMandelbrot(%d, %d, %d) = %f, %f, %f, %f; want %f, %f, %f, %f",
				tt.z, tt.x, tt.y, xmin, ymin, xmax, ymax, tt.xmin, tt.ymin, tt.xmax, tt.ymax)
		}
	}
}
