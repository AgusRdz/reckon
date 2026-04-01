package detect

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// Extensions returns the set of unique file extensions present under dir.
func Extensions(dir string) map[string]bool {
	exts := map[string]bool{}
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error { //nolint
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			return nil
		}
		if ext := strings.ToLower(filepath.Ext(name)); ext != "" {
			exts[ext] = true
		}
		return nil
	})
	return exts
}
