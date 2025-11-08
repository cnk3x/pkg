//go:build !windows

package cmdx

import (
	"os/exec"
	"syscall"
	"time"
)

func graceStop(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	cmd.WaitDelay = time.Second * 10
	cmd.Cancel = func() error {
		if cmd.Process == nil || cmd.Process.Pid <= 0 {
			return nil
		}

		// 先尝试发送 SIGINT 信号，允许进程优雅退出
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		if err != nil {
			return err
		}

		// 等待一段时间，让进程有时间优雅退出
		time.Sleep(time.Second * 3)

		// 再次检查进程是否仍在运行
		// 如果进程仍在运行，则发送 SIGKILL 强制终止
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	return nil
}
