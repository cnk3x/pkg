package cmdo

import (
	"os/exec"
	"strings"
)

// LookPath 依次尝试查找给定名称的可执行文件路径，返回第一个找到的路径或空字符串
func LookPath(names ...string) string {
	for _, name := range names {
		if name = strings.TrimSpace(name); name == "" {
			continue
		}
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}
