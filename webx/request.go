package webx

import (
	"net"
	"net/http"
	"strings"
)

const (
	HeaderUpgrade            = "Upgrade"
	HeaderXForwardedFor      = "X-Forwarded-For"
	HeaderXForwardedProto    = "X-Forwarded-Proto"
	HeaderXForwardedProtocol = "X-Forwarded-Protocol"
	HeaderXForwardedSsl      = "X-Forwarded-Ssl"
	HeaderXUrlScheme         = "X-Url-Scheme"
	HeaderXRealIP            = "X-Real-Ip"
)

func IsTLS(r *http.Request) bool { return r.TLS != nil }

func IsWebSocket(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get(HeaderUpgrade), "websocket")
}

func GetScheme(r *http.Request) string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if IsTLS(r) {
		return "https"
	}
	if scheme := r.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := r.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := r.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := r.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}

func GetRealIp(r *http.Request) string {
	// Fall back to legacy behavior
	if ip := r.Header.Get(HeaderXForwardedFor); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			xffIp := strings.TrimSpace(ip[:i])
			xffIp = strings.TrimPrefix(xffIp, "[")
			xffIp = strings.TrimSuffix(xffIp, "]")
			return xffIp
		}
		return ip
	}
	if ip := r.Header.Get(HeaderXRealIP); ip != "" {
		ip = strings.TrimPrefix(ip, "[")
		ip = strings.TrimSuffix(ip, "]")
		return ip
	}
	ra, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ra
}
