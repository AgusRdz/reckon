package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/agusrdz/reckon/config"
	"github.com/agusrdz/reckon/extract"
	"github.com/agusrdz/reckon/index"
	"github.com/agusrdz/reckon/metrics"
	"github.com/agusrdz/reckon/walk"
)

// BuildStats holds results from BuildIndex.
type BuildStats struct {
	Symbols   int
	Files     int
	Languages []string
}

// BuildIndex walks dir, extracts symbols, writes .codeindex, returns stats.
func BuildIndex(dir string) ([]extract.Symbol, BuildStats, error) {
	cfg, err := config.Load(dir)
	if err != nil {
		return nil, BuildStats{}, err
	}

	files, err := walk.Files(dir, cfg)
	if err != nil {
		return nil, BuildStats{}, err
	}

	extMap := map[string]extract.Extractor{}
	for _, e := range extract.All() {
		for _, ext := range e.Extensions() {
			extMap[ext] = e
		}
	}

	var symbols []extract.Symbol
	langSet := map[string]bool{}
	fileCount := 0

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		e, ok := extMap[ext]
		if !ok {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		rel, err := filepath.Rel(dir, file)
		if err != nil {
			rel = file
		}
		rel = filepath.ToSlash(rel)

		syms := e.Extract(rel, content)
		if len(syms) > 0 {
			symbols = append(symbols, syms...)
			fileCount++
			lang := extToLang(ext)
			langSet[lang] = true
		}
	}

	langs := make([]string, 0, len(langSet))
	for l := range langSet {
		langs = append(langs, l)
	}
	sort.Strings(langs)

	if len(symbols) > 0 {
		if err := index.Write(dir, symbols); err != nil {
			return symbols, BuildStats{}, err
		}
		ensureGitignore(dir, cfg)
		metrics.Record(dir, len(symbols), fileCount, langs)
	}

	return symbols, BuildStats{
		Symbols:   len(symbols),
		Files:     fileCount,
		Languages: langs,
	}, nil
}

// ensureGitignore adds .codeindex to either the local or global .gitignore.
func ensureGitignore(dir string, cfg *config.Config) {
	if cfg.Gitignore == "global" {
		ensureGlobalGitignore()
		removeFromLocalGitignore(dir)
	} else {
		ensureLocalGitignore(dir)
	}
}

// ensureLocalGitignore adds .codeindex to the project's .gitignore.
func ensureLocalGitignore(dir string) {
	addToGitignoreFile(filepath.Join(dir, ".gitignore"))
}

// ensureGlobalGitignore adds .codeindex to the user's global gitignore file.
func ensureGlobalGitignore() {
	path := globalGitignorePath()
	if path == "" {
		return
	}
	addToGitignoreFile(path)
}

// globalGitignorePath returns the path to the global gitignore file.
// It checks git config first, then falls back to OS-specific defaults.
func globalGitignorePath() string {
	// Try git config --global core.excludesFile
	out, err := exec.Command("git", "config", "--global", "core.excludesFile").Output()
	if err == nil {
		p := strings.TrimSpace(string(out))
		if p != "" {
			// Expand ~ manually since os.ExpandEnv won't handle it
			if strings.HasPrefix(p, "~/") {
				home, err := os.UserHomeDir()
				if err == nil {
					p = filepath.Join(home, p[2:])
				}
			}
			return p
		}
	}

	// Fall back to OS defaults
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(home, ".gitignore_global")
	}
	// macOS and Linux: prefer XDG location
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}
	return filepath.Join(xdgConfig, "git", "ignore")
}

// addToGitignoreFile appends .codeindex to path if not already present.
func addToGitignoreFile(path string) {
	entry := index.Filename

	data, err := os.ReadFile(path)
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(line) == entry {
				return
			}
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		prefix := ""
		if len(data) > 0 && data[len(data)-1] != '\n' {
			prefix = "\n"
		}
		fmt.Fprintf(f, "%s%s\n", prefix, entry)
	} else if os.IsNotExist(err) {
		// Create parent dirs if needed (important for ~/.config/git/ignore)
		_ = os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte(entry+"\n"), 0644)
	}
}

// removeFromLocalGitignore removes .codeindex from the project's .gitignore if present.
func removeFromLocalGitignore(dir string) {
	path := filepath.Join(dir, ".gitignore")
	entry := index.Filename

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	filtered := lines[:0]
	removed := false
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			removed = true
			continue
		}
		filtered = append(filtered, line)
	}
	if !removed {
		return
	}

	// Trim trailing blank lines added by the removal, but keep one final newline
	result := strings.Join(filtered, "\n")
	result = strings.TrimRight(result, "\n")
	if result != "" {
		result += "\n"
	}
	os.WriteFile(path, []byte(result), 0644)
}

func extToLang(ext string) string {
	switch ext {
	case ".go":
		return "Go"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".js", ".jsx":
		return "JavaScript"
	case ".cs":
		return "C#"
	case ".py":
		return "Python"
	case ".java":
		return "Java"
	case ".rs":
		return "Rust"
	case ".rb":
		return "Ruby"
	case ".php":
		return "PHP"
	default:
		return ext
	}
}

// formatNum formats an integer with comma separators (e.g. 1243 -> "1,243").
func formatNum(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, c)
	}
	return string(result)
}
