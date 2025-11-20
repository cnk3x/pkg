package urlx

import (
	"log/slog"
	"net/http"
	"time"
)

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

// Request 请求构造
type Request struct {
	options []func(*Request) error // options

	// request fields
	method  string         // 接口请求方法
	url     string         // 请求地址
	body    Body           // 请求内容
	headers []HeaderOption // 请求头处理

	// response fields
	uses []Process // 中间件

	// client fields
	tryTimes []time.Duration // 重试时间和时机
	client   *http.Client    // client

	//misc
	log      *slog.Logger
	logLevel slog.Level

	// 特别设置的头参数，优先级比通过 HeaderSet 方法设置的头参数高
	userAgent string
	referer   string
}

// New 以一些选项开始初始化请求器
func New(options ...Option) *Request { return (&Request{}).With(options...) }

/*请求公共设置*/

// With 增加选项
func (c *Request) With(options ...Option) *Request {
	c.options = append(c.options, options...)
	return c
}

// Method 设置请求方法
func (c *Request) Method(method string) *Request { c.method = method; return c }

// Url 设置请求链接
func (c *Request) Url(url string) *Request { c.url = url; return c }

// UserAgent 设置用户代理
func (c *Request) UserAgent(userAgent string) *Request { c.userAgent = userAgent; return c }

// Referer 设置请求来源
func (c *Request) Referer(referer string) *Request { c.referer = referer; return c }

// UserAgent 设置用户代理
func UserAgent(userAgent string) Option {
	return func(c *Request) error { c.UserAgent(userAgent); return nil }
}

// Referer 设置请求来源
func Referer(referer string) Option {
	return func(c *Request) error { c.UserAgent(referer); return nil }
}
