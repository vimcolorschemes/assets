package generator

import (
	"image/color"
	"testing"
)

func TestRenderRasterUsesComputedDimensionsAndCells(t *testing.T) {
	img := renderRaster(asset{Name: "unit", Padding: 20, Border: true}, testTheme(), []cell{
		{X: 0, Y: 0, Color: "#010203"},
		{X: 1, Y: 0, Color: "#040506", Opacity: shadowOpacity},
	}, 2, 1)

	if got := img.Bounds(); got.Dx() != 60 || got.Dy() != 58 {
		t.Fatalf("bounds = %v, want 60x58", got)
	}

	if got := img.NRGBAAt(0, 0); got == (color.NRGBA{R: 20, G: 23, B: 25, A: 255}) {
		t.Fatal("outer border left background-colored pixels at the image edge")
	}

	if got := img.NRGBAAt(5, 5); got == (color.NRGBA{R: 20, G: 23, B: 25, A: 255}) {
		t.Fatal("inner border was not drawn at the expected inset")
	}

	if got := img.NRGBAAt(20, 20); got != (color.NRGBA{R: 1, G: 2, B: 3, A: 255}) {
		t.Fatalf("first cell = %#v, want #010203", got)
	}

	if got := img.NRGBAAt(30, 20); got != (color.NRGBA{R: 13, G: 15, B: 17, A: 255}) {
		t.Fatalf("shadow cell = %#v, want blended #0d0f11", got)
	}
}

func TestRenderRasterScalesBorderMetrics(t *testing.T) {
	img := renderRaster(asset{Name: "unit", Padding: 20, Border: true, BorderScale: 2}, testTheme(), []cell{
		{X: 0, Y: 0, Color: "#010203"},
	}, 2, 1)

	background := color.NRGBA{R: 20, G: 23, B: 25, A: 255}
	if got := img.NRGBAAt(5, 5); got != background {
		t.Fatalf("unscaled inner border pixel = %#v, want background", got)
	}
	if got := img.NRGBAAt(10, 10); got == background {
		t.Fatal("scaled inner border was not drawn at the expected inset")
	}
}

func TestRenderRasterSquareAssetUsesEqualDimensionsAndCentersContent(t *testing.T) {
	img := renderRaster(asset{Name: "unit", Padding: 20, Square: true, Border: true}, testTheme(), []cell{
		{X: 0, Y: 0, Color: "#010203"},
	}, 2, 1)

	if got := img.Bounds(); got.Dx() != 60 || got.Dy() != 60 {
		t.Fatalf("bounds = %v, want 60x60", got)
	}

	if got := img.NRGBAAt(25, 21); got != (color.NRGBA{R: 1, G: 2, B: 3, A: 255}) {
		t.Fatalf("centered cell = %#v, want #010203", got)
	}
}

func TestRenderRasterBorderlessKeepsDimensionsAndOmitsBorder(t *testing.T) {
	img := renderRaster(asset{Name: "unit", Padding: 20}, testTheme(), []cell{
		{X: 0, Y: 0, Color: "#010203"},
	}, 2, 1)

	background := color.NRGBA{R: 20, G: 23, B: 25, A: 255}
	if got := img.Bounds(); got.Dx() != 60 || got.Dy() != 58 {
		t.Fatalf("bounds = %v, want 60x58", got)
	}
	if got := img.NRGBAAt(0, 0); got != background {
		t.Fatalf("outer border pixel = %#v, want background", got)
	}
	if got := img.NRGBAAt(5, 5); got != background {
		t.Fatalf("inner border pixel = %#v, want background", got)
	}
	if got := img.NRGBAAt(20, 20); got != (color.NRGBA{R: 1, G: 2, B: 3, A: 255}) {
		t.Fatalf("first cell = %#v, want #010203", got)
	}
}

func TestParseHexRGBRejectsInvalidColor(t *testing.T) {
	if _, err := parseHexRGB("nope"); err == nil {
		t.Fatal("parseHexRGB returned nil error for invalid color")
	}
}
