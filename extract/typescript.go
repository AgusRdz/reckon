package extract

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

func init() {
	Register(&tsExtractor{})
}

type tsExtractor struct{}

func (t *tsExtractor) Extensions() []string {
	return []string{".ts", ".tsx", ".js", ".jsx"}
}

var (
	tsClass  = regexp.MustCompile(`^\s*(?:export\s+(?:default\s+)?)?(?:abstract\s+)?class\s+([A-Za-z_$][A-Za-z0-9_$]*)`)
	tsIface  = regexp.MustCompile(`^\s*(?:export\s+)?interface\s+([A-Za-z_$][A-Za-z0-9_$]*)`)
	tsType   = regexp.MustCompile(`^\s*(?:export\s+)?type\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*[=<]`)
	tsEnum   = regexp.MustCompile(`^\s*(?:export\s+)?(?:const\s+)?enum\s+([A-Za-z_$][A-Za-z0-9_$]*)`)
	tsFunc   = regexp.MustCompile(`^\s*(?:export\s+(?:default\s+)?)?(?:async\s+)?function\s*\*?\s*([A-Za-z_$][A-Za-z0-9_$]*)`)
	tsArrow  = regexp.MustCompile(`^\s*(?:export\s+)?(?:const|let|var)\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*=\s*(?:async\s+)?(?:\([^)]*\)|[A-Za-z_$][A-Za-z0-9_$]*)\s*=>`)
	tsArrow2 = regexp.MustCompile(`^\s*(?:export\s+)?(?:const|let|var)\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*=\s*(?:async\s+)?function`)
	tsMethod = regexp.MustCompile(`^\s+(?:(?:public|private|protected|static|async|override|abstract|readonly|get|set)\s+)*([A-Za-z_$][A-Za-z0-9_$]*)\s*(?:<[^>]*>)?\s*\(`)
)

var tsKeywords = map[string]bool{
	"if": true, "else": true, "for": true, "while": true, "do": true,
	"switch": true, "case": true, "catch": true, "finally": true, "try": true,
	"return": true, "throw": true, "new": true, "delete": true, "typeof": true,
	"void": true, "await": true, "yield": true, "super": true, "import": true,
	"export": true, "from": true, "const": true, "let": true, "var": true,
	"class": true, "interface": true, "type": true, "enum": true, "function": true,
	"extends": true, "implements": true, "instanceof": true, "in": true, "of": true,
}

func (t *tsExtractor) Extract(file string, content []byte) []Symbol {
	var symbols []Symbol
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()

		if m := tsClass.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "class"})
			continue
		}
		if m := tsIface.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "interface"})
			continue
		}
		if m := tsType.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "type"})
			continue
		}
		if m := tsEnum.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "enum"})
			continue
		}
		if m := tsFunc.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "function"})
			continue
		}
		if m := tsArrow.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "function"})
			continue
		}
		if m := tsArrow2.FindStringSubmatch(text); m != nil {
			symbols = append(symbols, Symbol{Name: m[1], File: file, Line: line, Kind: "function"})
			continue
		}
		// Methods: indented lines
		if len(text) > 0 && (text[0] == ' ' || text[0] == '\t') {
			if m := tsMethod.FindStringSubmatch(text); m != nil {
				name := m[1]
				if !tsKeywords[name] && !strings.HasPrefix(strings.TrimSpace(text), "//") {
					symbols = append(symbols, Symbol{Name: name, File: file, Line: line, Kind: "method"})
				}
			}
		}
	}
	return symbols
}
