package extract

import (
	"testing"
)

func TestPHPExtractorClass(t *testing.T) {
	e := &phpExtractor{}

	src := `<?php

class Foo {
}

abstract class AbstractBase {
}
`
	syms := e.Extract("test.php", []byte(src))

	if !findSymbol(syms, "Foo", "class") {
		t.Errorf("expected Foo class, got %v", syms)
	}
	if !findSymbol(syms, "AbstractBase", "class") {
		t.Errorf("expected AbstractBase class, got %v", syms)
	}
}

func TestPHPExtractorInterface(t *testing.T) {
	e := &phpExtractor{}

	src := `<?php

interface IFoo {
    public function doSomething(): void;
}
`
	syms := e.Extract("test.php", []byte(src))

	if !findSymbol(syms, "IFoo", "interface") {
		t.Errorf("expected IFoo interface, got %v", syms)
	}
}

func TestPHPExtractorTrait(t *testing.T) {
	e := &phpExtractor{}

	src := `<?php

trait Reusable {
    public function sharedMethod(): void {
    }
}
`
	syms := e.Extract("test.php", []byte(src))

	// PHP extractor maps trait to "interface" kind
	if !findSymbol(syms, "Reusable", "interface") {
		t.Errorf("expected Reusable trait as interface kind, got %v", syms)
	}
}

func TestPHPExtractorStandaloneFunction(t *testing.T) {
	e := &phpExtractor{}

	src := `<?php

function standalone(int $x): int {
    return $x * 2;
}
`
	syms := e.Extract("test.php", []byte(src))

	if !findSymbol(syms, "standalone", "function") {
		t.Errorf("expected standalone function, got %v", syms)
	}
}

func TestPHPExtractorMethods(t *testing.T) {
	e := &phpExtractor{}

	src := `<?php

class MyController {
    public function doSomething(): void {
    }

    private static function helper(string $s): string {
        return $s;
    }

    protected function process(): bool {
        return true;
    }
}
`
	syms := e.Extract("test.php", []byte(src))

	if !findSymbol(syms, "doSomething", "method") {
		t.Errorf("expected doSomething method, got %v", syms)
	}
	if !findSymbol(syms, "helper", "method") {
		t.Errorf("expected helper method, got %v", syms)
	}
	if !findSymbol(syms, "process", "method") {
		t.Errorf("expected process method, got %v", syms)
	}
}

func TestPHPExtractorCommentsSkipped(t *testing.T) {
	e := &phpExtractor{}

	src := `<?php

// class NotAClass {}
/* interface NotAnInterface */
class RealClass {
    public function realMethod(): void {}
}
`
	syms := e.Extract("test.php", []byte(src))

	if findSymbol(syms, "NotAClass", "class") {
		t.Errorf("NotAClass from comment should not be extracted")
	}
	if !findSymbol(syms, "RealClass", "class") {
		t.Errorf("expected RealClass, got %v", syms)
	}
}
