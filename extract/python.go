package extract

import (
	"bufio"
	"bytes"
	"regexp"
)

func init() {
	Register(&pyExtractor{})
}

type pyExtractor struct{}

func (p *pyExtractor) Extensions() []string {
	return []string{".py"}
}

var (
	pyClass = regexp.MustCompile(`^class\s+([A-Za-z_][A-Za-z0-9_]*)`)
	pyDef   = regexp.MustCompile(`^(\s*)def\s+([A-Za-z_][A-Za-z0-9_]*)`)
)

func (p *pyExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()

		if m := pyClass.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "class"})
			continue
		}
		if m := pyDef.FindStringSubmatch(text); m != nil {
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
