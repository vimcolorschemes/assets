package main

import (
	"fmt"
	"os"

	"github.com/vimcolorschemes/assets/internal/generator"
)

func main() {
	if err := generator.Generate(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
