package extract

import (
	"testing"
)

func TestTSExtractorClass(t *testing.T) {
	e := &tsExtractor{}

	src := `class Foo {
  doSomething() {
  }
}

export class Bar {
  handleEvent(e: Event) {
  }
}

export default class Baz {
}
`
	syms := e.Extract("test.ts", []byte(src))

	if !findSymbol(syms, "Foo", "class") {
		t.Errorf("expected Foo class, got %v", syms)
	}
	if !findSymbol(syms, "Bar", "class") {
		t.Errorf("expected Bar class, got %v", syms)
	}
	if !findSymbol(syms, "Baz", "class") {
		t.Errorf("expected Baz class, got %v", syms)
	}
}

func TestTSExtractorInterface(t *testing.T) {
	e := &tsExtractor{}

	src := `interface IFoo {
  getName(): string;
}

export interface IBar {
  id: number;
}
`
	syms := e.Extract("test.ts", []byte(src))

	if !findSymbol(syms, "IFoo", "interface") {
		t.Errorf("expected IFoo interface, got %v", syms)
	}
	if !findSymbol(syms, "IBar", "interface") {
		t.Errorf("expected IBar interface, got %v", syms)
	}
}

func TestTSExtractorType(t *testing.T) {
	e := &tsExtractor{}

	src := `type MyType = string | number;

export type ID = string;
`
	syms := e.Extract("test.ts", []byte(src))

	if !findSymbol(syms, "MyType", "type") {
		t.Errorf("expected MyType type, got %v", syms)
	}
	if !findSymbol(syms, "ID", "type") {
		t.Errorf("expected ID type, got %v", syms)
	}
}

func TestTSExtractorEnum(t *testing.T) {
	e := &tsExtractor{}

	src := `enum Color {
  Red,
  Green,
  Blue,
}

export const enum Direction {
  Up,
  Down,
}
`
	syms := e.Extract("test.ts", []byte(src))

	if !findSymbol(syms, "Color", "enum") {
		t.Errorf("expected Color enum, got %v", syms)
	}
	if !findSymbol(syms, "Direction", "enum") {
		t.Errorf("expected Direction enum, got %v", syms)
	}
}

func TestTSExtractorFunction(t *testing.T) {
	e := &tsExtractor{}

	src := `function doSomething(x: number): void {
}

export function doOther(a: string): string {
  return a;
}

async function fetchData(): Promise<void> {
}

export async function loadUser(): Promise<User> {
  return {} as User;
}
`
	syms := e.Extract("test.ts", []byte(src))

	for _, name := range []string{"doSomething", "doOther", "fetchData", "loadUser"} {
		if !findSymbol(syms, name, "function") {
			t.Errorf("expected function %q, got %v", name, syms)
		}
	}
}

func TestTSExtractorArrowFunction(t *testing.T) {
	e := &tsExtractor{}

	src := `const handler = async () => {
};

const fn = function(x: number) {
  return x;
};

const simple = (a: string) => a.toUpperCase();
`
	syms := e.Extract("test.ts", []byte(src))

	for _, name := range []string{"handler", "fn", "simple"} {
		if !findSymbol(syms, name, "function") {
			t.Errorf("expected function %q, got %v", name, syms)
		}
	}
}

func TestTSExtractorMethod(t *testing.T) {
	e := &tsExtractor{}

	src := `class MyService {
  getData(): string[] {
    return [];
  }

  async saveRecord(record: Record): Promise<void> {
  }

  private helper() {
  }
}
`
	syms := e.Extract("test.ts", []byte(src))

	for _, name := range []string{"getData", "saveRecord", "helper"} {
		if !findSymbol(syms, name, "method") {
			t.Errorf("expected method %q, got %v", name, syms)
		}
	}
}

func TestTSExtractorKeywordsNotExtracted(t *testing.T) {
	e := &tsExtractor{}

	src := `class Foo {
  bar() {
    if (condition) {
      for (let i = 0; i < 10; i++) {
        while (true) {
          switch (x) {
          }
        }
      }
    }
  }
}
`
	syms := e.Extract("test.ts", []byte(src))

	keywords := []string{"if", "for", "while", "switch"}
	for _, kw := range keywords {
		if findSymbol(syms, kw, "method") {
			t.Errorf("keyword %q should not be extracted as method", kw)
		}
	}
}

func TestTSExtractorAbstractClass(t *testing.T) {
	e := &tsExtractor{}

	src := `abstract class BaseController {
  abstract handle(): void;
}
`
	syms := e.Extract("test.ts", []byte(src))

	if !findSymbol(syms, "BaseController", "class") {
		t.Errorf("expected BaseController abstract class, got %v", syms)
	}
}
