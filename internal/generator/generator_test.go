package generator

import (
	"strings"
	"testing"
)

func TestValidateSelectionRejectsUnknownAsset(t *testing.T) {
	err := validateSelection(selectedAssets([]string{"missing"}), defaultAssets)
	if err == nil {
		t.Fatal("validateSelection returned nil error for unknown asset")
	}
	if !strings.Contains(err.Error(), "unknown asset: missing") {
		t.Fatalf("error = %q, want unknown asset message", err.Error())
	}
}

func TestSelectedAssetsEmptyMeansAllAssets(t *testing.T) {
	selected := selectedAssets(nil)
	if selected != nil {
		t.Fatalf("selectedAssets(nil) = %#v, want nil", selected)
	}
}

func TestVariantPathAddsSuffixBeforeExtension(t *testing.T) {
	got := variantPath("out/v/v.svg", "borderless")
	if got != "out/v/v-borderless.svg" {
		t.Fatalf("variantPath() = %q, want out/v/v-borderless.svg", got)
	}
}
