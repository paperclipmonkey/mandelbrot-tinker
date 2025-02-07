package main

import (
	"io"
	"net/http"
	"net/http/httptest"
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

func TestLivezHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	handleLivez(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ok" {
		t.Errorf("expected ok got %v", string(data))
	}
}

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handleHealthz(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ok" {
		t.Errorf("expected ok got %v", string(data))
	}
}

func TestMandelbrotHandler(t *testing.T) {
	tests := []struct {
		name string
		z    string
		x    string
		y    string
	}{
		{"origin", "0", "0", "0"},
		{"zoom1_topleft", "1", "0", "0"},
		{"zoom1_bottomright", "1", "1", "1"},
		{"zoom2_center", "2", "1", "1"},
		{"zoom3_detail", "3", "3", "3"},
		{"zoom4_edge", "4", "8", "8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/mandelbrot/", nil)
			req.SetPathValue("z", tt.z)
			req.SetPathValue("x", tt.x)
			req.SetPathValue("y", tt.y)

			w := httptest.NewRecorder()
			handleMandelbrot(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK got %v", res.StatusCode)
			}

			contentType := res.Header.Get("Content-Type")
			if contentType != "image/png" {
				t.Errorf("expected Content-Type image/png got %v", contentType)
			}
		})
	}
}

func TestMandelbrotValidationHandler(t *testing.T) {
	tests := []struct {
		name string
		z    string
		x    string
		y    string
	}{
		{"text_in_int", "test", "10", "0"},
		{"negative_zoom", "-1", "0", "0"},
		{"deep_zoom", "100", "1", "1"},
		{"out_of_bounds_x", "1", "4", "1"},
		{"out_of_bounds_y", "1", "1", "4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/mandelbrot/", nil)
			req.SetPathValue("z", tt.z)
			req.SetPathValue("x", tt.x)
			req.SetPathValue("y", tt.y)

			w := httptest.NewRecorder()
			handleMandelbrot(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != 400 {
				t.Errorf("in test %v", tt.name)
				t.Errorf("expected status 400 got %v", res.StatusCode)
			}
		})
	}
}
