//go:build !windows

package cmdo

import (
	"os"
	"os/exec"
	"syscall"
)

func setPKill(c *exec.Cmd) {
	c.SysProcAttr.Setpgid = true
}

func terminate(proc *os.Process) error {
	return syscall.Kill(int(-proc.Pid), syscall.SIGTERM)
}
