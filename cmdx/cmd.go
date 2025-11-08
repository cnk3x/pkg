package cmdx

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 适用于需要更多控制的长时间运行的场景，如需要自定义输出、错误处理等
type Cmd struct {
	Name string
	Options

	onPrepare func(c *Cmd) error
	onStart   func(c *exec.Cmd) error
	onStarted func(pid int)
	onExited  func(pid int, code int, err error) error

	options []Option

	status    string // 运行的状态
	cancelCur context.CancelFunc
	cancelTop context.CancelFunc

	done    <-chan struct{}
	started atomic.Bool
}

func New(options ...Option) *Cmd {
	x := &Cmd{options: options}

	x.onPrepare = func(*Cmd) error { return nil }
	x.onStart = func(*exec.Cmd) error { return nil }
	x.onStarted = func(int) {}
	x.onExited = func(_ int, _ int, err error) error { return err }

	return x.With(options...)
}

func (x *Cmd) With(options ...Option) *Cmd { x.options = append(x.options, options...); return x }

func (x *Cmd) Stop()                 { x.status = "stopping"; x.cancelTop() }
func (x *Cmd) Restart()              { x.status = "restarting"; x.cancelCur() }
func (x *Cmd) Done() <-chan struct{} { return x.done }
func (x *Cmd) Status() string        { return x.status }

func (x *Cmd) Start(ctx context.Context) error {
	if !x.started.CompareAndSwap(false, true) {
		slog.Debug("cmdx already started", "name", x.Name)
		return fmt.Errorf("cmdx already started")
	}

	top, cancel := context.WithCancel(ctx)
	x.cancelTop = cancel

	done := make(chan struct{})
	closeDone := sync.OnceFunc(func() { close(done) })
	x.done = done

	go func() {
		defer cancel()
		defer closeDone()
		defer x.started.Store(false)

		for i := 0; ; i++ {
			if err := x.directRun(top); err != nil {
				slog.Error("cmdx run error", "name", x.Name, "err", err)
			}

			delay, restartable := x.Options.Restart.ShouldRestart(i)
			if !restartable {
				slog.Debug("cmdx restart not restartable", "name", x.Name, "i", i)
				return
			}
			slog.Debug("cmdx restart delay", "name", x.Name, "delay", delay.String())

			select {
			case <-top.Done():
				return
			case <-time.After(delay):
			}
		}
	}()

	return nil
}

func (x *Cmd) directRun(ctx context.Context) (err error) {
	if x.status != "restarting" {
		x.statusUp("starting")
	}

	for _, option := range x.options {
		if err = option(x); err != nil {
			return fmt.Errorf("cmdx options apply error: %w", err)
		}
	}

	if x.onPrepare != nil {
		if err = x.onPrepare(x); err != nil {
			return fmt.Errorf("cmdx before start error: %w", err)
		}
	}

	if len(x.Command) > 0 && x.Name == "" {
		x.Name = strings.TrimSuffix(filepath.Base(x.Command[0]), filepath.Ext(x.Command[0]))
	}

	cur, cancel := context.WithCancel(ctx)
	defer cancel()
	x.cancelCur = cancel

	cmd := exec.CommandContext(cur, x.Command[0], x.Command[1:]...)
	cmd.Env = x.Env
	cmd.Dir = x.Dir

	var logClose func()
	if cmd.Stdout, cmd.Stderr, logClose, err = x.Log.Open(); err != nil {
		return fmt.Errorf("cmdx log open error: %w", err)
	}
	defer logClose()

	if err = x.onStart(cmd); err != nil {
		return fmt.Errorf("cmdx onStart fail: %w", err)
	}

	slog.Debug("cmdx starting", "name", x.Name, "cmd", cmd.String(), "dir", cmd.Dir)
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("cmdx start fail: %w", err)
	}

	if x.onStarted != nil {
		x.onStarted(cmd.Process.Pid)
	}

	slog.Debug("cmdx started", "name", x.Name, "pid", cmd.Process.Pid)
	x.statusUp("running")
	if err = cmd.Wait(); err != nil {
		err = fmt.Errorf("cmdx exited with code %d: %w", cmd.ProcessState.ExitCode(), err)
	}

	if x.onExited != nil {
		err = x.onExited(cmd.Process.Pid, cmd.ProcessState.ExitCode(), err)
	}
	slog.Debug("cmdx exited", "name", x.Name, "err", err)
	if x.status != "restarting" {
		x.statusUp("exited")
	}
	return
}

func (x *Cmd) statusUp(status string) {
	x.status = status
}
