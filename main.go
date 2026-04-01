package main

import (
	"fmt"
	"os"

	"github.com/agusrdz/reckon/cmd"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		cmd.Help(version)
		return
	}

	switch os.Args[1] {
	case "hook":
		cmd.Root(version)
	case "--help", "help", "-h":
		cmd.Help(version)
	case "--version", "version":
		fmt.Printf("reckon %s\n", version)
	case "index":
		cmd.Index()
	case "stats":
		cmd.Stats()
	case "update":
		cmd.Update(version)
	case "init", "setup":
		cmd.Init(version)
	case "uninstall":
		cmd.Uninstall(version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\nrun 'reckon help' for usage\n", os.Args[1])
		os.Exit(1)
	}
}
