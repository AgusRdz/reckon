package extract

import (
	"bufio"
	"bytes"
	"regexp"
)

func init() {
	Register(&rbExtractor{})
}

type rbExtractor struct{}

func (r *rbExtractor) Extensions() []string {
	return []string{".rb"}
}

var (
	rbClass  = regexp.MustCompile(`^\s*class\s+([A-Za-z_][A-Za-z0-9_:]*)`)
	rbModule = regexp.MustCompile(`^\s*module\s+([A-Za-z_][A-Za-z0-9_:]*)`)
	rbDef    = regexp.MustCompile(`^(\s*)def\s+(?:self\.)?([A-Za-z_][A-Za-z0-9_?!]*)`)
)

func (r *rbExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()

		if m := rbClass.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "class"})
			continue
		}
		if m := rbModule.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := rbDef.FindStringSubmatch(text); m != nil {
			indent := m[1]
			name := m[2]
			kind := "function"
			if len(indent) > 0 {
				kind = "method"
			}
			symbols = append(symbols, Symbol{Name: name, File: file, Line: line, Kind: kind})
		}
	}
	return symbols
}
