package cmdo

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func setPKill(c *exec.Cmd) {
	c.SysProcAttr.CreationFlags |= syscall.CREATE_NEW_PROCESS_GROUP
	c.SysProcAttr.HideWindow = true
}

func terminate(proc *os.Process) error {
	return exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(int(proc.Pid))).Run()
}
