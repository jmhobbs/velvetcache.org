package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/pravj/geopattern"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

var backgrounds []string = []string{
	"chevrons",
	"diamonds",
	"hexagons",
	"mosaic-squares",
	"octagons",
	"overlapping-circles",
	// oksvg can not render this one
	//"plaid",
	"sine-waves",
	"tessellation",
	"triangles",
	"xes",
	// these aren't geopattern patterns, we implement them
	"gradient-linear-ltr",
	"gradient-linear-rtl",
	"gradient-linear-ttb",
	"gradient-linear-btt",
	"gradient-linear-tlbr",
	"gradient-linear-trbl",
	"gradient-linear-bltr",
	"gradient-linear-rltb",
}

func rasterizeSVG(svg string, baseColor string) (*image.RGBA, error) {
	red, green, blue, err := hexColorToRGB(baseColor)
	if err != nil {
		return nil, fmt.Errorf("could not convert color %q: %w", baseColor, err)
	}

	matches := regexp.MustCompile("<svg xmlns=.*? width='(\\d+)' height='(\\d+)'>").FindStringSubmatch(svg)
	if len(matches) != 3 {
		return nil, fmt.Errorf("could not extract dimensions from svg")
	}

	width, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to extract pattern width: %w", err)
	}
	height, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to extract pattern height: %w", err)
	}

	icon, err := oksvg.ReadIconStream(bytes.NewReader([]byte(svg)), oksvg.WarnErrorMode)
	if err != nil {
		return nil, fmt.Errorf("unable to read icon stream: %w", err)
	}

	icon.SetTarget(0, 0, float64(width), float64(height))
	rgba := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.RGBA{red, green, blue, 255}}, image.Point{}, draw.Src)

	icon.Draw(rasterx.NewDasher(int(width), int(height), rasterx.NewScannerGV(int(width), int(height), rgba, rgba.Bounds())), 1)

	return rgba, nil
}

func drawBackground(dc *gg.Context, title string, pattern string, baseColor string) error {
	if strings.HasPrefix(pattern, "gradient-") {
		return drawBackgroundGradient(dc, pattern, baseColor)
	}
	return drawBackgroundPattern(dc, title, pattern, baseColor)

}

func drawBackgroundGradient(dc *gg.Context, pattern string, baseColor string) error {
	red, green, blue, err := hexColorToRGB(baseColor)
	if err != nil {
		return fmt.Errorf("could not convert color %q: %w", baseColor, err)
	}

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
	}

	gradient := gg.NewLinearGradient(fromX, fromY, toX, toY)
	gradient.AddColorStop(0, color.RGBA{red, green, blue, 255})
	gradient.AddColorStop(1, color.RGBA{255, 255, 255, 255})
	dc.SetFillStyle(gradient)
	dc.DrawRectangle(0, 0, outputWidth, outputHeight)
	dc.Fill()
	return nil
}

func drawBackgroundPattern(dc *gg.Context, title string, pattern string, baseColor string) error {
	args := map[string]string{
		"phrase":    title,
		"generator": pattern,
		"baseColor": baseColor,
	}

	rgba, err := rasterizeSVG(geopattern.Generate(args), baseColor)
	if err != nil {
		return fmt.Errorf("unable to rasterize SVG: %w", err)
	}

	// draw pattern over whole image
	dc.SetFillStyle(gg.NewSurfacePattern(rgba, gg.RepeatBoth))
	dc.DrawRectangle(0, 0, outputWidth, outputHeight)
	dc.Fill()

	return nil
}
