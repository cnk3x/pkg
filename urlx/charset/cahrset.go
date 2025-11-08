package urlx

import (
	"io"
	"mime"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	HeaderContentType = "Content-Type"
	ParamCharset      = "charset"
)

type Process = func(resp *http.Response) error // 响应处理器
type ProcessMw = func(next Process) Process

// Charset 指定响应的编码，auto 或者空则通过 Content-Type 自动判断
func Charset(charset string) ProcessMw {
	getCharset := func(params map[string]string) string {
		if charset == "" || charset == "auto" {
			if len(params) > 0 {
				if charset = strings.TrimSpace(params[ParamCharset]); charset != "" {
					return charset
				}
			}
		}
		return strings.ToLower(charset)
	}

	return func(next Process) Process {
		return func(resp *http.Response) error {
			var body = io.Reader(resp.Body)
			mimeType, params, _ := mime.ParseMediaType(resp.Header.Get(HeaderContentType))
			if cs := getCharset(params); cs != "" && cs != "UTF-8" {
				if codec, err := htmlindex.Get(cs); err == nil && codec != unicode.UTF8 {
					body = transform.NewReader(body, codec.NewDecoder())
					resp.Header.Set(HeaderContentType, mimeType)
					resp.Body = io.NopCloser(body)
				}
			}
			return next(resp)
		}
	}
}

// AutoCharset 将响应解码成UTF-8
var AutoCharset = Charset("auto")
