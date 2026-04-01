package extract

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

func init() {
	Register(&javaExtractor{})
}

type javaExtractor struct{}

func (j *javaExtractor) Extensions() []string {
	return []string{".java"}
}

var (
	javaClass    = regexp.MustCompile(`\bclass\s+([A-Za-z_][A-Za-z0-9_]*)`)
	javaIface    = regexp.MustCompile(`\binterface\s+([A-Za-z_][A-Za-z0-9_]*)`)
	javaEnum     = regexp.MustCompile(`\benum\s+([A-Za-z_][A-Za-z0-9_]*)`)
	javaMethod   = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
	javaModifier = regexp.MustCompile(`\b(?:public|private|protected|static|final|abstract|synchronized|native|default|strictfp)\b`)
)

func (j *javaExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		trimmed := strings.TrimSpace(text)

		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
			continue
		}

		if m := javaClass.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "class"})
			continue
		}
		if m := javaIface.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := javaEnum.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "enum"})
			continue
		}
		if javaModifier.MatchString(text) {
			if m := javaMethod.FindStringSubmatch(text); m != nil {
				name := m[1]
				switch name {
				case "if", "for", "while", "switch", "catch", "new", "return", "class", "interface", "enum":
					continue
				}
				symbols = append(symbols, Symbol{Name: name, File: file, Line: line, Kind: "method"})
			}
		}
	}
	return symbols
}
