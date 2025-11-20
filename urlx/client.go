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

func clientDefault() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// Try 失败重试，等待休眠时间
func (c *Request) Try(times ...time.Duration) *Request { c.tryTimes = times; return c }

// Client 使用的客户端定义
func (c *Request) Client(client *http.Client) *Request { c.client = client; return c }

// RoundTrip 自定义 RoundTripper
func (c *Request) RoundTrip(transport http.RoundTripper) *Request {
	c.client.Transport = transport
	return c
}

// Proxy 设置代理
func Proxy(proxyUrlString string) Option {
	return func(r *Request) error {
		if tr, ok := r.client.Transport.(*http.Transport); ok {
			fixedURL, err := url.Parse(proxyUrlString)
			if err != nil {
				return err
			}
			tr.Proxy = http.ProxyURL(fixedURL)
			return nil
		}

		if tr, ok := r.client.Transport.(interface{ SetProxy(string) error }); ok {
			return tr.SetProxy(proxyUrlString)
		}

		return fmt.Errorf("不支持的 Transport 类型: %T", r.client.Transport)
	}
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
	return func(c *Request) error {
		c.client.Jar = jar
		return nil
	}
}

// UseClient 使用自定义的HTTP客户端
func Client(client *http.Client) Option {
	return func(r *Request) error {
		r.client = client
		return nil
	}
}

// Idempotent 幂等重试
func Idempotent(base time.Duration, maxTimes int) Option {
	trys := make([]time.Duration, maxTimes)
	for i := range maxTimes {
		trys[i] = base * (1 << i)
	}
	return func(r *Request) error {
		if len(trys) > 0 {
			r.Try(trys...)
		}
		return nil
	}
}
