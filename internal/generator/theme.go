package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/superstarryeyes/bit/ansifonts"
)

const themePath = "theme.toml"

type theme struct {
	Font            string                `toml:"font"`
	Background      string                `toml:"background"`
	LightBackground string                `toml:"light_background"`
	GradientStart   string                `toml:"gradient_start"`
	GradientMid     string                `toml:"gradient_mid"`
	GradientEnd     string                `toml:"gradient_end"`
	CharSpacing     int                   `toml:"char_spacing"`
	WordSpacing     int                   `toml:"word_spacing"`
	LineSpacing     int                   `toml:"line_spacing"`
	Shadow          bool                  `toml:"shadow"`
	ShadowOffsetX   int                   `toml:"shadow_offset_x"`
	ShadowOffsetY   int                   `toml:"shadow_offset_y"`
	ShadowStyle     string                `toml:"shadow_style"`
	Assets          map[string]assetTheme `toml:"assets"`
}

type assetTheme struct {
	Padding     int `toml:"padding"`
	BorderScale int `toml:"border_scale"`
	OffsetX     int `toml:"offset_x"`
	OffsetY     int `toml:"offset_y"`
}

func loadTheme() (theme, error) {
	path, err := resolveThemePath()
	if err != nil {
		return theme{}, err
	}

	var t theme
	if _, err := toml.DecodeFile(path, &t); err != nil {
		return theme{}, err
	}
	if err := t.validate(); err != nil {
		return theme{}, err
	}
	return t, nil
}

func resolveThemePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	start := dir

	for {
		path := filepath.Join(dir, themePath)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found from %s", themePath, start)
		}
		dir = parent
	}
}

func (t theme) validate() error {
	if t.Font == "" {
		return fmt.Errorf("theme font is required")
	}
	for name, value := range map[string]string{
		"background":       t.Background,
		"light_background": t.LightBackground,
		"gradient_start":   t.GradientStart,
		"gradient_mid":     t.GradientMid,
		"gradient_end":     t.GradientEnd,
	} {
		if _, err := parseHexRGB(value); err != nil {
			return fmt.Errorf("invalid %s color %q", name, value)
		}
	}
	if _, err := shadowStyle(t.ShadowStyle); err != nil {
		return err
	}
	return nil
}

func (t theme) renderOptions() ansifonts.RenderOptions {
	return ansifonts.RenderOptions{
		CharSpacing:            t.CharSpacing,
		WordSpacing:            t.WordSpacing,
		LineSpacing:            t.LineSpacing,
		Alignment:              ansifonts.CenterAlign,
		TextColor:              t.GradientStart,
		GradientColor:          t.GradientEnd,
		GradientDirection:      ansifonts.LeftRight,
		UseGradient:            true,
		ScaleFactor:            1,
		ShadowEnabled:          t.Shadow,
		ShadowHorizontalOffset: t.ShadowOffsetX,
		ShadowVerticalOffset:   t.ShadowOffsetY,
		ShadowStyle:            mustShadowStyle(t.ShadowStyle),
	}
}

func shadowStyle(value string) (ansifonts.ShadowStyle, error) {
	switch strings.ToLower(value) {
	case "light":
		return ansifonts.LightShade, nil
	case "medium":
		return ansifonts.MediumShade, nil
	case "dark":
		return ansifonts.DarkShade, nil
	default:
		return ansifonts.LightShade, fmt.Errorf("unknown shadow_style: %s", value)
	}
}

func mustShadowStyle(value string) ansifonts.ShadowStyle {
	style, err := shadowStyle(value)
	if err != nil {
		panic(err)
	}
	return style
}
