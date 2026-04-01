package config

import (
	"os"
	"path/filepath"
	"testing"
)

func containsPattern(patterns []string, p string) bool {
	for _, pat := range patterns {
		if pat == p {
			return true
		}
	}
	return false
}

func TestLoadNoFile(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load with no file should not error, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// Verify defaults include key patterns
	defaultPatterns := []string{
		"**/node_modules/**",
		"**/*.test.ts",
		"**/*.spec.ts",
		"**/dist/**",
	}
	for _, p := range defaultPatterns {
		if !containsPattern(cfg.SkipPatterns, p) {
			t.Errorf("expected default pattern %q in skip patterns %v", p, cfg.SkipPatterns)
		}
	}
}

func TestLoadWithFile(t *testing.T) {
	dir := t.TempDir()

	yaml := `skip_patterns:
  - "**/*.generated.go"
  - "**/testdata/**"
`
	if err := os.WriteFile(filepath.Join(dir, ".codeindex.yml"), []byte(yaml), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Custom patterns should be appended
	if !containsPattern(cfg.SkipPatterns, "**/*.generated.go") {
		t.Errorf("expected **/*.generated.go in patterns, got %v", cfg.SkipPatterns)
	}
	if !containsPattern(cfg.SkipPatterns, "**/testdata/**") {
		t.Errorf("expected **/testdata/** in patterns, got %v", cfg.SkipPatterns)
	}

	// Defaults should still be present
	if !containsPattern(cfg.SkipPatterns, "**/node_modules/**") {
		t.Errorf("expected **/node_modules/** still present after merge, got %v", cfg.SkipPatterns)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()

	// Write invalid YAML
	if err := os.WriteFile(filepath.Join(dir, ".codeindex.yml"), []byte("skip_patterns: [unclosed"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.SkipPatterns) == 0 {
		t.Error("expected non-empty default skip patterns")
	}
}
