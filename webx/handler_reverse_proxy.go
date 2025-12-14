package webx

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/cnk3x/pkg/x"
)

func SingleHostReverseProxy(urlString string) http.Handler {
	return httputil.NewSingleHostReverseProxy(x.May(url.Parse(urlString)))
}

func ReverseProxy(getUri func() *url.URL) http.Handler {
	return &httputil.ReverseProxy{Director: director(getUri)}
}

func director(getUri func() *url.URL) func(req *http.Request) {
	return func(req *http.Request) {
		u := getUri()
		targetQuery := u.RawQuery
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(u, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
}

func joinURLPath(a, b *url.URL) (path, rawPath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	aPath := a.EscapedPath()
	bPath := b.EscapedPath()

	aSlash := strings.HasSuffix(aPath, "/")
	bSlash := strings.HasPrefix(bPath, "/")

	switch {
	case aSlash && bSlash:
		return a.Path + b.Path[1:], aPath + bPath[1:]
	case !aSlash && !bSlash:
		return a.Path + "/" + b.Path, aPath + "/" + bPath
	}
	return a.Path + b.Path, aPath + bPath
}

func singleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	}
	return a + b
}
