package urlx

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const (
	HeaderRequestCookie = "Cookie" // Request Cookie
	// HeaderResponseCookie = "Set-Cookie" // Response Cookie
)

func CookieSet(cookieText string) Option {
	return func(c *Request) error { c.HeaderSet(HeaderRequestCookie, cookieText); return nil }
}

// CookieAddString 添加Cookie到请求
func CookieAddString(cookies ...string) Option {
	return Headers(func(headers http.Header) {
		for _, s := range cookies {
			if s != "" {
				if c := headers.Get(HeaderRequestCookie); c != "" {
					headers.Set(HeaderRequestCookie, c+"; "+s)
				} else {
					headers.Set(HeaderRequestCookie, s)
				}
			}
		}
	})
}

// CookieAdd 添加Cookie到请求
func CookieAdd(cookies ...*http.Cookie) Option {
	return Headers(func(headers http.Header) {
		for _, cookie := range cookies {
			if cookie != nil {
				s := fmt.Sprintf("%s=%s", sanitizeCookieName(cookie.Name), sanitizeCookieValue(cookie.Value))
				if c := headers.Get(HeaderRequestCookie); c != "" {
					headers.Set(HeaderRequestCookie, c+"; "+s)
				} else {
					headers.Set(HeaderRequestCookie, s)
				}
			}
		}
	})
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeCookieName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

func sanitizeCookieValue(v string) string {
	v = sanitizeOrWarn("Cookie.Value", validCookieValueByte, v)
	if len(v) == 0 {
		return v
	}
	if strings.ContainsAny(v, " ,") {
		return `"` + v + `"`
	}
	return v
}

func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}

func sanitizeOrWarn(fieldName string, valid func(byte) bool, v string) string {
	ok := true
	for i := 0; i < len(v); i++ {
		if valid(v[i]) {
			continue
		}
		slog.Debug(fmt.Sprintf("net/http: invalid byte %q in %s; dropping invalid bytes", v[i], fieldName), "pkg", "urlx")
		ok = false
		break
	}
	if ok {
		return v
	}
	buf := make([]byte, 0, len(v))
	for i := 0; i < len(v); i++ {
		if b := v[i]; valid(b) {
			buf = append(buf, b)
		}
	}
	return string(buf)
}
