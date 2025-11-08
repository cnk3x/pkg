package urlx

import (
	"context"
	"net/http"
	"strings"
	"time"
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

// AcceptLanguage 接受语言
func AcceptLanguage(acceptLanguages ...string) HeaderOption {
	return HeaderSet(HeaderAcceptLanguage, strings.Join(acceptLanguages, "; "))
}

// Accept 接受格式
func Accept(accepts ...string) HeaderOption {
	return HeaderSet(HeaderAccept, strings.Join(accepts, ", "))
}

// UserAgent 浏览器代理字符串
func UserAgent(userAgent string) HeaderOption {
	return HeaderSet(HeaderUserAgent, userAgent)
}

// Referer 引用地址
func Referer(referer string) HeaderOption {
	return HeaderSet(HeaderReferer, referer)
}

// Browser 浏览器
func Browser(ctx context.Context) *Request {
	ms := time.Millisecond
	return New(ctx).HeaderWith(AcceptHTML, AcceptChinese).TryAt(ms*300, ms*800, ms*1500)
}

// MacEdge Mac Edge 浏览器
func MacEdge(ctx context.Context) *Request {
	return Browser(ctx).HeaderWith(MacEdgeAgent)
}

// WindowsEdge Windows Edge 浏览器
func WindowsEdge(ctx context.Context) *Request {
	return Browser(ctx).HeaderWith(WindowsEdgeAgent)
}

// AndroidEdge Android Edge 浏览器
func AndroidEdge(ctx context.Context) *Request {
	return Browser(ctx).HeaderWith(AndroidEdgeAgent)
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
	NoCache = HeaderOption(func(headers http.Header) {
		headers.Set(HeaderCacheControl, "no-cache")
		headers.Set(HeaderPragma, "no-cache")
	})
)

// type BrowserBuilder struct {
// 	Mozilla     string
// 	Platform    []string
// 	Devices     string
// 	AppleWebKit string
// 	Apps        []BrowserApp
// }
//
// type BrowserApp struct {
// 	Name    string
// 	Version string
// }
//
// func (b BrowserBuilder) String() string {
// 	var ua bytes.Buffer
// 	if b.Mozilla == "" {
// 		b.Mozilla = "5.0"
// 	}
// 	ua.WriteString("Mozilla/")
// 	ua.WriteString(b.Mozilla)
//
// 	if len(b.Platform) > 0 || b.Devices != "" {
// 		ua.WriteByte('(')
// 		for i, p := range b.Platform {
// 			if i > 0 {
// 				ua.WriteString("; ")
// 			}
// 			ua.WriteString(p)
// 		}
// 		if b.Devices != "" {
// 			if len(b.Platform) > 0 {
// 				ua.WriteString("; ")
// 			}
// 			ua.WriteString(b.Devices)
// 		}
// 		ua.WriteByte(')')
// 	}
//
// 	if b.AppleWebKit != "" {
// 		ua.WriteString("AppleWebKit/")
// 		ua.WriteString(b.AppleWebKit)
// 		ua.WriteString("(KHTML, like Gecko)")
// 	}
//
// 	for _, app := range b.Apps {
// 		ua.WriteByte(' ')
// 		ua.WriteString(app.Name)
// 		ua.WriteByte('/')
// 		ua.WriteString(app.Version)
// 	}
//
// 	return ua.String()
// }
