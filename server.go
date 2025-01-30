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

func webserver() {
	http.HandleFunc("/mandelbrot", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		// Parse query parameters
		params := r.URL.Query()
		// xmin := params.Get("xmin")

		xminF, err := strconv.ParseFloat(params.Get("xmin"), 64); if err != nil { xminF = -2.0 }
		if xminF == 0 {
			xminF = -2.0
		}

		yminF, _ := strconv.ParseFloat(params.Get("ymin"), 64)
		if yminF == 0 {
			yminF = -2.0
		}

		xmaxF, _ := strconv.ParseFloat(params.Get("xmax"), 64)
		if xmaxF == 0 {
			xmaxF = 2.0
		}

		ymaxF, _ := strconv.ParseFloat(params.Get("ymax"), 64)
		if ymaxF == 0 {
			ymaxF = 2.0
		}

		widthI, err := strconv.Atoi(params.Get("width")); if err != nil { widthI = 800 }
		heightI, _ := strconv.Atoi(params.Get("height"))

		img, err := processInput(xminF, yminF, xmaxF, ymaxF, widthI, heightI)
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

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}
