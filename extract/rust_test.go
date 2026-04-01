package extract

import (
	"testing"
)

func TestRustExtractorFunctions(t *testing.T) {
	e := &rustExtractor{}

	src := `fn top_level_fn(x: i32) -> i32 {
    x + 1
}

pub fn public_fn() {
}

async fn async_fn() -> Result<(), Error> {
    Ok(())
}
`
	syms := e.Extract("test.rs", []byte(src))

	for _, name := range []string{"top_level_fn", "public_fn", "async_fn"} {
		if !findSymbol(syms, name, "function") {
			t.Errorf("expected function %q, got %v", name, syms)
		}
	}
}

func TestRustExtractorStruct(t *testing.T) {
	e := &rustExtractor{}

	src := `struct MyStruct {
    field: i32,
}

pub struct PublicStruct {
    pub value: String,
}
`
	syms := e.Extract("test.rs", []byte(src))

	if !findSymbol(syms, "MyStruct", "struct") {
		t.Errorf("expected MyStruct struct, got %v", syms)
	}
	if !findSymbol(syms, "PublicStruct", "struct") {
		t.Errorf("expected PublicStruct struct, got %v", syms)
	}
}

func TestRustExtractorEnum(t *testing.T) {
	e := &rustExtractor{}

	src := `enum MyEnum {
    VariantA,
    VariantB(i32),
}

pub enum Status {
    Active,
    Inactive,
}
`
	syms := e.Extract("test.rs", []byte(src))

	if !findSymbol(syms, "MyEnum", "enum") {
		t.Errorf("expected MyEnum enum, got %v", syms)
	}
	if !findSymbol(syms, "Status", "enum") {
		t.Errorf("expected Status enum, got %v", syms)
	}
}

func TestRustExtractorTrait(t *testing.T) {
	e := &rustExtractor{}

	src := `trait MyTrait {
    fn do_something(&self);
}

pub trait Drawable {
    fn draw(&self);
}
`
	syms := e.Extract("test.rs", []byte(src))

	if !findSymbol(syms, "MyTrait", "interface") {
		t.Errorf("expected MyTrait as interface kind, got %v", syms)
	}
	if !findSymbol(syms, "Drawable", "interface") {
		t.Errorf("expected Drawable as interface kind, got %v", syms)
	}
}

func TestRustExtractorImplMethod(t *testing.T) {
	e := &rustExtractor{}

	src := `struct Point {
    x: f64,
    y: f64,
}

impl Point {
    fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }

    pub fn distance(&self, other: &Point) -> f64 {
        0.0
    }
}
`
	syms := e.Extract("test.rs", []byte(src))

	if !findSymbol(syms, "Point", "struct") {
		t.Errorf("expected Point struct, got %v", syms)
	}
	if !findSymbol(syms, "new", "method") {
		t.Errorf("expected new as method (inside impl), got %v", syms)
	}
	if !findSymbol(syms, "distance", "method") {
		t.Errorf("expected distance as method (inside impl), got %v", syms)
	}
}
