package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/agusrdz/reckon/config"
	"github.com/agusrdz/reckon/extract"
	"github.com/agusrdz/reckon/index"
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

	if len(symbols) > 0 {
		if err := index.Write(dir, symbols); err != nil {
			return symbols, BuildStats{}, err
		}
		ensureGitignore(dir)
	}

	langs := make([]string, 0, len(langSet))
	for l := range langSet {
		langs = append(langs, l)
	}
	sort.Strings(langs)

	return symbols, BuildStats{
		Symbols:   len(symbols),
		Files:     fileCount,
		Languages: langs,
	}, nil
}

// ensureGitignore adds .codeindex to the project's .gitignore if not already present.
func ensureGitignore(dir string) {
	path := filepath.Join(dir, ".gitignore")
	entry := index.Filename

	data, err := os.ReadFile(path)
	if err == nil {
		// File exists — check if entry is already there
		for _, line := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(line) == entry {
				return
			}
		}
		// Append with a trailing newline
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		// Ensure we start on a new line
		suffix := ""
		if len(data) > 0 && data[len(data)-1] != '\n' {
			suffix = "\n"
		}
		fmt.Fprintf(f, "%s%s\n", suffix, entry)
	} else if os.IsNotExist(err) {
		// No .gitignore — create one
		os.WriteFile(path, []byte(entry+"\n"), 0644)
	}
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
