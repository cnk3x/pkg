package compress

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zstd"
)

type (
	Process      = func(resp *http.Response) error // 响应处理器
	HeaderOption = func(headers http.Header)       // 请求头处理
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
)

// Decode 解压Body
func Decode(next Process) Process {
	return func(resp *http.Response) (err error) {
		body := resp.Body
		if cEncoding := resp.Header.Get(HeaderContentEncoding); cEncoding != "" {
			switch cEncoding {
			case "br":
				body = io.NopCloser(brotli.NewReader(body))
			case "deflate":
				body = flate.NewReader(body)
			case "gzip":
				body, err = gzip.NewReader(body)
			case "s2":
				body = io.NopCloser(s2.NewReader(body))
			case "snappy":
				body = io.NopCloser(snappy.NewReader(body))
			case "zstd", "zst":
				var b *zstd.Decoder
				if b, err = zstd.NewReader(body); err == nil {
					body = b.IOReadCloser()
				}
			}
			if err != nil {
				return
			}
		}

		if body != resp.Body {
			resp.Header.Del(HeaderContentEncoding)
			defer closes(body)
			resp.Body = body
		}

		return next(resp)
	}
}

// AcceptEncoding 接受编码
func AcceptEncoding(acceptEncodings ...string) HeaderOption {
	return func(headers http.Header) {
		headers.Set(HeaderAcceptEncoding, strings.Join(acceptEncodings, ","))
	}
}

var (
	// AcceptAllEncodings 接受所有的编码格式
	AcceptAllEncodings = AcceptEncoding("zstd", "br", "gzip", "deflate", "snappy", "s2")
	// DefaultEncodings 默认接受所有的编码格式
	DefaultEncodings = AcceptEncoding("gzip", "deflate", "br")
)

//closes 静默关闭 io.Closer
func closes(closer io.Closer, errPrintPrefix ...string) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			if len(errPrintPrefix) > 0 && errPrintPrefix[0] != "" {
				log.Printf("「%s」 %s", errPrintPrefix, err)
			}
		}
	}
}
