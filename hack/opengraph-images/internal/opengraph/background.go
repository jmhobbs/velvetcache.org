package opengraph

import (
	"hash/crc64"
	"log"
	"strings"

	"github.com/fogleman/gg"
	"github.com/mazznoer/colorgrad"
	"github.com/mazznoer/csscolorparser"
	"github.com/ojrac/opensimplex-go"
)

var gradients []string = []string{
	"gradient-linear-ltr",
	"gradient-linear-rtl",
	"gradient-linear-ttb",
	"gradient-linear-btt",
	"gradient-linear-tlbr",
	"gradient-linear-trbl",
	"gradient-linear-bltr",
	"gradient-linear-brtl",
	"gradient-noise",
	"gradient-noise-sharp",
}

func drawBackground(dc *gg.Context, title string, pattern string, baseColor string) error {
	return drawBackgroundGradient(dc, pattern, baseColor, makeNoiseSeedFromString(title))
}

const GRADIENT_DENSITY uint = 100

func drawBackgroundGradient(dc *gg.Context, pattern string, baseColor string, noiseSeed int64) error {
	var fromX, fromY, toX, toY float64
	switch pattern {
	case "gradient-linear-ltr":
		fromX, fromY, toX, toY = 0, 0, outputWidth, 0
	case "gradient-linear-rtl":
		fromX, fromY, toX, toY = outputWidth, 0, 0, 0
	case "gradient-linear-ttb":
		fromX, fromY, toX, toY = 0, 0, 0, outputHeight
	case "gradient-linear-btt":
		fromX, fromY, toX, toY = 0, outputHeight, 0, 0
	case "gradient-linear-tlbr":
		fromX, fromY, toX, toY = 0, 0, outputWidth, outputHeight
	case "gradient-linear-trbl":
		fromX, fromY, toX, toY = outputWidth, 0, 0, outputHeight
	case "gradient-linear-bltr":
		fromX, fromY, toX, toY = 0, outputHeight, outputWidth, 0
	case "gradient-linear-brtl":
		fromX, fromY, toX, toY = outputWidth, outputHeight, 0, 0
	case "gradient-noise-sharp":
	case "gradient-noise":
		// noop
	default:
		log.Fatalf("unknown gradient layout %q", pattern)
	}

	var (
		colors []csscolorparser.Color
		grad   colorgrad.Gradient
		err    error
	)

	switch baseColor {
	case "RdBu":
		grad = colorgrad.RdBu()
	case "RdYlBu":
		grad = colorgrad.RdYlBu()
	case "RdYlGn":
		grad = colorgrad.RdYlGn()
	case "Spectral":
		grad = colorgrad.Spectral()
	case "Turbo":
		grad = colorgrad.Turbo()
	case "Viridis":
		grad = colorgrad.Viridis()
	case "Inferno":
		grad = colorgrad.Inferno()
	case "Plasma":
		grad = colorgrad.Plasma()
	case "Warm":
		grad = colorgrad.Warm()
	case "Cool":
		grad = colorgrad.Cool()
	case "YlOrRd":
		grad = colorgrad.YlOrRd()
	case "Rainbow":
		grad = colorgrad.Rainbow()
	case "Sinebow":
		grad = colorgrad.Sinebow()
	default:
		grad, err = colorgrad.NewGradient().
			HtmlColors(strings.Split(baseColor, "-")...).
			Build()
		if err != nil {
			log.Fatalf("could not build gradient %q: %v", baseColor, err)
		}
	}

	if pattern == "gradient-noise" || pattern == "gradient-noise-sharp" {
		if pattern == "gradient-noise-sharp" {
			grad = grad.Sharp(7, 0.2)
		}
		noise := opensimplex.NewNormalized(noiseSeed)
		for y := 0; y < int(outputHeight); y++ {
			for x := 0; x < int(outputWidth); x++ {
				t := noise.Eval2(float64(x)*0.0025, float64(y)*0.0025)
				dc.SetColor(grad.At(t))
				dc.SetPixel(x, y)
			}
		}
		return nil
	}

	colors = grad.Colors(GRADIENT_DENSITY)
	gradient := gg.NewLinearGradient(fromX, fromY, toX, toY)

	for i, c := range colors {
		gradient.AddColorStop(float64(i)/float64(len(colors)), c)
	}
	dc.SetFillStyle(gradient)
	dc.DrawRectangle(0, 0, outputWidth, outputHeight)
	dc.Fill()
	return nil
}

func makeNoiseSeedFromString(in string) int64 {
	table := crc64.MakeTable(crc64.ISO)
	return int64(crc64.Checksum([]byte(in), table))
}
