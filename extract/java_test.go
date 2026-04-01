package extract

import (
	"testing"
)

func TestJavaExtractorClass(t *testing.T) {
	e := &javaExtractor{}

	src := `package com.example;

public class Foo {
}

class Bar {
}
`
	syms := e.Extract("test.java", []byte(src))

	if !findSymbol(syms, "Foo", "class") {
		t.Errorf("expected Foo class, got %v", syms)
	}
	if !findSymbol(syms, "Bar", "class") {
		t.Errorf("expected Bar class, got %v", syms)
	}
}

func TestJavaExtractorInterface(t *testing.T) {
	e := &javaExtractor{}

	src := `public interface Bar {
    void doSomething();
}
`
	syms := e.Extract("test.java", []byte(src))

	if !findSymbol(syms, "Bar", "interface") {
		t.Errorf("expected Bar interface, got %v", syms)
	}
}

func TestJavaExtractorEnum(t *testing.T) {
	e := &javaExtractor{}

	src := `public enum Status {
    ACTIVE, INACTIVE
}
`
	syms := e.Extract("test.java", []byte(src))

	if !findSymbol(syms, "Status", "enum") {
		t.Errorf("expected Status enum, got %v", syms)
	}
}

func TestJavaExtractorMethods(t *testing.T) {
	e := &javaExtractor{}

	src := `public class MyService {
    public void doSomething(String input) {
    }

    private static int compute(int a, int b) {
        return a + b;
    }

    protected abstract String format(Object obj);
}
`
	syms := e.Extract("test.java", []byte(src))

	if !findSymbol(syms, "doSomething", "method") {
		t.Errorf("expected doSomething method, got %v", syms)
	}
	if !findSymbol(syms, "compute", "method") {
		t.Errorf("expected compute method, got %v", syms)
	}
	if !findSymbol(syms, "format", "method") {
		t.Errorf("expected format method, got %v", syms)
	}
}

func TestJavaExtractorCommentsSkipped(t *testing.T) {
	e := &javaExtractor{}

	src := `// public class NotAClass
/**
 * public interface NotAnInterface
 */
public class RealClass {
    public void realMethod() {}
}
`
	syms := e.Extract("test.java", []byte(src))

	if findSymbol(syms, "NotAClass", "class") {
		t.Errorf("NotAClass from comment should not be extracted")
	}
	if !findSymbol(syms, "RealClass", "class") {
		t.Errorf("expected RealClass, got %v", syms)
	}
}
