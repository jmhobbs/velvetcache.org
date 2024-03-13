package opengraph

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"

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

func Generate(title, background, baseColor string) (io.Reader, error) {
	useGradient, generatorIndex, colorIndex := stringToIndexes(title)

	if background == "" {
		if useGradient {
			background = gradients[generatorIndex%len(gradients)]
		} else {
			background = patterns[generatorIndex%len(patterns)]
		}
	}

	if baseColor == "" {
		if useGradient {
			baseColor = gradientColors[colorIndex%len(gradientColors)]
		} else {
			baseColor = colors[colorIndex%len(colors)]
		}
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

	err := drawBackground(dc, title, background, baseColor)
	if err != nil {
		return nil, fmt.Errorf("unable to draw background SVG: %w", err)
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

	ob := bytes.NewBuffer(nil)
	dc.EncodePNG(ob)

	return ob, nil
}

// take a string, hash it, then use the MSB to return values for flipping switches
func stringToIndexes(s string) (bool, int, int) {
	sum := sha256.Sum256([]byte(s))
	return sum[0] > 128, int(binary.BigEndian.Uint32(sum[1:5])), int(binary.BigEndian.Uint32(sum[5:9]))
}
