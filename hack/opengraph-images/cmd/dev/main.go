package main

import (
	"log"
	"net/http"

	"github.com/jmhobbs/velvetcache.org/hack/opengraph-images/internal/opengraph"
)

func main() {
	http.HandleFunc("/", opengraph.Handler)

	log.Println("Listening on :9191")
	log.Fatal(http.ListenAndServe(":9191", nil))
}
