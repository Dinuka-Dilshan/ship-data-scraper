package main

import (
	"log"
	"net/http"

	scraper "github.com/Dinuka-Dilshan/ship-data-scraper"
)

func main() {
	http.HandleFunc("/", scraper.Handler)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
