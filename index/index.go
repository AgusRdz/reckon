package index

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/agusrdz/reckon/extract"
)

// Filename is the name of the symbol index file.
const Filename = ".codeindex"

// Write writes symbols to .codeindex in dir.
func Write(dir string, symbols []extract.Symbol) error {
	path := filepath.Join(dir, Filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, s := range symbols {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", s.Name, s.File, s.Line, s.Kind)
	}
	return w.Flush()
}

// Read reads .codeindex from dir and returns symbols.
func Read(dir string) ([]extract.Symbol, error) {
	path := filepath.Join(dir, Filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var symbols []extract.Symbol
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 4 {
			continue
		}
		line, _ := strconv.Atoi(parts[2])
		symbols = append(symbols, extract.Symbol{
			Name: parts[0],
			File: parts[1],
			Line: line,
			Kind: parts[3],
		})
	}
	return symbols, scanner.Err()
}
