package generator

import "testing"

func TestParseANSITruecolorCells(t *testing.T) {
	lines := []string{
		"\x1b[38;2;1;2;3m█\x1b[0m \x1b[38;2;4;5;6m░\x1b[0m",
		"  \x1b[38;2;255;128;0m█\x1b[0m",
	}

	cells, cols, rows, err := parseANSI(lines)
	if err != nil {
		t.Fatalf("parseANSI returned error: %v", err)
	}
	if cols != 3 || rows != 2 {
		t.Fatalf("dimensions = %dx%d, want 3x2", cols, rows)
	}

	want := []cell{
		{X: 0, Y: 0, Color: "#010203"},
		{X: 2, Y: 0, Color: "#040506", Opacity: shadowOpacity},
		{X: 2, Y: 1, Color: "#ff8000"},
	}
	if len(cells) != len(want) {
		t.Fatalf("got %d cells, want %d", len(cells), len(want))
	}
	for i := range want {
		if cells[i] != want[i] {
			t.Fatalf("cell %d = %#v, want %#v", i, cells[i], want[i])
		}
	}
}

func TestParseANSIRejectsUnterminatedSequence(t *testing.T) {
	_, _, _, err := parseANSI([]string{"\x1b[38;2;1;2;3"})
	if err == nil {
		t.Fatal("parseANSI returned nil error for unterminated ANSI sequence")
	}
}
