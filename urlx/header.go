package urlx

import (
	"net/http"
	"strings"
)

/* headers */

type HeaderOption = func(headers http.Header) // 请求头处理

func (c *Request) HeaderSet(key, value string) *Request {
	c.headers = append(c.headers, func(headers http.Header) {
		if key = strings.TrimSpace(key); key != "" {
			switch key[0] {
			case '+':
				headers.Add(key[1:], strings.TrimSpace(value))
			case '-':
				headers.Del(key[1:])
			default:
				headers.Set(key, strings.TrimSpace(value))
			}
		}
	})
	return c
}

func (c *Request) HeaderSets(lines ...string) *Request {
	c.headers = append(c.headers, headerSets(lines))
	return c
}

// HeaderSet 设置请求头
func HeaderSet(key, value string) Option {
	return func(c *Request) error { c.HeaderSet(key, value); return nil }
}

// HeaderSets 设置请求头
func HeaderSets(lines ...string) Option {
	return func(c *Request) error { c.HeaderSets(lines...); return nil }
}

func Headers(options ...HeaderOption) Option {
	return func(c *Request) error { ; c.headers = append(c.headers, options...); return nil }
}

func Accept(accept string) Option {
	return func(c *Request) error { c.HeaderSet("Accept", accept); return nil }
}
