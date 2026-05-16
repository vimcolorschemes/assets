package generator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const shadowOpacity = "0.42"

type cell struct {
	X       int
	Y       int
	Color   string
	Opacity string
}

func parseANSI(lines []string) ([]cell, int, int, error) {
	var cells []cell
	maxCols := 0

	for y, line := range lines {
		x := 0
		color := ""

		for i := 0; i < len(line); {
			if line[i] == '\x1b' {
				seqEnd := strings.IndexByte(line[i:], 'm')
				if seqEnd == -1 {
					return nil, 0, 0, fmt.Errorf("unterminated ANSI sequence on line %d", y+1)
				}

				params := line[i+2 : i+seqEnd]
				color = applyANSI(params, color)
				i += seqEnd + 1
				continue
			}

			r, size := utf8.DecodeRuneInString(line[i:])
			if r == utf8.RuneError && size == 0 {
				break
			}

			if r != ' ' && color != "" {
				opacity := ""
				if r == '░' {
					opacity = shadowOpacity
				}
				cells = append(cells, cell{X: x, Y: y, Color: color, Opacity: opacity})
			}

			x++
			i += size
		}

		if x > maxCols {
			maxCols = x
		}
	}

	return cells, maxCols, len(lines), nil
}

func applyANSI(params string, current string) string {
	parts := strings.Split(params, ";")
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "0":
			current = ""
		case "38":
			if i+4 < len(parts) && parts[i+1] == "2" {
				r, rOK := parseColorPart(parts[i+2])
				g, gOK := parseColorPart(parts[i+3])
				b, bOK := parseColorPart(parts[i+4])
				if rOK && gOK && bOK {
					current = fmt.Sprintf("#%02x%02x%02x", r, g, b)
				}
				i += 4
			}
		}
	}
	return current
}

func parseColorPart(value string) (int, bool) {
	color, err := strconv.Atoi(value)
	if err != nil || color < 0 || color > 255 {
		return 0, false
	}
	return color, true
}
