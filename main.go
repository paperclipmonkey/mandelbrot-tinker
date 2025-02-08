package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	port := 80
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	log.Printf("Server opening on port %d", port)

	webserver(port)
}
