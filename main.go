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
		// No subcommand: show help if stdout is a terminal (interactive),
		// otherwise hook mode (stdout is piped to Claude Code).
		if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
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
