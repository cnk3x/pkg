package strx

import (
	"path/filepath"
	"strings"
)

func GlobMatch(pattern, name string) bool {
	if strings.EqualFold(name, pattern) {
		return true
	}

	if strings.ContainsAny(pattern, "*?[]") {
		matched, _ := filepath.Match(pattern, name)
		return matched
	}

	return false
}
