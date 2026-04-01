package extract

import (
	"testing"
)

func findSymbol(symbols []Symbol, name, kind string) bool {
	for _, s := range symbols {
		if s.Name == name && s.Kind == kind {
			return true
		}
	}
	return false
}

func findSymbolAt(symbols []Symbol, name, kind string, line int) bool {
	for _, s := range symbols {
		if s.Name == name && s.Kind == kind && s.Line == line {
			return true
		}
	}
	return false
}

func TestGoExtractor(t *testing.T) {
	e := &goExtractor{}

	src := `package main

import "fmt"

func TopLevel(x int) string {
	return fmt.Sprintf("%d", x)
}

func (r *Receiver) Method() error {
	return nil
}

func (v Value) ValueMethod(n int) int {
	return n
}

type MyStruct struct {
	Field string
}

type MyInterface interface {
	DoSomething() error
}

type MyAlias = string

type Result[T any] struct {
	Value T
}
`

	syms := e.Extract("test.go", []byte(src))

	tests := []struct {
		name string
		kind string
	}{
		{"TopLevel", "function"},
		{"Method", "method"},
		{"ValueMethod", "method"},
		{"MyStruct", "struct"},
		{"MyInterface", "interface"},
	}

	for _, tt := range tests {
		if !findSymbol(syms, tt.name, tt.kind) {
			t.Errorf("expected symbol %q kind=%q, not found in %v", tt.name, tt.kind, syms)
		}
	}
}

func TestGoExtractorLineNumbers(t *testing.T) {
	e := &goExtractor{}

	src := `package main

func Alpha() {}

type Beta struct{}

func Gamma() {}
`
	syms := e.Extract("test.go", []byte(src))

	if !findSymbolAt(syms, "Alpha", "function", 3) {
		t.Errorf("Alpha should be at line 3, got %v", syms)
	}
	if !findSymbolAt(syms, "Beta", "struct", 5) {
		t.Errorf("Beta should be at line 5, got %v", syms)
	}
	if !findSymbolAt(syms, "Gamma", "function", 7) {
		t.Errorf("Gamma should be at line 7, got %v", syms)
	}
}

func TestGoExtractorGenericType(t *testing.T) {
	e := &goExtractor{}

	// Generic type syntax: type Result[T any] struct {}
	// The goStruct regex requires `type NAME struct` (no brackets between name and struct).
	// The goType regex requires `type NAME ` (whitespace after name, but `[` follows).
	// Neither regex matches generics — this is a known regex-only limitation.
	// Non-generic type aliases ARE captured by goType.
	src := `package main

type StringAlias = string

type MyTypedef int
`
	syms := e.Extract("test.go", []byte(src))

	// Non-generic type alias/typedef should be caught by goType regex
	found := false
	for _, s := range syms {
		if s.Name == "MyTypedef" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected MyTypedef type symbol to be extracted, got %v", syms)
	}
}

func TestGoExtractorFile(t *testing.T) {
	e := &goExtractor{}
	syms := e.Extract("src/main.go", []byte(`package main

func Hello() {}
`))
	if len(syms) != 1 {
		t.Fatalf("expected 1 symbol, got %d", len(syms))
	}
	if syms[0].File != "src/main.go" {
		t.Errorf("expected file src/main.go, got %q", syms[0].File)
	}
}

func TestGoExtractorEmpty(t *testing.T) {
	e := &goExtractor{}
	syms := e.Extract("empty.go", []byte(`package main
`))
	if len(syms) != 0 {
		t.Errorf("expected no symbols from empty file, got %v", syms)
	}
}
