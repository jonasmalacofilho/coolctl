package driver

import (
	"encoding/hex"
	"image/color"
	"log"
)

func colorFromHexString(c string) color.RGBA {
	b, err := hex.DecodeString(c)
	if err != nil {
		log.Fatal(err)
	}

	return color.RGBA{b[0], b[1], b[2], 1}
}

func paletteFromColors(colors []string) color.Palette {
	var palette color.Palette
	if colors != nil {
		for _, c := range colors {
			palette = append(palette, colorFromHexString(c))
		}
	}

	return palette
}
