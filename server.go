package main

import (
	"log"
	"net/http"
	"strconv"
)

func webserver(port int) {
	http.HandleFunc("/mandelbrot/{z}/{x}/{y}", handleMandelbrot)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/livez", handleLivez)

	http.HandleFunc("/healthz", handleHealthz)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
