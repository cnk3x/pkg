package cmdo

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/samber/lo"
)

// Option 定义了一个函数类型，用于配置exec.Cmd实例
type Option func(*exec.Cmd)

// Apply 将一系列选项应用到exec.Cmd实例上，并确保SysProcAttr和Env已正确初始化
//
// 参数:
//   - c: 要应用选项的exec.Cmd实例
//   - opts: 要应用的一系列选项
//
// 返回值:
//   - *exec.Cmd: 配置完成的exec.Cmd实例
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

// With 允许将多个选项组合成一个单一选项
//
// 参数:
//   - opts: 要组合的一系列选项
//
// 返回值:
//   - Option: 组合后的选项函数
func With(opts ...Option) Option {
	return func(c *exec.Cmd) {
		Apply(c, opts...)
	}
}

// PKill 设置命令的进程终止行为，提供自定义的取消逻辑
// 优先使用平台特定的终止方式，失败时回退到标准的Kill方法
//
// 参数:
//   - c: 要设置终止行为的exec.Cmd实例
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

// CancelWith 使用提供的context创建一个可取消的命令
//
// 参数:
//   - ctx: 用于控制命令生命周期的context
//
// 返回值:
//   - Option: 可将context应用到命令的选项函数
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

// Dir 创建一个设置工作目录的选项
//
// 参数:
//   - dir: 要设置为工作目录的路径
//
// 返回值:
//   - Option: 可将目录设置应用到命令的选项函数
func Dir(dir string) Option {
	return func(c *exec.Cmd) {
		c.Dir = dir
	}
}

// Env 创建一个设置环境变量的选项
//
// 参数:
//   - env: 包含要设置的环境变量键值对的映射
//
// 返回值:
//   - Option: 可将环境变量设置应用到命令的选项函数
func Env(env map[string]string) Option {
	return func(c *exec.Cmd) {
		if len(env) > 0 {
			c.Env = append(c.Env, lo.MapToSlice(env, func(key string, value string) string { return key + "=" + value })...)
		}
	}
}

// Std 将命令的标准输入、输出和错误流连接到操作系统对应的标准流
//
// 参数:
//   - c: 要设置标准流的exec.Cmd实例
func Std(c *exec.Cmd) {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
}
