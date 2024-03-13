package opengraph

import (
	_ "embed"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

//go:embed SourceCodePro-Regular.ttf
var fontFile []byte
var font *truetype.Font

func init() {
	var err error

	font, err = truetype.Parse(fontFile)
	if err != nil {
		panic(err)
	}
}

func fitTypeToBox(text string, maxWidth, maxHeight, lineSpacing, baseFontSize float64) (float64, float64, float64) {
	fontSize := baseFontSize

	dc := gg.NewContext(int(maxWidth), int(maxHeight))
	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	dc.SetFontFace(face)

	wrapped := dc.WordWrap(text, maxWidth)
	textWidth, textHeight := dc.MeasureMultilineString(strings.Join(wrapped, "\n"), lineSpacing)

	if textHeight > maxHeight || textWidth > maxWidth {
		for {
			fontSize -= 1.0

			face = truetype.NewFace(font, &truetype.Options{Size: fontSize})
			dc.SetFontFace(face)

			wrapped = dc.WordWrap(text, maxWidth)
			textWidth, textHeight = dc.MeasureMultilineString(strings.Join(wrapped, "\n"), lineSpacing)
			if textHeight < maxHeight && textWidth < maxWidth {
				break
			}
		}
	} else {
		lastTextWidth := textWidth
		lastTextHeight := textHeight
		for {
			fontSize += 1.0
			if fontSize > maxFontSize {
				break
			}

			face = truetype.NewFace(font, &truetype.Options{Size: fontSize})
			dc.SetFontFace(face)

			wrapped = dc.WordWrap(text, maxWidth)
			textWidth, textHeight = dc.MeasureMultilineString(strings.Join(wrapped, "\n"), lineSpacing)
			if textHeight > maxHeight || textWidth > maxWidth {
				textWidth = lastTextWidth
				textHeight = lastTextHeight
				fontSize -= 1.0
				break
			} else {
				lastTextWidth = textWidth
				lastTextHeight = textHeight
			}
		}
	}

	return fontSize, textWidth, textHeight
}
