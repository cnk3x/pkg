package archive

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/cnk3x/pkg/filex"
)

type options struct {
	stripComponents int
	skipEmptyDir    bool
	progress        func(index int, name string, cur, total int64)
	filters         []string
}

const pathSeparator = string(filepath.Separator)

func Extract(dir string, extractOptions ...Option) ProcessFunc {
	var eop options
	for _, o := range extractOptions {
		o(&eop)
	}

	return func(ctx context.Context, item Item) error {
		fpath := filepath.Clean(item.Path())

		if eop.stripComponents > 0 {
			paths := strings.Split(fpath, pathSeparator)
			if len(paths) <= eop.stripComponents {
				return nil
			}
			fpath = filepath.Join(paths[eop.stripComponents:]...)
		}

		if fpath == "" {
			return nil
		}

		if len(eop.filters) > 0 {
			for _, f := range eop.filters {
				re, e := regexp.Compile(f)
				if e != nil {
					return e
				}
				if !re.MatchString(fpath) {
					return nil
				}
			}
		}

		target := filepath.Join(dir, fpath)

		if item.IsDir() && !eop.skipEmptyDir {
			err := os.MkdirAll(target, item.Mode())
			return err
		}

		it, err := item.Open()
		if err != nil {
			return err
		}

		var p filex.ProgressFunc
		if eop.progress != nil {
			current, total, index := int64(0), item.Size(), item.Index()
			p = func(n int64) { eop.progress(index, fpath, atomic.AddInt64(&current, n), total) }
		}
		return filex.Process(target, filex.WriteFrom(ctx, it, p), filex.CreateMode(item.Mode().Perm()))
	}
}

type Option func(option *options)

func Filter(filters ...string) Option { return func(option *options) { option.filters = filters } }

func StripComponents(deep int) Option {
	return func(option *options) { option.stripComponents = deep }
}

func SkipEmptyDir(skip bool) Option {
	return func(option *options) { option.skipEmptyDir = skip }
}

func Progress(progress func(index int, name string, cur, total int64)) Option {
	return func(option *options) { option.progress = progress }
}
