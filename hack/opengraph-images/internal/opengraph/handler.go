package opengraph

import (
	_ "embed"
	"io"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	pattern := r.URL.Query().Get("pattern")
	baseColor := r.URL.Query().Get("color")

	png, err := Generate(title, pattern, baseColor)
	if err != nil {
		log.Printf("unable to generate opengraph image: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "image/png")
	io.Copy(w, png)
}
