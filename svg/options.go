package svg

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

type Options struct {
	ID func(file string) string
}

type Option func(opts *Options)

func GenerateID(id func(file string) string) Option {
	return func(opts *Options) {
		opts.ID = id
	}
}

func NameFromBase(base string) Option {
	return GenerateID(func(file string) string {
		rel, err := filepath.Rel(base, file)
		if err != nil {
			slog.Error(fmt.Sprintf("获取 %s 相对路径失败: %v", file, err))
			return ""
		}

		name := strings.TrimSuffix(rel, filepath.Ext(rel))
		return cleanRe.ReplaceAllString(name, "-")
	})
}
