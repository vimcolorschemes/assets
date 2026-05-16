package generator

import (
	"image/color"
	"testing"
)

func TestRenderRasterUsesComputedDimensionsAndCells(t *testing.T) {
	img := renderRaster(asset{Name: "unit", Padding: 20}, []cell{
		{X: 0, Y: 0, Color: "#010203"},
		{X: 1, Y: 0, Color: "#040506", Opacity: shadowOpacity},
	}, 2, 1)

	if got := img.Bounds(); got.Dx() != 60 || got.Dy() != 58 {
		t.Fatalf("bounds = %v, want 60x58", got)
	}

	if got := img.RGBAAt(0, 0); got == (color.RGBA{R: 9, G: 14, B: 19, A: 255}) {
		t.Fatal("outer border left background-colored pixels at the image edge")
	}

	if got := img.RGBAAt(5, 5); got == (color.RGBA{R: 9, G: 14, B: 19, A: 255}) {
		t.Fatal("inner border was not drawn at the expected inset")
	}

	if got := img.RGBAAt(20, 20); got != (color.RGBA{R: 1, G: 2, B: 3, A: 255}) {
		t.Fatalf("first cell = %#v, want #010203", got)
	}

	if got := img.RGBAAt(30, 20); got != (color.RGBA{R: 7, G: 10, B: 14, A: 255}) {
		t.Fatalf("shadow cell = %#v, want blended #070a0e", got)
	}
}

func TestParseHexRGBRejectsInvalidColor(t *testing.T) {
	if _, err := parseHexRGB("nope"); err == nil {
		t.Fatal("parseHexRGB returned nil error for invalid color")
	}
}
