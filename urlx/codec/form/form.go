package form

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type (
	Process = func(resp *http.Response) error // 响应处理器
	Body    = func() (contentType string, body io.Reader, err error)
)

// Decode 处理Form表单响应
func Decode(out *url.Values) Process {
	return func(resp *http.Response) error {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		values, err := url.ParseQuery(string(respBody))
		if err != nil {
			return err
		}
		*out = values
		return nil
	}
}

// Encode 提交Form表单
func Encode(in any) Body {
	return func() (contentType string, body io.Reader, err error) {
		contentType = "application/x-www-form-urlencoded; charset=utf-8"
		switch o := in.(type) {
		case io.Reader:
			body = o
		case []byte:
			body = bytes.NewReader(o)
		case string:
			body = strings.NewReader(o)
		case url.Values:
			body = strings.NewReader(o.Encode())
		case *url.Values:
			body = strings.NewReader(o.Encode())
		case map[string]string:
			values := url.Values{}
			for k, v := range o {
				values.Set(k, v)
			}
			body = strings.NewReader(values.Encode())
		default:
			if r, ok := o.(io.Reader); ok {
				body = r
			} else {
				var values url.Values
				if values, err = query.Values(in); err == nil {
					body = strings.NewReader(values.Encode())
				}
			}
		}
		return
	}
}
