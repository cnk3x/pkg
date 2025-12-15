package filex

import (
	"path/filepath"
	"testing"
)

func TestRel(t *testing.T) {
	base := `..\filex/`
	path := `path_test1.go`
	t.Logf("base: %s", base)
	t.Logf("path: %s", path)

	// if !PathIsDir(base) {
	// 	base = filepath.Dir(base)
	// 	t.Logf("trimed base: %s", base)
	// }

	base, _ = filepath.Abs(base)
	path, _ = filepath.Abs(path)

	base = filepath.ToSlash(base + "/")

	t.Logf("base abs: %s", base)
	t.Logf("path abs: %s", path)

	rel, e := filepath.Rel(base, path)
	if e != nil {
		t.Log(e)
	}
	t.Logf("rel: %s", rel)

	joined := filepath.Join(base, rel)
	t.Logf("joined: %s", joined)
}
