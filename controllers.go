package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func logRequest(r *http.Request) {
	log.Printf("%s - - [%s] \"%s %s %s\" \"%s\" \"%s\"\n",
		r.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.RequestURI,
		r.Proto,
		r.Referer(),
		r.UserAgent(),
	)
}

func handleMandelbrot(w http.ResponseWriter, r *http.Request) {
	// Extract z, x, y from URL path
	logRequest(r)
	z, err := strconv.Atoi(r.PathValue("z"))
	if err != nil {
		http.Error(w, "Invalid z parameter", http.StatusBadRequest)
		return
	}
	if z < 0 || z > 24 {
		http.Error(w, "z parameter must be between 0 and 24", http.StatusBadRequest)
		return
	}
	x, err := strconv.Atoi(r.PathValue("x"))
	if err != nil {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}
	maxCoord := 1 << z
	if x < 0 || x >= maxCoord {
		http.Error(w, "x parameter must be within range for zoom level", http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(r.PathValue("y"))
	if err != nil {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}
	if y < 0 || y >= maxCoord {
		http.Error(w, "y parameter must be within range for zoom level", http.StatusBadRequest)
		return
	}
	xmin, ymin, xmax, ymax := slippyToMandelbrot(z, x, y)

	img, err := processInput(xmin, ymin, xmax, ymax, 256, 256)
	if err != nil {
		http.Error(w, "Error generating image", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	img.WriteTo(w)
}

func handleLivez(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok")) // Stub health check. If we have dependencies, check them here
}
