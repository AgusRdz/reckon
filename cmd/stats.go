package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/agusrdz/reckon/index"
)

// Stats reads .codeindex and prints a breakdown by kind and language.
func Stats() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	symbols, err := index.Read(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reckon: .codeindex not found — run 'reckon index' first")
		os.Exit(1)
	}

	if len(symbols) == 0 {
		fmt.Println("no symbols in index")
		return
	}

	byKind := map[string]int{}
	byLang := map[string]int{}
	fileSet := map[string]bool{}

	for _, s := range symbols {
		byKind[s.Kind]++
		fileSet[s.File] = true
		ext := strings.ToLower(filepath.Ext(s.File))
		byLang[extToLang(ext)]++
	}

	fmt.Printf("total:    %s symbols across %s files\n\n",
		formatNum(len(symbols)),
		formatNum(len(fileSet)),
	)

	kinds := make([]string, 0, len(byKind))
	for k := range byKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	fmt.Println("by kind:")
	for _, k := range kinds {
		fmt.Printf("  %-12s %s\n", k, formatNum(byKind[k]))
	}

	langs := make([]string, 0, len(byLang))
	for l := range byLang {
		langs = append(langs, l)
	}
	sort.Strings(langs)

	fmt.Println("\nby language:")
	for _, l := range langs {
		fmt.Printf("  %-12s %s\n", l, formatNum(byLang[l]))
	}
}
