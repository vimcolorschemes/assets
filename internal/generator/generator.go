package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/superstarryeyes/bit/ansifonts"
)

type asset struct {
	Name        string
	Text        string
	OutPath     string
	Padding     int
	Square      bool
	Border      bool
	BorderScale int
	Background  string
	Transparent bool
	OffsetX     int
	OffsetY     int
}

var defaultAssets = []asset{
	{Name: "v", Text: "v", OutPath: "out/v/v.svg", Square: true, Border: true},
	{Name: "vimcolorschemes", Text: "vimcolorschemes", OutPath: "out/vimcolorschemes/vimcolorschemes.svg", Border: true},
}

// Generate writes all assets, or only the named assets when names is non-empty.
func Generate(names []string) error {
	theme, err := loadTheme()
	if err != nil {
		return fmt.Errorf("load theme: %w", err)
	}

	selected := selectedAssets(names)
	if err := validateSelection(selected, defaultAssets); err != nil {
		return err
	}

	font, err := ansifonts.LoadFont(theme.Font)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	for _, item := range defaultAssets {
		if len(selected) > 0 && !selected[item.Name] {
			continue
		}
		if err := generateAsset(theme.apply(item), theme, font); err != nil {
			return fmt.Errorf("generate %s: %w", item.Name, err)
		}
	}

	return nil
}

func (t theme) apply(item asset) asset {
	assetTheme, ok := t.Assets[item.Name]
	if !ok {
		return item
	}
	item.Padding = assetTheme.Padding
	item.BorderScale = assetTheme.BorderScale
	item.OffsetX = assetTheme.OffsetX
	item.OffsetY = assetTheme.OffsetY
	return item
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

func generateAsset(item asset, theme theme, font *ansifonts.Font) error {
	lines := ansifonts.RenderTextWithOptions(item.Text, font, theme.renderOptions())
	cells, cols, rows, err := parseANSI(lines)
	if err != nil {
		return err
	}

	for _, variant := range assetVariants(item, theme) {
		if err := writeAssetFiles(variant, theme, cells, cols, rows); err != nil {
			return err
		}
	}
	return nil
}

func assetVariants(item asset, theme theme) []asset {
	item.Background = theme.Background
	borderless := item
	borderless.Name += " borderless"
	borderless.OutPath = variantPath(item.OutPath, "borderless")
	borderless.Border = false

	light := item
	light.Name += " light"
	light.OutPath = variantPath(item.OutPath, "light")
	light.Background = theme.LightBackground

	lightBorderless := light
	lightBorderless.Name += " borderless"
	lightBorderless.OutPath = variantPath(item.OutPath, "light-borderless")
	lightBorderless.Border = false

	transparent := item
	transparent.Name += " transparent"
	transparent.OutPath = variantPath(item.OutPath, "transparent")
	transparent.Transparent = true

	transparentBorderless := transparent
	transparentBorderless.Name += " borderless"
	transparentBorderless.OutPath = variantPath(item.OutPath, "transparent-borderless")
	transparentBorderless.Border = false

	return []asset{item, borderless, light, lightBorderless, transparent, transparentBorderless}
}

func writeAssetFiles(item asset, theme theme, cells []cell, cols int, rows int) error {
	svg := renderSVG(item, theme, cells, cols, rows)
	if err := os.MkdirAll(filepath.Dir(item.OutPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(item.OutPath, svg, 0o644); err != nil {
		return err
	}

	raster := renderRaster(item, theme, cells, cols, rows)
	basePath := strings.TrimSuffix(item.OutPath, filepath.Ext(item.OutPath))
	pngPath := basePath + ".png"
	webpPath := basePath + ".webp"
	if err := writePNG(pngPath, raster); err != nil {
		return err
	}
	if item.Transparent {
		return nil
	}
	return writeWebP(webpPath, raster)
}

func variantPath(path string, variant string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + "-" + variant + ext
}

func assetBackground(item asset, theme theme) string {
	if item.Background != "" {
		return item.Background
	}
	return theme.Background
}
