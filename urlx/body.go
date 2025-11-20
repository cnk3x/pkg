package urlx

import (
	"context"
	"io"
	"net/url"
	"strings"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	CharsetUTF8       = "charset=utf-8"
)

type Body = func(ctx context.Context) (body io.Reader, contentType string, err error) // 请求提交内容构造方法

// Body 设置请求提交内容
func (c *Request) Body(body Body) *Request { c.body = body; return c }

func (c *Request) Form(formBody io.Reader) *Request {
	return c.Body(func(context.Context) (io.Reader, string, error) { return formBody, ContentTypeForm, nil })
}

func (c *Request) FormValues(formBody url.Values) *Request {
	return c.Form(strings.NewReader(formBody.Encode()))
}
