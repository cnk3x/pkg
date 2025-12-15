package filex

import (
	"io/fs"
	"path/filepath"
)

func Files(dir string, match func(string) bool) ([]string, error) {
	var files []string
	err := Walk(dir, match, func(fullPath string) error { files = append(files, fullPath); return nil })
	return files, err
}

func Walk(dir string, match func(string) bool, walkFn func(fullPath string) error) error {
	return filepath.WalkDir(dir, func(fullPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			fullPath += string(filepath.Separator)
		}

		if match(fullPath) {
			return walkFn(fullPath)
		}

		if d.IsDir() {
			return fs.SkipDir
		}

		return nil
	})
}

func PathIsDir(s string) bool {
	return len(s) >= 1 && (s[len(s)-1] == '\\' || s[len(s)-1] == '/')
}
