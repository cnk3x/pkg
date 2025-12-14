package urlx

import (
	"context"
)

type Logger = func(ctx context.Context, msg string, args ...any)

var nop = func(ctx context.Context, msg string, args ...any) {}

func (c *Request) logger() Logger {
	if c.log != nil {
		return c.log
	}
	return nop
}

func (c *Request) Log(log Logger) *Request {
	c.log = log
	return c
}

func Log(log Logger) Option {
	return func(c *Request) error { c.log = log; return nil }
}
