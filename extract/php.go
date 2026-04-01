package extract

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

func init() {
	Register(&phpExtractor{})
}

type phpExtractor struct{}

func (p *phpExtractor) Extensions() []string {
	return []string{".php"}
}

var (
	phpClass    = regexp.MustCompile(`\bclass\s+([A-Za-z_][A-Za-z0-9_]*)`)
	phpIface    = regexp.MustCompile(`\binterface\s+([A-Za-z_][A-Za-z0-9_]*)`)
	phpTrait    = regexp.MustCompile(`\btrait\s+([A-Za-z_][A-Za-z0-9_]*)`)
	phpFunc     = regexp.MustCompile(`^\s*function\s+([A-Za-z_][A-Za-z0-9_]*)`)
	phpMethod   = regexp.MustCompile(`^\s+(?:(?:public|private|protected|static|abstract|final)\s+)*function\s+([A-Za-z_][A-Za-z0-9_]*)`)
	phpModifier = regexp.MustCompile(`\b(?:public|private|protected|static|abstract|final)\b`)
)

func (p *phpExtractor) Extract(file string, content []byte) []Symbol {
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

		if m := phpClass.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "class"})
			continue
		}
		if m := phpIface.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := phpTrait.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if phpModifier.MatchString(text) {
			if m := phpMethod.FindStringSubmatch(text); m != nil {
				symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "method"})
				continue
			}
		}
		if m := phpFunc.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "function"})
		}
	}
	return symbols
}
