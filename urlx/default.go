package urlx

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

const Version = "0.0.1"

const ms = time.Millisecond

func DefaultAgent(r *Request) error {
	type x struct{}
	r.userAgent = fmt.Sprintf("urlx/%s (%s) golang/%s(%s %s)", Version, reflect.TypeOf(x{}).PkgPath(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
	return nil
}

var (
	WindowsEdgeAgent   = UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36 Edg/96.0.1054.43")
	MacSafariAgent     = UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.3 Safari/605.1.15")
	IOSSafariAgent     = UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12F70 Safari/600.1.4")
	AndroidWebkitAgent = UserAgent("Mozilla/5.0 (Linux; Android 11; SM-G9910) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30")
)

// Default 默认的请求器
func Default() *Request { return New().ErrTry(300*ms, 800*ms, 1500*ms).With(DefaultAgent) }
func Windows() *Request { return New().ErrTry(300*ms, 800*ms, 1500*ms).With(WindowsEdgeAgent) }
func Mac() *Request     { return New().ErrTry(300*ms, 800*ms, 1500*ms).With(MacSafariAgent) }
func IOS() *Request     { return New().ErrTry(300*ms, 800*ms, 1500*ms).With(IOSSafariAgent) }
func Android() *Request { return New().ErrTry(300*ms, 800*ms, 1500*ms).With(AndroidWebkitAgent) }
