package main

import (
	"net/http"
)

func webserver() {
	http.HandleFunc("/mandelbrot", func(w http.ResponseWriter, r *http.Request) {
		img, err := processInput()
		if err != nil {
			http.Error(w, "Error generating image", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		img.WriteTo(w)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}
