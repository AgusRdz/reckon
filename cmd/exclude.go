package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/agusrdz/reckon/config"
)

// Exclude manages skip_patterns in .codeindex.yml.
//   - no args: print usage
//   - --list: show active patterns (defaults + user-defined)
//   - --remove <pattern>: remove a user-defined pattern
//   - <pattern>: add pattern
func Exclude(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: reckon exclude <pattern>\n       reckon exclude --list\n       reckon exclude --remove <pattern>")
		os.Exit(1)
	}

	switch args[0] {
	case "--list":
		excludeList(cwd)
	case "--remove":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "reckon exclude --remove: pattern required")
			os.Exit(1)
		}
		excludeRemove(cwd, args[1])
	default:
		excludeAdd(cwd, args[0])
	}
}

func excludeList(dir string) {
	merged, _ := config.Load(dir)
	user, _ := config.LoadFile(dir)

	userSet := map[string]bool{}
	for _, p := range user.SkipPatterns {
		userSet[p] = true
	}

	fmt.Println("active skip patterns:")
	for _, p := range merged.SkipPatterns {
		if userSet[p] {
			fmt.Printf("  %s  %s\n", p, dim("(user-defined)"))
		} else {
			fmt.Printf("  %s  %s\n", p, dim("(default)"))
		}
	}
}

func excludeAdd(dir, pattern string) {
	cfg, err := config.LoadFile(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	for _, p := range cfg.SkipPatterns {
		if strings.TrimSpace(p) == pattern {
			fmt.Printf("%s is already in skip_patterns\n", pattern)
			return
		}
	}

	cfg.SkipPatterns = append(cfg.SkipPatterns, pattern)
	if err := config.SaveFile(dir, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("added %q to skip_patterns in .codeindex.yml\n", pattern)
}

func excludeRemove(dir, pattern string) {
	cfg, err := config.LoadFile(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	filtered := cfg.SkipPatterns[:0]
	found := false
	for _, p := range cfg.SkipPatterns {
		if strings.TrimSpace(p) == pattern {
			found = true
			continue
		}
		filtered = append(filtered, p)
	}

	if !found {
		fmt.Fprintf(os.Stderr, "%q not found in user-defined skip_patterns (cannot remove defaults)\n", pattern)
		os.Exit(1)
	}

	cfg.SkipPatterns = filtered
	if err := config.SaveFile(dir, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("removed %q from skip_patterns in .codeindex.yml\n", pattern)
}
