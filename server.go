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

func webserver(port int) {
	http.HandleFunc("/mandelbrot/{z}/{x}/{y}", func(w http.ResponseWriter, r *http.Request) {
		// Extract z, x, y from URL path
		logRequest(r)
		z, err := strconv.Atoi(r.PathValue("z"))
		if err != nil {
			http.Error(w, "Invalid z parameter", http.StatusBadRequest)
			return
		}
		x, err := strconv.Atoi(r.PathValue("x"))
		if err != nil {
			http.Error(w, "Invalid x parameter", http.StatusBadRequest)
			return
		}
		y, err := strconv.Atoi(r.PathValue("y"))
		if err != nil {
			http.Error(w, "Invalid y parameter", http.StatusBadRequest)
			return
		}
		log.Printf("%d %d %d", z, x, y)
		xmin, ymin, xmax, ymax := slippyToMandelbrot(z, x, y)

		log.Printf("%f %f %f %f", xmin, ymin, xmax, ymax)

		img, err := processInput(xmin, ymin, xmax, ymax, 256, 256)
		if err != nil {
			http.Error(w, "Error generating image", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		img.WriteTo(w)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok")) // Stub health check. If we have dependencies, check them here
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
