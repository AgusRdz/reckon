package main

import (
	"fmt"
	"os"

	"github.com/agusrdz/reckon/cmd"
	"github.com/mattn/go-isatty"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		// No subcommand: hook mode if stdin is a pipe, else help.
		if isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			cmd.Help(version)
			return
		}
		cmd.Root(version)
		return
	}

	switch os.Args[1] {
	case "--help", "help", "-h":
		cmd.Help(version)
	case "--version", "version":
		fmt.Printf("reckon %s\n", version)
	case "index":
		cmd.Index()
	case "stats":
		cmd.Stats()
	case "init", "setup":
		cmd.Init(version)
	case "uninstall":
		cmd.Uninstall(version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\nrun 'reckon help' for usage\n", os.Args[1])
		os.Exit(1)
	}
}
