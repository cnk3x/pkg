package cmdx

import "log/slog"

type Option func(*Program)

func (x Option) apply(p *Program) { x(p) }

func Prepare(onPrepare func(p *Config) error) Option {
	return func(p *Program) { p.onPrepare = onPrepare }
}

func ProcessInline(processInline func(s string, stdout bool)) Option {
	return func(p *Program) { p.cfg.Log.processInline = processInline }
}

func Log(logger *slog.Logger) Option {
	return func(p *Program) { p.log = logger }
}
