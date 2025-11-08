package strx

import (
	"path/filepath"
	"strings"
)

func Match(pattern, name string) bool {
	if strings.ContainsAny(pattern, "*?[]") {
		matched, _ := filepath.Match(pattern, name)
		return matched
	}
	return strings.EqualFold(name, pattern)
}
