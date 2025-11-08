package json

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type (
	Process = func(resp *http.Response) error // 响应处理器
	Body    = func() (contentType string, body io.Reader, err error)
)

// Decode 处理JSON响应
func Decode(out any) Process {
	return func(resp *http.Response) error {
		return NewDecoder(resp.Body).Decode(out)
	}
}

// Encode 提交JSON
func Encode(in any) Body {
	return func() (contentType string, body io.Reader, err error) {
		contentType = "application/json; charset=utf-8"
		switch o := in.(type) {
		case io.Reader:
			body = o
		case []byte:
			body = bytes.NewReader(o)
		case string:
			body = strings.NewReader(o)
		case bytes.Buffer:
			body = bytes.NewReader(o.Bytes())
		case *bytes.Buffer:
			body = bytes.NewReader(o.Bytes())
		default:
			var data []byte
			if data, err = Marshal(in); err == nil {
				body = bytes.NewReader(data)
			}
		}
		return
	}
}
