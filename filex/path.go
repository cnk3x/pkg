package filex

import (
	"os"
	"path/filepath"
)

func Exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func PathSplit(p string) (dir, name, ext string) {
	dir, name = filepath.Split(p)
	ext = filepath.Ext(name)
	name = name[:len(name)-len(ext)]
	return
}

func IsRegular(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.Mode().IsRegular()
}
