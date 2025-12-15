package filex

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func PathIsDir(s string) bool {
	return len(s) >= 1 && (s[len(s)-1] == '\\' || s[len(s)-1] == '/')
}

func List(root string, matches ...func(string, fs.DirEntry) error) (files []string, err error) {
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, match := range matches {
			if err = match(p, d); err != nil {
				return err
			}
		}

		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}

		files = append(files, strings.ReplaceAll(rel, "\\", "/"))
		return nil
	})
	return
}
