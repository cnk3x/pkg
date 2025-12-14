package rex

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

func CompileMatch(patterns ...string) func(string) bool {
	if len(patterns) == 0 {
		return nil
	}

	type Matcher struct {
		Match  func(string) bool
		Invert bool
	}

	matches := lo.FilterMap(patterns, func(pattern string, _ int) (*Matcher, bool) {
		invert := strings.HasPrefix(pattern, "!")
		if invert {
			pattern = pattern[1:]
		}
		if re, e := regexp.Compile(pattern); e == nil {
			return &Matcher{Match: re.MatchString, Invert: invert}, true
		}
		return nil, false
	})

	return func(fullPath string) (match bool) {
		if fullPath != "" {
			fullPath = filepath.ToSlash(filepath.Clean(fullPath))
		}

		return lo.Reduce(matches, func(acc bool, matcher *Matcher, i int) bool {
			if !matcher.Invert {
				return acc || matcher.Match(fullPath)
			}
			if i > 0 && !acc {
				return false
			}
			return !matcher.Match(fullPath)
		}, false)
	}
}
