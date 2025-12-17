//go:build !windows

package cmdo

import "os"

// GetShell 根据当前操作系统返回合适的 shell 命令及其参数
func GetShell() []string {
	if s := LookPath(os.Getenv("SHELL"), "bash", "sh", "ash"); s != "" {
		return []string{s}
	}
	return nil
}
