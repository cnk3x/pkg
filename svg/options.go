package svg

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
)

type Options struct {
	ID     func(file string) string
	Pretty bool
}

type Option func(opts *Options)

func Pretty(pretty bool) Option                     { return func(opts *Options) { opts.Pretty = pretty } }
func GenerateID(id func(file string) string) Option { return func(opts *Options) { opts.ID = id } }

func NameFromBase(base string) Option {
	return GenerateID(func(file string) string {
		rel, err := filepath.Rel(base, file)
		if err != nil {
			slog.Error(fmt.Sprintf("获取 %s 相对路径失败: %v", file, err))
			return ""
		}

		return cleanRe.ReplaceAllString(strings.TrimSuffix(rel, filepath.Ext(rel)), "-")
	})
}

var cleanRe = regexp.MustCompile(`[\\/:*?"<>|.-]+`)
var spaceRe = regexp.MustCompile(`\s+`)
