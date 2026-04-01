package extract

import (
	"testing"
)

func TestRubyExtractorClass(t *testing.T) {
	e := &rbExtractor{}

	src := `class MyClass
  def initialize
  end
end

class AnotherClass < BaseClass
end
`
	syms := e.Extract("test.rb", []byte(src))

	if !findSymbol(syms, "MyClass", "class") {
		t.Errorf("expected MyClass class, got %v", syms)
	}
	if !findSymbol(syms, "AnotherClass", "class") {
		t.Errorf("expected AnotherClass class, got %v", syms)
	}
}

func TestRubyExtractorModule(t *testing.T) {
	e := &rbExtractor{}

	src := `module MyModule
  def helper
  end
end
`
	syms := e.Extract("test.rb", []byte(src))

	if !findSymbol(syms, "MyModule", "interface") {
		t.Errorf("expected MyModule as interface kind, got %v", syms)
	}
}

func TestRubyExtractorTopLevelDef(t *testing.T) {
	e := &rbExtractor{}

	src := `def top_level_method
  puts "hello"
end

def another_method(x, y)
  x + y
end
`
	syms := e.Extract("test.rb", []byte(src))

	if !findSymbol(syms, "top_level_method", "function") {
		t.Errorf("expected top_level_method as function, got %v", syms)
	}
	if !findSymbol(syms, "another_method", "function") {
		t.Errorf("expected another_method as function, got %v", syms)
	}
}

func TestRubyExtractorInstanceMethod(t *testing.T) {
	e := &rbExtractor{}

	src := `class MyClass
  def instance_method
  end

  def another_instance_method(arg)
  end
end
`
	syms := e.Extract("test.rb", []byte(src))

	if !findSymbol(syms, "instance_method", "method") {
		t.Errorf("expected instance_method as method, got %v", syms)
	}
	if !findSymbol(syms, "another_instance_method", "method") {
		t.Errorf("expected another_instance_method as method, got %v", syms)
	}
}

func TestRubyExtractorClassMethod(t *testing.T) {
	e := &rbExtractor{}

	src := `class MyClass
  def self.class_method
  end
end

def self.top_level_class
end
`
	syms := e.Extract("test.rb", []byte(src))

	// Inside class: indented self.method → method kind
	if !findSymbol(syms, "class_method", "method") {
		t.Errorf("expected class_method as method, got %v", syms)
	}
}
