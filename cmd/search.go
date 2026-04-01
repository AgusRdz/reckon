package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/agusrdz/reckon/index"
	"github.com/agusrdz/reckon/metrics"
)

// Search looks up a pattern in .codeindex, logs the query, and prints matches.
func Search(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: reckon search <pattern>")
		os.Exit(1)
	}
	pattern := strings.Join(args, " ")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	symbols, err := index.Read(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: .codeindex not found — run 'reckon index' first\n")
		os.Exit(1)
	}

	lower := strings.ToLower(pattern)
	var results []metrics.Result
	for _, s := range symbols {
		if strings.Contains(strings.ToLower(s.Name), lower) {
			results = append(results, metrics.Result{
				Symbol: s.Name,
				File:   s.File,
				Line:   s.Line,
				Kind:   s.Kind,
			})
		}
	}

	metrics.RecordSearch(cwd, pattern, results)

	if len(results) == 0 {
		fmt.Printf("no matches for %q\n", pattern)
		return
	}

	for _, r := range results {
		fmt.Printf("%s\t%s\t%d\t%s\n", r.Symbol, r.File, r.Line, r.Kind)
	}
}
