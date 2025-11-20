package urlx

type Option = func(*Request) error

func Options(options ...Option) Option {
	return func(c *Request) error {
		c.With(options...)
		return nil
	}
}
