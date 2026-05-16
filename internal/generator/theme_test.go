package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func testTheme() theme {
	return theme{
		Font:            "pressstart",
		Background:      "#141719",
		LightBackground: "#e7e8e4",
		GradientStart:   "#938aa9",
		GradientMid:     "#a88e90",
		GradientEnd:     "#b6927b",
		CharSpacing:     2,
		WordSpacing:     2,
		LineSpacing:     1,
		Shadow:          true,
		ShadowOffsetX:   1,
		ShadowOffsetY:   1,
		ShadowStyle:     "light",
		Assets: map[string]assetTheme{
			"v":               {Padding: 20, OffsetX: -5, OffsetY: -5},
			"vimcolorschemes": {Padding: 40},
		},
	}
}

func TestThemeApplyOverridesAssetLayout(t *testing.T) {
	got := testTheme().apply(asset{Name: "v", Square: true})
	if got.Padding != 20 || got.OffsetX != -5 || got.OffsetY != -5 || !got.Square {
		t.Fatalf("theme.apply() = %#v, want themed padding/offsets with original fields kept", got)
	}
}

func TestThemeRejectsUnknownShadowStyle(t *testing.T) {
	theme := testTheme()
	theme.ShadowStyle = "heavy"
	if err := theme.validate(); err == nil {
		t.Fatal("validate returned nil error for unknown shadow style")
	}
}

func TestLoadThemeFindsThemeInParentDirectory(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatal(err)
		}
	})

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, themePath), []byte(`font = "pressstart"
background = "#141719"
light_background = "#e7e8e4"
gradient_start = "#938aa9"
gradient_mid = "#a88e90"
gradient_end = "#b6927b"
char_spacing = 2
word_spacing = 2
line_spacing = 1
shadow = true
shadow_offset_x = 1
shadow_offset_y = 1
shadow_style = "light"

[assets.v]
padding = 20
offset_x = -5
offset_y = -5
`), 0o644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(dir, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}

	theme, err := loadTheme()
	if err != nil {
		t.Fatal(err)
	}
	if theme.Font != "pressstart" || theme.Assets["v"].Padding != 20 {
		t.Fatalf("loadTheme() = %#v, want parent theme", theme)
	}
}
