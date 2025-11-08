package xml

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type (
	Process = func(resp *http.Response) error // 响应处理器
	Body    = func() (contentType string, body io.Reader, err error)
)

// Decode 处理xml响应
func Decode(out any) Process {
	return func(resp *http.Response) error {
		return xml.NewDecoder(resp.Body).Decode(out)
	}
}

func Encode(in any) Body {
	return func() (contentType string, body io.Reader, err error) {
		contentType = "application/xml; charset=utf-8"
		if in == nil {
			return
		}
		var ok bool
		if body, ok = in.(io.Reader); ok {
			return
		}
		switch o := in.(type) {
		case io.Reader:
			body = o
		case []byte:
			body = bytes.NewReader(o)
		case string:
			body = strings.NewReader(o)
		case bytes.Buffer:
			body = &o
		case *bytes.Buffer:
			body = o
		default:
			var data []byte
			if data, err = xml.Marshal(in); err == nil {
				body = bytes.NewReader(data)
			}
		}
		return
	}
}
