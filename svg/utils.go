package svg

import (
	"cmp"
	"os"
	"regexp"
)

func openWrite(destPath string, write func(file *os.File) error) (err error) {
	tempFile := destPath + ".temp"
	defer os.Remove(tempFile)

	var file *os.File
	if file, err = os.Create(tempFile); err != nil {
		return
	}

	if e1, e2 := write(file), file.Close(); e1 != nil || e2 != nil {
		err = cmp.Or(e1, e2)
		return
	}

	if err = os.Remove(destPath); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
	}

	err = os.Rename(tempFile, destPath)
	return
}

var cleanRe = regexp.MustCompile(`[\\/:*?"<>|.-]+`)
