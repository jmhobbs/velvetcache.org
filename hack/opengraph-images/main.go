package main

import (
	"bytes"
	"crypto/sha1"
	_ "embed"
	"encoding/binary"
	"log"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

const (
	outputWidth  float64 = 2000
	outputHeight float64 = 1047

	padding float64 = 40
	margin  float64 = 60

	maxFontSize float64 = 400
)

func main() {
	http.HandleFunc("/bg.png", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")

		generatorIndex, colorIndex := stringToIndexes(title)
		pattern := backgrounds[generatorIndex%len(backgrounds)]
		baseColor := colors[colorIndex%len(colors)]

		if r.URL.Query().Get("pattern") != "" {
			pattern = r.URL.Query().Get("pattern")
		}
		if r.URL.Query().Get("color") != "" {
			baseColor = "#" + r.URL.Query().Get("color")
		}

		/*
			┌───────────────────────────────────────┐
			│                                       │
			│  ┌───────────────────────────────┐    │
			│  │ ┌───────────────────────────┐ ├─┐  │
			│  │ │                           │ │x│  │
			│  │ │                           │ │x│  │
			│  │ │         Text Here         │ │x│  │
			│  │ │                           │ │x│  │
			│  │ │                           │ │x│  │
			│  │ └───────────────────────────┘ │x│  │
			│  └──┬────────────────────────────┘x│  │
			│     │xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx│  │
			│     └──────────────────────────────┘  │
			│                                       │
			└───────────────────────────────────────┘

			Built of four boxes:
			1. The pattern background, (0,0) to (outputWidth, outputHeight)

			2. The white box, width of outputWidth - margins, vertically
			   fit to text but no taller than outputHeight - margins. Centered
			   vertically based on final text size.

			3. The shadow box, same dimensions as white box but offset by
			   shadowOffset on both X and Y

			4. The text box, fit to text but with maximum dimensions of
			   the box width and height, less padding
		*/

		// Set the baseline coordinates and dimensions
		// heights and Y values will be adjusted when text height
		// is calculated
		var (
			boxX float64 = margin
			boxY float64 = margin

			textX float64 = margin + padding
			textY float64 = margin + padding

			boxWidth float64 = outputWidth - 2*margin

			textMaxWidth  float64 = boxWidth - 2*padding
			textMaxHeight float64 = outputHeight - 2*margin - 2*padding

			shadowOffset float64 = 20
			shadowX      float64 = margin + shadowOffset
			shadowY      float64 = margin + shadowOffset

			lineSpacing float64 = 1.0
			fontSize    float64 = 120
		)

		dc := gg.NewContext(int(outputWidth), int(outputHeight))

		err := drawBackground(dc, title, pattern, baseColor)
		if err != nil {
			log.Printf("unable to draw background SVG: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// find best fit text dimensions
		fontSize, _, textHeight := fitTypeToBox(title, textMaxWidth, textMaxHeight, lineSpacing, fontSize)
		dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: fontSize}))

		// calculate dependent dimensions and vertical align
		boxHeight := textHeight + padding*2
		boxY = (outputHeight - boxHeight) / 2
		textY = (outputHeight - textHeight) / 2
		shadowY = boxY + shadowOffset

		// our font has overhead space at _about_ this ratio
		// shifting the text up pulls the descenders into the
		// calculated text box and removes the extra "padding"
		var textShiftY float64 = -1 * (22.5 / 120.0) * fontSize

		// draw shadow box
		dc.SetRGB(0.15, 0.15, 0.15)
		dc.DrawRectangle(shadowX, shadowY, boxWidth, boxHeight)
		dc.Fill()

		// draw white box
		dc.SetRGB(1, 1, 1)
		dc.DrawRectangle(boxX, boxY, boxWidth, boxHeight)
		dc.Fill()

		// draw text
		dc.SetRGB(0, 0, 0)
		dc.DrawStringWrapped(title, textX, textY+textShiftY, 0, 0, textMaxWidth, lineSpacing, gg.AlignLeft)

		w.Header().Add("Content-Type", "image/png")
		dc.EncodePNG(w)
	})

	log.Println("Listening on :9191")
	log.Fatal(http.ListenAndServe(":9191", nil))
}

// TODO: Is there a better option for this? Does SHA-1 have MSB at the end?
func stringToIndexes(s string) (int, int) {
	sha1 := sha1.New()
	sha1.Write([]byte(s))
	sum := sha1.Sum(nil)
	i, err := binary.ReadVarint(bytes.NewReader(sum))
	if err != nil {
		panic(err)
	}
	if i < 0 {
		i = -i
	}
	j, err := binary.ReadVarint(bytes.NewReader(sum))
	if err != nil {
		panic(err)
	}
	if j < 0 {
		j = -j
	}
	return int(i), int(j)
}
