package index

import (
	"testing"

	"github.com/agusrdz/reckon/extract"
)

func TestWriteReadRoundtrip(t *testing.T) {
	dir := t.TempDir()

	want := []extract.Symbol{
		{Name: "LoginAsync", File: "src/controllers/UserController.cs", Line: 145, Kind: "method"},
		{Name: "AuthService", File: "src/services/AuthService.cs", Line: 1, Kind: "class"},
		{Name: "useAuthStore", File: "src/stores/auth.store.ts", Line: 12, Kind: "function"},
	}

	if err := Write(dir, want); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got, err := Read(dir)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d symbols, got %d", len(want), len(got))
	}

	for i, w := range want {
		g := got[i]
		if g.Name != w.Name || g.File != w.File || g.Line != w.Line || g.Kind != w.Kind {
			t.Errorf("symbol[%d] mismatch: want %+v, got %+v", i, w, g)
		}
	}
}

func TestWriteReadEmpty(t *testing.T) {
	dir := t.TempDir()

	if err := Write(dir, []extract.Symbol{}); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got, err := Read(dir)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestWriteReadAllKinds(t *testing.T) {
	dir := t.TempDir()

	kinds := []string{"class", "method", "function", "interface", "type", "enum", "struct", "const"}
	var symbols []extract.Symbol
	for i, k := range kinds {
		symbols = append(symbols, extract.Symbol{
			Name: "Sym" + k,
			File: "src/file.go",
			Line: i + 1,
			Kind: k,
		})
	}

	if err := Write(dir, symbols); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got, err := Read(dir)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(got) != len(symbols) {
		t.Fatalf("expected %d symbols, got %d", len(symbols), len(got))
	}

	for _, k := range kinds {
		found := false
		for _, s := range got {
			if s.Kind == k {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("kind %q not found after roundtrip", k)
		}
	}
}

func TestReadMissingFile(t *testing.T) {
	dir := t.TempDir()

	_, err := Read(dir)
	if err == nil {
		t.Error("expected error reading non-existent index, got nil")
	}
}

func TestWriteNilSymbols(t *testing.T) {
	dir := t.TempDir()

	if err := Write(dir, nil); err != nil {
		t.Fatalf("Write with nil symbols failed: %v", err)
	}

	got, err := Read(dir)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result for nil symbols, got %v", got)
	}
}
