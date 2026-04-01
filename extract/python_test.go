package extract

import (
	"testing"
)

func TestPyExtractorClass(t *testing.T) {
	e := &pyExtractor{}

	src := `class MyClass:
    pass

class AnotherClass(BaseClass):
    pass
`
	syms := e.Extract("test.py", []byte(src))

	if !findSymbol(syms, "MyClass", "class") {
		t.Errorf("expected MyClass class, got %v", syms)
	}
	if !findSymbol(syms, "AnotherClass", "class") {
		t.Errorf("expected AnotherClass class, got %v", syms)
	}
}

func TestPyExtractorTopLevelFunction(t *testing.T) {
	e := &pyExtractor{}

	src := `def top_level_func(x, y):
    return x + y

def another_func():
    pass
`
	syms := e.Extract("test.py", []byte(src))

	if !findSymbol(syms, "top_level_func", "function") {
		t.Errorf("expected top_level_func as function, got %v", syms)
	}
	if !findSymbol(syms, "another_func", "function") {
		t.Errorf("expected another_func as function, got %v", syms)
	}
}

func TestPyExtractorMethod(t *testing.T) {
	e := &pyExtractor{}

	src := `class MyService:
    def method(self, x):
        pass

    def another_method(self):
        pass

    def __init__(self):
        pass
`
	syms := e.Extract("test.py", []byte(src))

	for _, name := range []string{"method", "another_method", "__init__"} {
		if !findSymbol(syms, name, "method") {
			t.Errorf("expected %q as method, got %v", name, syms)
		}
	}
}

func TestPyExtractorNestedClass(t *testing.T) {
	e := &pyExtractor{}

	src := `class Outer:
    class Inner:
        def inner_method(self):
            pass

    def outer_method(self):
        pass
`
	syms := e.Extract("test.py", []byte(src))

	if !findSymbol(syms, "Outer", "class") {
		t.Errorf("expected Outer class, got %v", syms)
	}
	// Inner is indented — pyClass requires `^class` so it won't match
	// This tests realistic extractor behavior (regex-only, no AST)
	if !findSymbol(syms, "inner_method", "method") {
		t.Errorf("expected inner_method as method, got %v", syms)
	}
	if !findSymbol(syms, "outer_method", "method") {
		t.Errorf("expected outer_method as method, got %v", syms)
	}
}

func TestPyExtractorFunctionVsMethod(t *testing.T) {
	e := &pyExtractor{}

	src := `def standalone():
    pass

class MyClass:
    def instance_method(self):
        pass
`
	syms := e.Extract("test.py", []byte(src))

	if !findSymbol(syms, "standalone", "function") {
		t.Errorf("expected standalone as function kind, got %v", syms)
	}
	if !findSymbol(syms, "instance_method", "method") {
		t.Errorf("expected instance_method as method kind, got %v", syms)
	}
}

func TestPyExtractorLineNumbers(t *testing.T) {
	e := &pyExtractor{}

	src := `class Alpha:
    pass

def beta():
    pass
`
	syms := e.Extract("test.py", []byte(src))

	if !findSymbolAt(syms, "Alpha", "class", 1) {
		t.Errorf("Alpha should be at line 1, got %v", syms)
	}
	if !findSymbolAt(syms, "beta", "function", 4) {
		t.Errorf("beta should be at line 4, got %v", syms)
	}
}
