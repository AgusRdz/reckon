package walk

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agusrdz/reckon/config"
)

func defaultConfig() *config.Config {
	return &config.Config{
		SkipPatterns: []string{
			"**/*.test.ts",
			"**/*.spec.ts",
			"**/__mocks__/**",
			"**/node_modules/**",
			"**/bin/**",
			"**/obj/**",
			"**/dist/**",
			"**/.git/**",
		},
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll %q: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile %q: %v", path, err)
	}
}

func containsPath(files []string, suffix string) bool {
	for _, f := range files {
		if filepath.ToSlash(f) != "" {
			rel := filepath.ToSlash(f)
			if len(rel) >= len(suffix) && rel[len(rel)-len(suffix):] == suffix {
				return true
			}
		}
	}
	return false
}

func TestFilesReturnsExpectedFiles(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")
	mustWriteFile(t, filepath.Join(dir, "src", "app.ts"), "export class App {}")

	cfg := defaultConfig()
	files, err := Files(dir, cfg)
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}

	if !containsPath(files, "main.go") {
		t.Errorf("expected main.go in results, got %v", files)
	}
	if !containsPath(files, "app.ts") {
		t.Errorf("expected src/app.ts in results, got %v", files)
	}
}

func TestFilesSkipsGitDir(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, ".git", "config"), "[core]")
	mustWriteFile(t, filepath.Join(dir, "main.go"), "package main")

	cfg := defaultConfig()
	files, err := Files(dir, cfg)
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}

	for _, f := range files {
		if filepath.ToSlash(f) != "" {
			// .git is a hidden dir — skipped by the HasPrefix(".") check
			base := filepath.Base(filepath.Dir(f))
			if base == ".git" {
				t.Errorf("expected .git directory to be skipped, found %q", f)
			}
		}
	}
	if !containsPath(files, "main.go") {
		t.Errorf("expected main.go, got %v", files)
	}
}

func TestFilesSkipsNodeModules(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "node_modules", "lodash", "index.js"), "module.exports = {}")
	mustWriteFile(t, filepath.Join(dir, "src", "index.ts"), "export default {}")

	cfg := defaultConfig()
	files, err := Files(dir, cfg)
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}

	for _, f := range files {
		if containsPath([]string{f}, "node_modules/lodash/index.js") {
			t.Errorf("node_modules should be skipped, found %q", f)
		}
	}
	if !containsPath(files, "index.ts") {
		t.Errorf("expected src/index.ts, got %v", files)
	}
}

func TestFilesSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "src", "app.ts"), "export class App {}")
	mustWriteFile(t, filepath.Join(dir, "src", "app.test.ts"), "describe('App', () => {})")
	mustWriteFile(t, filepath.Join(dir, "src", "service.spec.ts"), "describe('Service', () => {})")

	cfg := defaultConfig()
	files, err := Files(dir, cfg)
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}

	for _, f := range files {
		if containsPath([]string{f}, "app.test.ts") {
			t.Errorf("*.test.ts should be skipped, found %q", f)
		}
		if containsPath([]string{f}, "service.spec.ts") {
			t.Errorf("*.spec.ts should be skipped, found %q", f)
		}
	}
	if !containsPath(files, "app.ts") {
		t.Errorf("expected app.ts, got %v", files)
	}
}

func TestFilesSkipsHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	// Hidden dir — walk skips dirs starting with "."
	mustWriteFile(t, filepath.Join(dir, ".hidden", "file.go"), "package hidden")
	mustWriteFile(t, filepath.Join(dir, "visible.go"), "package main")

	cfg := defaultConfig()
	files, err := Files(dir, cfg)
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}

	for _, f := range files {
		if containsPath([]string{f}, ".hidden/file.go") {
			t.Errorf("hidden dir should be skipped, found %q", f)
		}
	}
	if !containsPath(files, "visible.go") {
		t.Errorf("expected visible.go, got %v", files)
	}
}

func TestFilesNoSkipMatch(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "handler.go"), "package main")
	mustWriteFile(t, filepath.Join(dir, "util.py"), "def helper(): pass")
	mustWriteFile(t, filepath.Join(dir, "style.css"), "body { margin: 0 }")

	cfg := &config.Config{SkipPatterns: []string{}}
	files, err := Files(dir, cfg)
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}

	if !containsPath(files, "handler.go") {
		t.Errorf("expected handler.go, got %v", files)
	}
	if !containsPath(files, "util.py") {
		t.Errorf("expected util.py, got %v", files)
	}
	if !containsPath(files, "style.css") {
		t.Errorf("expected style.css, got %v", files)
	}
}
