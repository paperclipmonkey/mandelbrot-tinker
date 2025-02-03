package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Welcome to my website!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestStaticFiles(t *testing.T) {
	req, err := http.NewRequest("GET", "/static/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	fs := http.FileServer(http.Dir("static/"))
	handler := http.StripPrefix("/static/", fs)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHealthz(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

func TestLivez(t *testing.T) {
	req, err := http.NewRequest("GET", "/livez", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

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
