package generator

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/eringen/gowebper"
)

type rgb struct {
	r uint8
	g uint8
	b uint8
}

func renderRaster(item asset, theme theme, cells []cell, cols int, rows int) *image.NRGBA {
	layout := assetLayout(item, cells, cols, rows)
	img := image.NewNRGBA(image.Rect(0, 0, layout.Width, layout.Height))
	if !item.Transparent {
		fillRect(img, 0, 0, layout.Width, layout.Height, mustRGB(assetBackground(item, theme)), 1)
	}

	if item.Border {
		drawStrokeRect(img, theme, 0, 0, layout.Width, layout.Height, 2, 0.9)
		drawStrokeRect(img, theme, 5, 5, layout.Width-10, layout.Height-10, 1, 0.45)
	}

	for _, c := range cells {
		opacity := 1.0
		if c.Opacity != "" {
			if parsed, err := strconv.ParseFloat(c.Opacity, 64); err == nil {
				opacity = parsed
			}
		}
		fillRect(img, layout.OffsetX+c.X*cellWidth, layout.OffsetY+c.Y*cellHeight, cellWidth, cellHeight, mustRGB(c.Color), opacity)
	}

	return img
}

func writePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func writeWebP(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gowebper.Encode(file, img, &gowebper.Options{Level: gowebper.LevelDefault})
}

func fillRect(img *image.NRGBA, x int, y int, width int, height int, c rgb, opacity float64) {
	for py := y; py < y+height; py++ {
		for px := x; px < x+width; px++ {
			blendPixel(img, px, py, c, opacity)
		}
	}
}

func drawStrokeRect(img *image.NRGBA, theme theme, x int, y int, width int, height int, strokeWidth int, opacity float64) {
	for py := y; py < y+height; py++ {
		for px := x; px < x+width; px++ {
			inTop := py < y+strokeWidth
			inBottom := py >= y+height-strokeWidth
			inLeft := px < x+strokeWidth
			inRight := px >= x+width-strokeWidth
			if !inTop && !inBottom && !inLeft && !inRight {
				continue
			}
			blendPixel(img, px, py, borderColorAt(theme, px, py, img.Bounds().Dx(), img.Bounds().Dy()), opacity)
		}
	}
}

func borderColorAt(theme theme, x int, y int, width int, height int) rgb {
	start := mustRGB(theme.GradientStart)
	mid := mustRGB(theme.GradientMid)
	end := mustRGB(theme.GradientEnd)
	t := float64(x+y) / float64(width+height-2)
	if t <= 0.5 {
		return mixRGB(start, mid, t*2)
	}
	return mixRGB(mid, end, (t-0.5)*2)
}

func mixRGB(a rgb, b rgb, t float64) rgb {
	return rgb{
		r: uint8(math.Round(float64(a.r) + (float64(b.r)-float64(a.r))*t)),
		g: uint8(math.Round(float64(a.g) + (float64(b.g)-float64(a.g))*t)),
		b: uint8(math.Round(float64(a.b) + (float64(b.b)-float64(a.b))*t)),
	}
}

func blendPixel(img *image.NRGBA, x int, y int, c rgb, opacity float64) {
	if !(image.Point{X: x, Y: y}.In(img.Bounds())) {
		return
	}
	if opacity >= 1 {
		img.SetNRGBA(x, y, color.NRGBA{R: c.r, G: c.g, B: c.b, A: 255})
		return
	}

	offset := img.PixOffset(x, y)
	dstAlpha := float64(img.Pix[offset+3]) / 255
	outAlpha := opacity + dstAlpha*(1-opacity)
	if outAlpha == 0 {
		return
	}
	img.Pix[offset+0] = compositeChannel(img.Pix[offset+0], dstAlpha, c.r, opacity, outAlpha)
	img.Pix[offset+1] = compositeChannel(img.Pix[offset+1], dstAlpha, c.g, opacity, outAlpha)
	img.Pix[offset+2] = compositeChannel(img.Pix[offset+2], dstAlpha, c.b, opacity, outAlpha)
	img.Pix[offset+3] = uint8(math.Round(outAlpha * 255))
}

func compositeChannel(dst uint8, dstAlpha float64, src uint8, srcAlpha float64, outAlpha float64) uint8 {
	return uint8(math.Round((float64(src)*srcAlpha + float64(dst)*dstAlpha*(1-srcAlpha)) / outAlpha))
}

func mustRGB(value string) rgb {
	c, err := parseHexRGB(value)
	if err != nil {
		panic(err)
	}
	return c
}

func parseHexRGB(value string) (rgb, error) {
	value = strings.TrimPrefix(value, "#")
	if len(value) != 6 {
		return rgb{}, strconv.ErrSyntax
	}

	parsed, err := strconv.ParseUint(value, 16, 32)
	if err != nil {
		return rgb{}, err
	}
	return rgb{r: uint8(parsed >> 16), g: uint8(parsed >> 8), b: uint8(parsed)}, nil
}
