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
func NameFromBase(base string) Option               { return GenerateID(NameFunc(base)) }

func NameFunc(base string) func(string) string {
	var cleanRe = regexp.MustCompile(`[\s\\/:*?"<>|.-]+`)

	if base == "" {
		return func(file string) string {
			return cleanRe.ReplaceAllString(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)), "-")
		}
	}

	return func(file string) string {
		rel, err := filepath.Rel(base, file)
		if err != nil {
			slog.Warn(fmt.Sprintf("get relative path for %q fail: %v", file, err))
			return ""
		}
		return cleanRe.ReplaceAllString(strings.TrimSuffix(rel, filepath.Ext(rel)), "-")
	}
}
