package generator

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eringen/gowebper"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type rgb struct {
	r uint8
	g uint8
	b uint8
}

func renderRaster(item asset, theme theme, cells []cell, cols int, rows int) *image.NRGBA {
	if item.OpenGraph {
		return renderOpenGraphRaster(item, theme, cells, cols, rows)
	}

	layout := assetLayout(item, cells, cols, rows)
	img := image.NewNRGBA(image.Rect(0, 0, layout.Width, layout.Height))
	if !item.Transparent {
		fillRect(img, 0, 0, layout.Width, layout.Height, mustRGB(assetBackground(item, theme)), 1)
	}

	if item.Border {
		border := borderMetrics(item)
		drawStrokeRect(img, theme, 0, 0, layout.Width, layout.Height, border.OuterStroke, 0.9)
		drawStrokeRect(img, theme, border.InnerInset, border.InnerInset, layout.Width-border.InnerInset*2, layout.Height-border.InnerInset*2, border.InnerStroke, 0.45)
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

func renderOpenGraphRaster(item asset, theme theme, cells []cell, cols int, rows int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, item.Width, item.Height))
	fillRect(img, 0, 0, item.Width, item.Height, mustRGB(assetBackground(item, theme)), 1)
	cardLayer := image.NewNRGBA(img.Bounds())
	drawOpenGraphCodeRaster(cardLayer, theme)
	blurNRGBA(cardLayer, 6)
	compositeLayer(img, cardLayer, 1)

	contentWidth := cols * cellWidth
	contentHeight := rows * cellHeight
	scale := math.Min(0.38, float64(item.Width-220)/float64(contentWidth))
	logoWidth := float64(contentWidth) * scale
	logoHeight := float64(contentHeight) * scale
	logoX := (float64(item.Width) - logoWidth) / 2
	logoY := (float64(item.Height) - logoHeight) / 2
	panelPaddingX := 44.0
	panelPaddingY := 36.0
	panelX := int(math.Round(logoX - panelPaddingX))
	panelY := int(math.Round(logoY - panelPaddingY))
	panelWidth := int(math.Round(logoWidth + panelPaddingX*2))
	panelHeight := int(math.Round(logoHeight + panelPaddingY*2))
	fillRect(img, panelX, panelY, panelWidth, panelHeight, mustRGB(assetBackground(item, theme)), 0.82)
	drawStrokeRect(img, theme, panelX, panelY, panelWidth, panelHeight, 3, 0.8)

	logo := renderCellsRaster(cells, cols, rows)
	drawScaledNearest(img, logo, int(math.Round(logoX)), int(math.Round(logoY)), int(math.Round(logoWidth)), int(math.Round(logoHeight)))

	return img
}

func compositeLayer(dst *image.NRGBA, src *image.NRGBA, opacity float64) {
	bounds := dst.Bounds().Intersect(src.Bounds())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBAAt(x, y)
			if c.A == 0 {
				continue
			}
			blendPixel(dst, x, y, rgb{r: c.R, g: c.G, b: c.B}, opacity*float64(c.A)/255)
		}
	}
}

func renderCellsRaster(cells []cell, cols int, rows int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, cols*cellWidth, rows*cellHeight))
	for _, c := range cells {
		opacity := 1.0
		if c.Opacity != "" {
			if parsed, err := strconv.ParseFloat(c.Opacity, 64); err == nil {
				opacity = parsed
			}
		}
		fillRect(img, c.X*cellWidth, c.Y*cellHeight, cellWidth, cellHeight, mustRGB(c.Color), opacity)
	}
	return img
}

func drawScaledNearest(dst *image.NRGBA, src *image.NRGBA, x int, y int, width int, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	for dy := 0; dy < height; dy++ {
		sy := dy * srcHeight / height
		for dx := 0; dx < width; dx++ {
			sx := dx * srcWidth / width
			c := src.NRGBAAt(sx, sy)
			if c.A == 0 {
				continue
			}
			blendPixel(dst, x+dx, y+dy, rgb{r: c.R, g: c.G, b: c.B}, float64(c.A)/255)
		}
	}
}

func blurNRGBA(img *image.NRGBA, radius int) {
	if radius <= 0 {
		return
	}
	bounds := img.Bounds()
	src := image.NewNRGBA(bounds)
	copy(src.Pix, img.Pix)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a int
			var count int
			for oy := -radius; oy <= radius; oy++ {
				py := y + oy
				if py < bounds.Min.Y || py >= bounds.Max.Y {
					continue
				}
				for ox := -radius; ox <= radius; ox++ {
					px := x + ox
					if px < bounds.Min.X || px >= bounds.Max.X {
						continue
					}
					c := src.NRGBAAt(px, py)
					r += int(c.R)
					g += int(c.G)
					b += int(c.B)
					a += int(c.A)
					count++
				}
			}
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(r / count), G: uint8(g / count), B: uint8(b / count), A: uint8(a / count)})
		}
	}
}

func drawOpenGraphCodeRaster(img *image.NRGBA, theme theme) {
	lines := openGraphCodeLines()
	d := &font.Drawer{Dst: img, Face: basicfont.Face7x13}
	for cardIndex, card := range openGraphCards() {
		drawOpenGraphRasterCard(img, d, card, lines, cardIndex)
	}
}

func drawOpenGraphRasterCard(img *image.NRGBA, d *font.Drawer, card openGraphCard, lines []string, cardIndex int) {
	fillRect(img, card.X, card.Y, card.Width, card.Height, mustRGB(card.Bg), 1)
	drawSolidStrokeRect(img, card.X+1, card.Y+1, card.Width-2, card.Height-2, 2, mustRGB(card.Border), 0.75)
	drawOpenGraphString(d, card.X+18, card.Y+28, card.Border, 1, card.Name)

	for i := 0; i < 7; i++ {
		y := card.Y + 56 + i*19
		if y > card.Y+card.Height-18 {
			break
		}
		line := truncateOpenGraphCodeLine(lines[(cardIndex*3+i)%len(lines)], card.Width)
		drawOpenGraphString(d, card.X+18, y, card.Palette[i%len(card.Palette)], 1, line)
	}
}

func drawOpenGraphString(d *font.Drawer, x int, y int, value string, opacity float64, text string) {
	c := mustRGB(value)
	d.Src = image.NewUniform(color.NRGBA{R: c.r, G: c.g, B: c.b, A: uint8(math.Round(opacity * 255))})
	d.Dot = fixed.P(x, y)
	d.DrawString(text)
}

func drawSolidStrokeRect(img *image.NRGBA, x int, y int, width int, height int, strokeWidth int, c rgb, opacity float64) {
	fillRect(img, x, y, width, strokeWidth, c, opacity)
	fillRect(img, x, y+height-strokeWidth, width, strokeWidth, c, opacity)
	fillRect(img, x, y, strokeWidth, height, c, opacity)
	fillRect(img, x+width-strokeWidth, y, strokeWidth, height, c, opacity)
}

func writePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func writeWebP(path string, img image.Image, preferCompatible bool) error {
	if preferCompatible {
		if err := writeWebPWithMagick(path, img); err == nil {
			return nil
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gowebper.Encode(file, flattenOpaque(img), &gowebper.Options{Level: gowebper.LevelDefault})
}

func writeWebPWithMagick(path string, img image.Image) error {
	if _, err := exec.LookPath("magick"); err != nil {
		return err
	}

	var pngBytes bytes.Buffer
	if err := png.Encode(&pngBytes, flattenOpaque(img)); err != nil {
		return err
	}

	cmd := exec.Command("magick", "png:-", "-quality", "90", path)
	cmd.Stdin = &pngBytes
	return cmd.Run()
}

func writePNGFromSVG(path string, svg []byte) error {
	return writeImageFromSVG(path, svg)
}

func writeWebPFromSVG(path string, svg []byte) error {
	return writeImageFromSVG(path, svg, "-quality", "90")
}

func writeImageFromSVG(path string, svg []byte, args ...string) error {
	if _, err := exec.LookPath("magick"); err != nil {
		return err
	}

	cmdArgs := append([]string{"svg:-"}, args...)
	cmdArgs = append(cmdArgs, path)
	cmd := exec.Command("magick", cmdArgs...)
	cmd.Stdin = bytes.NewReader(svg)
	return cmd.Run()
}

func flattenOpaque(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	flattened := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a == 0 {
				flattened.SetRGBA(x, y, color.RGBA{A: 255})
				continue
			}
			flattened.SetRGBA(x, y, color.RGBA{
				R: uint8((r * 0xffff / a) >> 8),
				G: uint8((g * 0xffff / a) >> 8),
				B: uint8((b * 0xffff / a) >> 8),
				A: 255,
			})
		}
	}
	return flattened
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
