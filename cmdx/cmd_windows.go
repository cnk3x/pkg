package cmdx

import (
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func graceStop(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags |= syscall.CREATE_NEW_PROCESS_GROUP
	cmd.SysProcAttr.HideWindow = true
	cmd.WaitDelay = time.Second * 10
	cmd.Cancel = func() error {
		if cmd.Process == nil || cmd.Process.Pid <= 0 {
			return nil
		}

		if err := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run(); err != nil {
			return err
		}

		// 使用标准的 Kill 方法终止主进程
		return cmd.Process.Kill()
	}

	return nil
}
