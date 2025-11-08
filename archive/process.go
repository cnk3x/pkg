package archive

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/cnk3x/gopkg/filex"
)

type Options struct {
	StripComponents int
	SkipEmptyDir    bool
	Report          func(index int, cur, total int64)
}

type Option func(option *Options)

func Extract(dir string, options ...Option) ProcessFunc {
	var extOpts Options
	for _, o := range options {
		o(&extOpts)
	}

	return func(item Item) error {
		path := item.Path()

		if extOpts.StripComponents > 0 {
			const sep = string(filepath.Separator)
			paths := Compact(strings.Split(strings.Trim(filepath.Clean(path), sep), sep))
			if len(paths) <= extOpts.StripComponents {
				return nil
			}
			path = filepath.Join(paths[extOpts.StripComponents:]...)
		}

		target := filepath.Join(dir, path)

		if item.IsDir() && !extOpts.SkipEmptyDir {
			err := os.MkdirAll(target, item.Mode())
			return err
		}

		it, err := item.Open()
		if err != nil {
			return err
		}

		current, total, index := int64(0), item.Size(), item.Index()
		progress := Iif(extOpts.Report != nil, func(n int64) { extOpts.Report(index, atomic.AddInt64(&current, n), total) }, nil)

		return filex.Process(target, filex.WriteFrom(it, progress), filex.CreateMode(item.Mode().Perm()))
	}
}

func Compact[T comparable](s []T) []T {
	var zero T
	return slices.DeleteFunc(s, func(it T) bool { return it == zero })
}

func May[T any](v T, _ error) T { return v }

func Iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}
