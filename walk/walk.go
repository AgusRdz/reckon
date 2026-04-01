package walk

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/agusrdz/reckon/config"
)

// Files returns all files under dir that are not excluded by cfg.
func Files(dir string, cfg *config.Config) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()

		if d.IsDir() {
			if strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			if skipDir(name, cfg.SkipPatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)
		if !skipFile(rel, cfg.SkipPatterns) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func skipDir(name string, patterns []string) bool {
	for _, p := range patterns {
		if strings.HasPrefix(p, "**/") && strings.HasSuffix(p, "/**") {
			dirName := p[3 : len(p)-3]
			if name == dirName {
				return true
			}
		}
	}
	return false
}

func skipFile(rel string, patterns []string) bool {
	base := filepath.Base(rel)
	for _, p := range patterns {
		if strings.HasPrefix(p, "**/*.") {
			// e.g. "**/*.test.ts" → suffix = ".test.ts"
			suffix := p[4:] // skip "**/*"
			if strings.HasSuffix(rel, suffix) {
				return true
			}
		}
		if !strings.Contains(p, "/") {
			if matched, _ := filepath.Match(p, base); matched {
				return true
			}
		}
		if matched, _ := filepath.Match(p, rel); matched {
			return true
		}
	}
	return false
}
