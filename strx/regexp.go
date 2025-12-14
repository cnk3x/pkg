package strx

import (
	"log/slog"
	"regexp"
	"strings"
	"sync"
)

func createRegPool() func(expr string) (*regexp.Regexp, error) {
	type CompileResult struct {
		re  *regexp.Regexp
		err error
	}
	regPool := map[string]CompileResult{}
	mu := &sync.Mutex{}
	return func(expr string) (*regexp.Regexp, error) {
		mu.Lock()
		defer mu.Unlock()
		if v, find := regPool[expr]; find {
			return v.re, v.err
		}
		re, err := regexp.Compile(expr)
		regPool[expr] = CompileResult{re: re, err: err}
		return re, nil
	}
}

var RegexpE = createRegPool()

func Regexp(expr string) *regexp.Regexp { re, _ := RegexpE(expr); return re }

func ReReplace(src, pattern, repl string) string {
	if re, err := RegexpE(pattern); err == nil {
		return re.ReplaceAllString(src, repl)
	}
	return src
}

// Match: src === pattern || src contains pattern || pattern regexp match src
// 如果 pattern 以 !开头，则匹配反转。
func Match(src string, pattern string) bool {
	invert := pattern != "" && pattern[0] == '!'
	if invert {
		pattern = pattern[1:]
	}

	if pattern == "" || pattern == "*" {
		return !invert
	}

	if strings.Contains(src, pattern) {
		return !invert
	}

	re, err := RegexpE(pattern)
	if err != nil {
		slog.Warn("compile filter error", "filter", pattern, "err", err)
		return false
	}

	return re.MatchString(src) != invert
}

// MatchAll 匹配所有(and)，单个匹配规则同 Match
func MatchAll(src string, patterns []string) bool { return matches(src, patterns, false) }

// MatchAny 匹配任意一个(or)，单个匹配规则同 Match
func MatchAny(src string, patterns []string) bool { return matches(src, patterns, true) }

func matches(src string, patterns []string, isOr bool) bool {
	for _, pattern := range patterns {
		r := Match(src, pattern)
		if isOr {
			if r {
				return true
			}
		} else {
			if !r {
				return false
			}
		}
	}
	return !isOr
}
