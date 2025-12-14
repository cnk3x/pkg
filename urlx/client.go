package urlx

import (
	"cmp"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

/* 设置客户端 */

type ClientOption func(cli *http.Client) error

func transportDefault() *http.Transport {
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// Client 使用自定义的HTTP客户端
func Client(client ClientOption) Option {
	return func(r *Request) error {
		r.clientOptions = append(r.clientOptions, client)
		return nil
	}
}

// Client 使用自定义的HTTP客户端
func (c *Request) Client(clientOptions ...ClientOption) *Request {
	c.clientOptions = append(c.clientOptions, clientOptions...)
	return c
}

// RoundTrip 自定义 RoundTripper
func (c *Request) RoundTrip(transport http.RoundTripper) *Request {
	return c.Client(func(cli *http.Client) error {
		cli.Transport = transport
		return nil
	})
}

// Proxy 设置代理
func Proxy(proxyUrlString string) Option {
	return Client(func(cli *http.Client) error {
		if cli.Transport == nil {
			cli.Transport = transportDefault()
			return nil
		}
		if tr, ok := cli.Transport.(*http.Transport); ok {
			fixedURL, err := url.Parse(proxyUrlString)
			if err != nil {
				return err
			}
			tr.Proxy = http.ProxyURL(fixedURL)
			return nil
		}

		if tr, ok := cli.Transport.(interface{ SetProxy(string) error }); ok {
			return tr.SetProxy(proxyUrlString)
		}

		return fmt.Errorf("不支持的 Transport 类型: %T", cli.Transport)
	})
}

// CookieEnabled 开关 Cookie
func CookieEnabled(enabled ...bool) Option {
	if cmp.Or(enabled...) {
		return Jar(errDel(cookiejar.New(nil)))
	}
	return Jar(nil)
}

// Jar 设置Cookie容器
func Jar(jar http.CookieJar) Option {
	return Client(func(cli *http.Client) error {
		cli.Jar = jar
		return nil
	})
}

// ErrTry 失败重试，等待休眠时间
func (c *Request) ErrTry(times ...time.Duration) *Request {
	c.tryTimes = times
	return c
}

// Idempotent 幂等重试
func Idempotent(base time.Duration, maxTimes int) Option {
	tryAts := make([]time.Duration, maxTimes)
	for i := range maxTimes {
		tryAts[i] = base * (1 << i)
	}
	return func(r *Request) error {
		if len(tryAts) > 0 {
			r.ErrTry(tryAts...)
		}
		return nil
	}
}

func ErrTry(tryAt ...time.Duration) Option {
	return func(c *Request) error { c.ErrTry(tryAt...); return nil }
}
