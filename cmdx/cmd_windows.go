package cmdx

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func setupCmd(c *exec.Cmd) {
	c.SysProcAttr.CreationFlags |= syscall.CREATE_NEW_PROCESS_GROUP
	c.SysProcAttr.HideWindow = true
}

func cancelProc(proc *os.Process) error {
	return exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(proc.Pid)).Run()
}
