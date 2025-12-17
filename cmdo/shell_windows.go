package cmdo

// GetShell 根据当前操作系统返回合适的 shell 命令及其参数
func GetShell() []string {
	s := LookPath("pwsh", "powershell")
	if s != "" {
		return []string{s, "-noprofile", "-nologo"}
	}
	return []string{"cmd"}
}
