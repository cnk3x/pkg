package fsw

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

type WalkOptions struct {
	skipDir   bool                                       //skip dir
	skipFile  bool                                       //skip file
	defFilter bool                                       //skip `[.-~_]*`
	match     func(string) bool                          //
	debug     bool                                       //
	walk      func(fullPath string, d fs.DirEntry) error //
	base      string
}

// FileWalkRel base filepath.WalkDir but
//
//	relpath is the relative path base [base] with slash separator and add the slash prefix
//	walkFn not has [err] param in, if the err not nil will return direct
func FileWalk(dir string, options ...func(*WalkOptions)) (err error) {
	wo := WalkOptions{defFilter: true, debug: os.Getenv("CNK3X_FSW_DEBUG") == "1"}
	for _, apply := range options {
		apply(&wo)
	}

	if dir, err = filepath.Abs(dir); err != nil {
		return
	}

	if wo.base == "" {
		wo.base = dir
	} else if wo.base, err = filepath.Abs(wo.base); err != nil {
		return
	}

	return filepath.WalkDir(dir, func(fullPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if fullPath != wo.base {
			if d.IsDir() {
				if wo.skipDir {
					return nil
				}
			} else {
				if wo.skipFile {
					return nil
				}
			}

			if wo.match != nil && !wo.match(fullPath) {
				if wo.debug {
					slog.Debug("skip by filter", "full", fullPath)
				}
				return lo.Ternary(d.IsDir(), fs.SkipDir, nil)
			}
		}

		if wo.walk != nil {
			return wo.walk(fullPath, d)
		}

		return fs.SkipAll
	})
}

func OptionJoin(options ...func(*WalkOptions)) func(*WalkOptions) {
	return func(o *WalkOptions) {
		for _, apply := range options {
			apply(o)
		}
	}
}

func Walk[T func(string, fs.DirEntry) error | func(string, fs.DirEntry) | func(string) error | func(string)](walk T) func(*WalkOptions) {
	var walkFn func(string, fs.DirEntry) error
	if f, ok := any(walk).(func(string, fs.DirEntry) error); ok {
		walkFn = f
	} else if f, ok := any(walk).(func(string, fs.DirEntry)); ok {
		walkFn = func(s string, de fs.DirEntry) error { f(s, de); return nil }
	} else if f, ok := any(walk).(func(string) error); ok {
		walkFn = func(s string, de fs.DirEntry) error { return f(s) }
	} else if f, ok := any(walk).(func(string)); ok {
		walkFn = func(s string, de fs.DirEntry) error { f(s); return nil }
	} else {
		return func(*WalkOptions) {}
	}
	return func(o *WalkOptions) { o.walk = walkFn }
}

func WalkFilter(filter func(string) bool) func(*WalkOptions) {
	return func(o *WalkOptions) { o.match = filter }
}

func SkipDefaultFilter(o *WalkOptions) { o.defFilter = false }
func WithSkipFile(o *WalkOptions)      { o.skipFile = true }
func WithSkipDir(o *WalkOptions)       { o.skipDir = true }
