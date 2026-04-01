package cmd

import (
	"fmt"
	"os"
	"strings"
)

// Index rebuilds .codeindex and prints stats (no hook response).
func Index() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	symbols, stats, err := BuildIndex(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	if len(symbols) == 0 {
		fmt.Println("reckon: no symbols found")
		return
	}

	langs := strings.Join(stats.Languages, ", ")
	fmt.Printf("rebuilt .codeindex: %s symbols across %s files (%s)\n",
		formatNum(stats.Symbols),
		formatNum(stats.Files),
		langs,
	)
}
