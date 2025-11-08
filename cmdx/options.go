package cmdx

import (
	"os/exec"

	"github.com/cnk3x/gopkg/jsonx"
)

type Options struct {
	Command jsonx.Strings  `json:"command,omitempty"`
	Env     jsonx.Strings  `json:"env,omitempty"`
	Dir     string         `json:"dir,omitempty"`
	Log     LogOptions     `json:"log,omitempty"`
	Restart RestartOptions `json:"restart,omitempty"`
}

func Args(args ...string) []string { return args }

type Option func(*Cmd) error

func (x Option) apply(c *Cmd) error { return x(c) }

func Name(name string) Option  { return func(c *Cmd) error { c.Name = name; return nil } }
func Env(env ...string) Option { return func(c *Cmd) error { c.Env = env; return nil } }
func Dir(dir string) Option    { return func(c *Cmd) error { c.Dir = dir; return nil } }
func GraceStop(c *Cmd) error   { return OnStart(graceStop).apply(c) }

func OnStart(onStart func(*exec.Cmd) error) Option {
	return func(c *Cmd) error { c.onStart = onStart; return nil }
}

func Prepare(onPrepare func(c *Cmd) error) Option {
	return func(c *Cmd) error { c.onPrepare = onPrepare; return nil }
}

func OnStarted(onStarted func(pid int)) Option {
	return func(c *Cmd) error { c.onStarted = onStarted; return nil }
}

func OnExited(onExited func(pid int, code int, err error) error) Option {
	return func(c *Cmd) error { c.onExited = onExited; return nil }
}

func Log(log string) Option {
	return func(c *Cmd) error { c.Log.Out, c.Log.Err = log, log; return nil }
}
