package browser

import (
	"strings"
	"time"

	"github.com/cnk3x/pkg/urlx"
)

const (
	HeaderAccept         = "Accept"
	HeaderAcceptLanguage = "Accept-Language"
	HeaderUserAgent      = "User-Agent"
	HeaderContentType    = "Content-Type"
	HeaderReferer        = "Referer"
	HeaderCacheControl   = "Cache-Control" // no-cache
	HeaderPragma         = "Pragma"        // no-cache
)

var (
	HeaderSet = urlx.HeaderSet
	Options   = urlx.Options
)

type (
	Option  = urlx.Option
	Request = urlx.Request
)

// AcceptLanguage 接受语言
func AcceptLanguage(acceptLanguages ...string) Option {
	return HeaderSet(HeaderAcceptLanguage, strings.Join(acceptLanguages, "; "))
}

// Accept 接受格式
func Accept(accepts ...string) Option {
	return HeaderSet(HeaderAccept, strings.Join(accepts, ", "))
}

// UserAgent 浏览器代理字符串
func UserAgent(userAgent string) Option {
	return HeaderSet(HeaderUserAgent, userAgent)
}

// Referer 引用地址
func Referer(referer string) Option {
	return HeaderSet(HeaderReferer, referer)
}

// Browser 浏览器
func Browser() *Request {
	ms := time.Millisecond
	return urlx.New().With(AcceptHTML, AcceptChinese).ErrTry(ms*300, ms*800, ms*1500)
}

// MacEdge Mac Edge 浏览器
func MacEdge() *Request {
	return Browser().With(MacEdgeAgent)
}

// WindowsEdge Windows Edge 浏览器
func WindowsEdge() *Request {
	return Browser().With(WindowsEdgeAgent)
}

// AndroidEdge Android Edge 浏览器
func AndroidEdge() *Request {
	return Browser().With(AndroidEdgeAgent)
}

var (
	MacChromeAgent  = UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36")
	MacFirefoxAgent = UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:65.0) Gecko/20100101 Firefox/65.0")
	MacSafariAgent  = UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.3 Safari/605.1.15")
	MacEdgeAgent    = UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36 Edg/96.0.1054.43")

	WindowsChromeAgent = UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	WindowsEdgeAgent   = UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36 Edg/96.0.1054.43")
	WindowsIEAgent     = UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko")

	IOSChromeAgent = UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) CriOS/31.0.1650.18 Mobile/11B554a Safari/8536.25")
	IOSSafariAgent = UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12F70 Safari/600.1.4")
	IOSEdgAgent    = UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/96.0.4664.55")

	AndroidChromeAgent = UserAgent("Mozilla/5.0 (Linux; Android 11; SM-G9910) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.59 Mobile Safari/537.36")
	AndroidWebkitAgent = UserAgent("Mozilla/5.0 (Linux; Android 11; SM-G9910) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30")
	AndroidEdgeAgent   = UserAgent("Mozilla/5.0 (Linux; Android 11; SM-G9910) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Mobile Safari/537.36 Edge/95.0.1020.55")
)

var (
	// AcceptChinese 接受中文
	AcceptChinese = AcceptLanguage("zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5")

	// AcceptHTML 接受网页浏览器格式
	AcceptHTML = Accept("text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

	// AcceptJSON 接受JSON格式
	AcceptJSON = Accept("application/json")

	// AcceptXML 接受XML格式
	AcceptXML = Accept("application/xml,text/xml")

	// AcceptAny 接受任意格式
	AcceptAny = Accept("*/*")

	// NoCache 无缓存
	NoCache = Options(HeaderSet(HeaderCacheControl, "no-cache"), HeaderSet(HeaderPragma, "no-cache"))
)
