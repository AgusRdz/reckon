package cmd

import (
	"fmt"
	"os"

	"github.com/agusrdz/reckon/config"
	"github.com/agusrdz/reckon/index"
)

// Ignore adds .codeindex to the local or global gitignore.
// scope must be "local" or "global".
func Ignore(scope string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reckon: %v\n", err)
		os.Exit(1)
	}

	// If not specified, read from .codeindex.yml (or default to local)
	if scope == "" {
		cfg, _ := config.Load(cwd)
		scope = cfg.Gitignore
	}

	switch scope {
	case "global":
		path := globalGitignorePath()
		if path == "" {
			fmt.Fprintln(os.Stderr, "reckon: could not determine global gitignore path")
			os.Exit(1)
		}
		addToGitignoreFile(path)
		removeFromLocalGitignore(cwd)
		fmt.Printf("added %s to global gitignore: %s\n", index.Filename, path)
	case "local":
		ensureLocalGitignore(cwd)
		fmt.Printf("added %s to local .gitignore\n", index.Filename)
	default:
		fmt.Fprintf(os.Stderr, "reckon ignore: unknown scope %q (use --local or --global)\n", scope)
		os.Exit(1)
	}
}
