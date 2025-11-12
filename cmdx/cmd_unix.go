//go:build !windows

package cmdx

import (
	"os"
	"os/exec"
	"syscall"
)

func setupCmd(c *exec.Cmd) {
	c.SysProcAttr.Setpgid = true
}

func cancelProc(proc *os.Process) error {
	return syscall.Kill(-proc.Pid, syscall.SIGINT)
}
