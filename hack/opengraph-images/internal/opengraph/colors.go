package opengraph

import "strconv"

var colors []string = []string{
	"#e5b567",
	"#b4d273",
	"#e87d3e",
	"#9e86c8",
	"#b05279",
	"#6c99bb",
}

var gradientColors []string = []string{
	"#e5b567",
	"#b4d273",
	"#e87d3e",
	"#9e86c8",
	"#b05279",
	"#6c99bb",
	"RdBu",
	"RdYlBu",
	"RdYlGn",
	"Spectral",
	"Turbo",
	"Viridis",
	"Inferno",
	"Plasma",
	"Warm",
	"Cool",
	"YlOrRd",
	"Rainbow",
	"Sinebow",
	"gold-hotpink-darkturquoise",
	"deeppink-gold-seagreen",
}

// #xxxxxx -> (r, g, b)
func hexColorToRGB(color string) (uint8, uint8, uint8, error) {
	red, err := strconv.ParseUint(color[1:3], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	green, err := strconv.ParseUint(color[3:5], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	blue, err := strconv.ParseUint(color[5:7], 16, 8)
	return uint8(red), uint8(green), uint8(blue), err
}
