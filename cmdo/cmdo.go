package cmdo

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/samber/lo"
)

type Option func(*exec.Cmd)

func Apply(c *exec.Cmd, opts ...Option) *exec.Cmd {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}

	if c.Env == nil {
		c.Env = os.Environ()
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func With(opts ...Option) Option {
	return func(c *exec.Cmd) {
		Apply(c, opts...)
	}
}

func PKill(c *exec.Cmd) {
	setPKill(c)
	c.Cancel = func() (err error) {
		if c.Process == nil || c.ProcessState != nil {
			slog.Debug("[cancel] 未运行")
			return
		}

		if err = terminate(c.Process); err == nil {
			slog.Debug("[cancel] 自定义调用，返回成功")
			return
		}

		slog.Debug("[cancel] 自定义调用，返回错误，改默认调用", "err", err)
		if err = c.Process.Kill(); err != nil {
			slog.Debug("[cancel] 默认调用，返回错误", "err", err)
			return
		}
		return
	}
}

func CancelWith(ctx context.Context) Option {
	return func(c *exec.Cmd) {
		cc := exec.CommandContext(ctx, c.Path, c.Args[1:]...)
		cc.Dir, cc.Env = c.Dir, c.Env
		cc.Stdin = c.Stdin
		cc.Stdout, cc.Stderr = c.Stdout, c.Stderr
		cc.SysProcAttr, cc.Cancel = c.SysProcAttr, c.Cancel
		*c = *cc
	}
}

func Dir(dir string) Option {
	return func(c *exec.Cmd) {
		c.Dir = dir
	}
}

func Env(env map[string]string) Option {
	return func(c *exec.Cmd) {
		if len(env) > 0 {
			c.Env = append(c.Env, lo.MapToSlice(env, func(key string, value string) string { return key + "=" + value })...)
		}
	}
}

func Std(c *exec.Cmd) {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
}
