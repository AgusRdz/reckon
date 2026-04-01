package extract

import (
	"testing"
)

func TestCSExtractorClass(t *testing.T) {
	e := &csExtractor{}

	src := `using System;

namespace MyApp
{
    public class UserController
    {
    }

    internal class Helper
    {
    }
}
`
	syms := e.Extract("test.cs", []byte(src))

	if !findSymbol(syms, "UserController", "class") {
		t.Errorf("expected UserController class, got %v", syms)
	}
	if !findSymbol(syms, "Helper", "class") {
		t.Errorf("expected Helper class, got %v", syms)
	}
}

func TestCSExtractorInterface(t *testing.T) {
	e := &csExtractor{}

	src := `public interface IAuthService
{
    Task<User> LoginAsync(string username, string password);
}
`
	syms := e.Extract("test.cs", []byte(src))

	if !findSymbol(syms, "IAuthService", "interface") {
		t.Errorf("expected IAuthService interface, got %v", syms)
	}
}

func TestCSExtractorEnum(t *testing.T) {
	e := &csExtractor{}

	src := `public enum StatusCode
{
    OK = 200,
    NotFound = 404,
}
`
	syms := e.Extract("test.cs", []byte(src))

	if !findSymbol(syms, "StatusCode", "enum") {
		t.Errorf("expected StatusCode enum, got %v", syms)
	}
}

func TestCSExtractorStruct(t *testing.T) {
	e := &csExtractor{}

	src := `public struct Point
{
    public int X;
    public int Y;
}
`
	syms := e.Extract("test.cs", []byte(src))

	if !findSymbol(syms, "Point", "struct") {
		t.Errorf("expected Point struct, got %v", syms)
	}
}

func TestCSExtractorMethods(t *testing.T) {
	e := &csExtractor{}

	src := `public class UserController
{
    public async Task<User> LoginAsync(string username, string password)
    {
        return null;
    }

    private void HandleError(Exception ex)
    {
    }

    protected internal static int Compute(int a, int b)
    {
        return a + b;
    }
}
`
	syms := e.Extract("test.cs", []byte(src))

	if !findSymbol(syms, "LoginAsync", "method") {
		t.Errorf("expected LoginAsync method, got %v", syms)
	}
	if !findSymbol(syms, "HandleError", "method") {
		t.Errorf("expected HandleError method, got %v", syms)
	}
	if !findSymbol(syms, "Compute", "method") {
		t.Errorf("expected Compute method, got %v", syms)
	}
}

func TestCSExtractorCommentsNotExtracted(t *testing.T) {
	e := &csExtractor{}

	src := `// public class NotAClass
/* public interface NotAnInterface */
public class RealClass
{
    // private void NotAMethod()
    public void RealMethod()
    {
    }
}
`
	syms := e.Extract("test.cs", []byte(src))

	if findSymbol(syms, "NotAClass", "class") {
		t.Errorf("NotAClass from comment should not be extracted")
	}
	if findSymbol(syms, "NotAnInterface", "interface") {
		t.Errorf("NotAnInterface from comment should not be extracted")
	}
	if findSymbol(syms, "NotAMethod", "method") {
		t.Errorf("NotAMethod from comment should not be extracted")
	}
	if !findSymbol(syms, "RealClass", "class") {
		t.Errorf("expected RealClass, got %v", syms)
	}
	if !findSymbol(syms, "RealMethod", "method") {
		t.Errorf("expected RealMethod, got %v", syms)
	}
}

func TestCSExtractorNoModifierNoMethod(t *testing.T) {
	e := &csExtractor{}

	// Lines without access modifiers should not be extracted as methods
	src := `public class Foo
{
    public void ValidMethod()
    {
        someLocalCall();
        anotherCall(x);
    }
}
`
	syms := e.Extract("test.cs", []byte(src))

	if findSymbol(syms, "someLocalCall", "method") {
		t.Errorf("someLocalCall (no modifier) should not be extracted as method")
	}
	if findSymbol(syms, "anotherCall", "method") {
		t.Errorf("anotherCall (no modifier) should not be extracted as method")
	}
	if !findSymbol(syms, "ValidMethod", "method") {
		t.Errorf("expected ValidMethod, got %v", syms)
	}
}
