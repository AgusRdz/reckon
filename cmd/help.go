package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agusrdz/reckon/hooks"
)

// Help prints usage information.
func Help(version string) {
	const colW = 32
	section := func(name string) string { return bold(cyan(name)) + "\n" }
	row := func(cmd, desc string) string {
		return fmt.Sprintf("  %-*s%s\n", colW, cmd, dim(desc))
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s %s — symbol index for Claude Code\n\n", bold("reckon"), version))

	b.WriteString(bold("Usage") + "\n")
	b.WriteString(row("reckon", "SessionStart hook: rebuild index, inject pointer"))
	b.WriteString(row("reckon <subcommand>", "Run a management subcommand"))
	b.WriteString("\n")

	b.WriteString(section("Setup"))
	b.WriteString(row("init", "Install Claude Code SessionStart hook"))
	b.WriteString(row("ignore [--local|--global]", "Add .codeindex to local or global gitignore"))
	b.WriteString(row("update", "Update reckon to the latest release"))
	b.WriteString(row("uninstall", "Remove hook and config"))
	b.WriteString("\n")

	b.WriteString(section("Index"))
	b.WriteString(row("index", "Rebuild .codeindex, print stats"))
	b.WriteString(row("search <pattern>", "Search .codeindex for a symbol (logged for metrics)"))
	b.WriteString(row("stats", "Show symbol count, file count, language breakdown"))
	b.WriteString(row("exclude <pattern>", "Add a glob pattern to skip_patterns"))
	b.WriteString(row("exclude --list", "Show all active skip patterns"))
	b.WriteString(row("exclude --remove <pattern>", "Remove a user-defined skip pattern"))
	b.WriteString("\n")

	b.WriteString(section("Metrics"))
	b.WriteString(row("metrics", "Show build history: total builds, per-project stats"))
	b.WriteString(row("metrics --clear", "Clear the usage log"))
	b.WriteString("\n")

	b.WriteString(section("Other"))
	b.WriteString(row("version", "Show version"))
	b.WriteString(row("help", "Show this help"))

	fmt.Print(b.String())
}

// Init installs the SessionStart hook.
func Init(version string) {
	hooks.Install(version)
}

// Uninstall removes the hook and config.
func Uninstall(version string) {
	hooks.Uninstall()
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "reckon")
	os.RemoveAll(configDir)
	fmt.Println("reckon uninstalled")
	fmt.Println("  hook removed from ~/.claude/settings.json")
	fmt.Printf("  config removed: %s\n", configDir)
	fmt.Println("\nbinary not removed — delete manually or via your package manager")
}
