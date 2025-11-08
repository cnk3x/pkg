package urlx

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// 一些特定方法的定义
type (
	Option       = func(*Request) error                                   // 请求选项
	Body         = func() (contentType string, body io.Reader, err error) // 请求提交内容构造方法
	HeaderOption = func(headers http.Header)                              // 请求头处理
)

// Request 请求构造
type Request struct {
	ctx     context.Context        // Context
	options []func(*Request) error // options

	// request fields
	method    string         // 接口请求方法
	url       string         // 请求地址
	query     string         // 请求链接参数
	buildBody Body           // 请求内容
	headers   []HeaderOption // 请求头处理

	// response fields
	beforeMw []ProcessMw // 中间件

	// client fields
	tryTimes []time.Duration // 重试时间和时机
	client   *http.Client    // client
}

// New 以一些选项开始初始化请求器
func New(ctx context.Context, options ...Option) *Request {
	return (&Request{ctx: ctx}).With(options...)
}

/*请求公共设置*/

// With 增加选项
func (c *Request) With(options ...Option) *Request {
	c.options = append(c.options, options...)
	return c
}

// Method 设置请求方法
func (c *Request) Method(method string) *Request {
	c.method = method
	return c
}

// Url 设置请求链接
func (c *Request) Url(url string) *Request {
	c.url = url
	return c
}

// Query 设置请求Query参数
func (c *Request) Query(query string) *Request {
	c.query = query
	return c
}

// Body 设置请求提交内容
func (c *Request) Body(body Body) *Request {
	c.buildBody = body
	return c
}

func (c *Request) Form(formBody io.Reader) *Request {
	return c.Body(func() (contentType string, body io.Reader, err error) {
		contentType = "application/x-www-form-urlencoded; charset=utf-8"
		body = formBody
		return
	})
}

func (c *Request) FormValues(formBody url.Values) *Request {
	return c.Form(strings.NewReader(formBody.Encode()))
}

/* headers */

// HeaderWith 设置请求头
func (c *Request) HeaderWith(options ...HeaderOption) *Request {
	c.headers = append(c.headers, options...)
	return c
}

// HeaderSet 设置请求头
func HeaderSet(key string, values ...string) HeaderOption {
	return func(headers http.Header) {
		headers.Set(key, strings.Join(values, ","))
	}
}

// HeaderDel 删除请求头
func HeaderDel(keys ...string) HeaderOption {
	return func(headers http.Header) {
		for _, key := range keys {
			headers.Del(key)
		}
	}
}
