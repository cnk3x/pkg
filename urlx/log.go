package urlx

import "log/slog"

func Logger(log *slog.Logger, level slog.Level) Option {
	return func(c *Request) error {
		c.log = log
		c.logLevel = level
		return nil
	}
}
