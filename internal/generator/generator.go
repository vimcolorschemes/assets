package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/superstarryeyes/bit/ansifonts"
)

const (
	startColor = "#938aa9"
	endColor   = "#b6927b"
)

type asset struct {
	Name    string
	Text    string
	OutPath string
	Padding int
}

var defaultAssets = []asset{
	{Name: "v", Text: "v", OutPath: "out/v.svg", Padding: 20},
	{Name: "vimcolorschemes", Text: "vimcolorschemes", OutPath: "out/vimcolorschemes.svg", Padding: 40},
}

// Generate writes all assets, or only the named assets when names is non-empty.
func Generate(names []string) error {
	selected := selectedAssets(names)
	if err := validateSelection(selected, defaultAssets); err != nil {
		return err
	}

	font, err := ansifonts.LoadFont("pressstart")
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	for _, item := range defaultAssets {
		if len(selected) > 0 && !selected[item.Name] {
			continue
		}
		if err := generateAsset(item, font); err != nil {
			return fmt.Errorf("generate %s: %w", item.Name, err)
		}
	}

	return nil
}

func selectedAssets(names []string) map[string]bool {
	if len(names) == 0 {
		return nil
	}

	selected := make(map[string]bool, len(names))
	for _, name := range names {
		selected[name] = true
	}
	return selected
}

func validateSelection(selected map[string]bool, assets []asset) error {
	if len(selected) == 0 {
		return nil
	}

	known := make(map[string]bool, len(assets))
	for _, item := range assets {
		known[item.Name] = true
	}

	for name := range selected {
		if !known[name] {
			return fmt.Errorf("unknown asset: %s", name)
		}
	}

	return nil
}

func generateAsset(item asset, font *ansifonts.Font) error {
	lines := ansifonts.RenderTextWithOptions(item.Text, font, renderOptions())
	cells, cols, rows, err := parseANSI(lines)
	if err != nil {
		return err
	}

	svg := renderSVG(item, cells, cols, rows)
	if err := os.MkdirAll(filepath.Dir(item.OutPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(item.OutPath, svg, 0o644); err != nil {
		return err
	}

	raster := renderRaster(item, cells, cols, rows)
	basePath := strings.TrimSuffix(item.OutPath, filepath.Ext(item.OutPath))
	if err := writePNG(basePath+".png", raster); err != nil {
		return err
	}
	return writeWebP(basePath+".webp", raster)
}

func renderOptions() ansifonts.RenderOptions {
	return ansifonts.RenderOptions{
		CharSpacing:            2,
		WordSpacing:            2,
		LineSpacing:            1,
		Alignment:              ansifonts.CenterAlign,
		TextColor:              startColor,
		GradientColor:          endColor,
		GradientDirection:      ansifonts.LeftRight,
		UseGradient:            true,
		ScaleFactor:            1,
		ShadowEnabled:          true,
		ShadowHorizontalOffset: 1,
		ShadowVerticalOffset:   1,
		ShadowStyle:            ansifonts.LightShade,
	}
}
