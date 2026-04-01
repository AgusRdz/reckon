package extract

import (
	"bufio"
	"bytes"
	"regexp"
)

func init() {
	Register(&goExtractor{})
}

type goExtractor struct{}

func (g *goExtractor) Extensions() []string {
	return []string{".go"}
}

var (
	goFunc     = regexp.MustCompile(`^func\s+([A-Z][A-Za-z0-9_]*|[a-z][A-Za-z0-9_]*)\s*\(`)
	goMethod   = regexp.MustCompile(`^func\s+\([^)]+\)\s+([A-Za-z][A-Za-z0-9_]*)\s*\(`)
	goStruct   = regexp.MustCompile(`^type\s+([A-Za-z][A-Za-z0-9_]*)\s+struct\b`)
	goIface    = regexp.MustCompile(`^type\s+([A-Za-z][A-Za-z0-9_]*)\s+interface\b`)
	goType     = regexp.MustCompile(`^type\s+([A-Za-z][A-Za-z0-9_]*)\s+`)
)

func (g *goExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()

		if m := goMethod.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "method"})
			continue
		}
		if m := goFunc.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "function"})
			continue
		}
		if m := goStruct.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "struct"})
			continue
		}
		if m := goIface.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := goType.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "type"})
		}
	}
	return symbols
}
