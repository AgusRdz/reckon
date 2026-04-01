package extract

import (
	"bufio"
	"bytes"
	"regexp"
)

func init() {
	Register(&rustExtractor{})
}

type rustExtractor struct{}

func (r *rustExtractor) Extensions() []string {
	return []string{".rs"}
}

var (
	rsFn     = regexp.MustCompile(`^\s*(?:pub\s+(?:\([^)]*\)\s+)?)?(?:async\s+)?fn\s+([A-Za-z_][A-Za-z0-9_]*)`)
	rsStruct = regexp.MustCompile(`^\s*(?:pub\s+(?:\([^)]*\)\s+)?)?struct\s+([A-Za-z_][A-Za-z0-9_]*)`)
	rsEnum   = regexp.MustCompile(`^\s*(?:pub\s+(?:\([^)]*\)\s+)?)?enum\s+([A-Za-z_][A-Za-z0-9_]*)`)
	rsTrait  = regexp.MustCompile(`^\s*(?:pub\s+(?:\([^)]*\)\s+)?)?trait\s+([A-Za-z_][A-Za-z0-9_]*)`)
	rsImpl   = regexp.MustCompile(`^\s*(?:pub\s+)?impl`)
)

func (r *rustExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	inImpl := false
	for scanner.Scan() {
		line++
		text := scanner.Text()

		if rsImpl.MatchString(text) {
			inImpl = true
		}
		// Rough impl block exit heuristic: closing brace at start of line
		if text == "}" {
			inImpl = false
		}

		if m := rsStruct.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "struct"})
			continue
		}
		if m := rsEnum.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "enum"})
			continue
		}
		if m := rsTrait.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := rsFn.FindStringSubmatch(text); m != nil {
			kind := "function"
			if inImpl {
				kind = "method"
			}
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: kind})
		}
	}
	return symbols
}
