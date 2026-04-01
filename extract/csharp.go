package extract

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

func init() {
	Register(&csExtractor{})
}

type csExtractor struct{}

func (c *csExtractor) Extensions() []string {
	return []string{".cs"}
}

var (
	csClass    = regexp.MustCompile(`\bclass\s+([A-Za-z_][A-Za-z0-9_]*)`)
	csIface    = regexp.MustCompile(`\binterface\s+([A-Za-z_][A-Za-z0-9_]*)`)
	csEnum     = regexp.MustCompile(`\benum\s+([A-Za-z_][A-Za-z0-9_]*)`)
	csStruct   = regexp.MustCompile(`\bstruct\s+([A-Za-z_][A-Za-z0-9_]*)`)
	csMethod   = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
	csModifier = regexp.MustCompile(`\b(?:public|private|protected|internal|static|virtual|override|abstract|async|sealed|extern|partial|new)\b`)
)

func (c *csExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		trimmed := strings.TrimSpace(text)

		// Skip comments and empty
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		if m := csClass.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "class"})
			continue
		}
		if m := csIface.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := csEnum.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "enum"})
			continue
		}
		if m := csStruct.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "struct"})
			continue
		}
		// Methods: line with an access modifier and an identifier followed by (
		if csModifier.MatchString(text) {
			if m := csMethod.FindStringSubmatch(text); m != nil {
				name := m[1]
				// Skip keywords that look like methods
				switch name {
				case "if", "for", "while", "switch", "catch", "new", "return", "using", "class", "interface", "enum", "struct", "namespace":
					continue
				}
				symbols = append(symbols, Symbol{Name: name, File: file, Line: line, Kind: "method"})
			}
		}
	}
	return symbols
}
