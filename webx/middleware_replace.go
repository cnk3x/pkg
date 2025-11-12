package webx

import (
	"net/http"
	"net/url"
	"strings"
)

func Strip(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := strings.CutPrefix(r.URL.Path, prefix)
		if !ok {
			http.NotFound(w, r)
			return
		}

		var rp string
		if r.URL.RawPath != "" {
			if rp, ok = strings.CutPrefix(r.URL.RawPath, prefix); !ok {
				http.NotFound(w, r)
				return
			}
		}

		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		h.ServeHTTP(w, r2)
	})
	// return http.StripPrefix(prefix, handler)
}

// func StripMw(prefix string) Middleware {
// 	if prefix == "" {
// 		return Nop
// 	}
// 	return ReplaceMw(f2(strings.CutPrefix, prefix))
// }

// func ReplaceMw(mod ...func(string) (string, bool)) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			p1, p2 := r.URL.Path, r.URL.RawPath
// 			if !applyMod(mod, &p1, &p2) {
// 				http.NotFound(w, r)
// 				return
// 			}
// 			r1 := Clone(r)
// 			r1.URL.Path, r1.URL.RawPath = p1, p2
// 			next.ServeHTTP(w, r1)
// 		})
// 	}
// }

// func applyMod(mod []func(string) (string, bool), dst ...*string) (ok bool) {
// 	for _, s := range dst {
// 		for _, m := range mod {
// 			if *s, ok = m(*s); !ok {
// 				break
// 			}
// 		}
// 	}
// 	return
// }

// func f2[A, B, C, D any](f func(A, B) (C, D), b B) func(A) (C, D) {
// 	return func(a A) (C, D) { return f(a, b) }
// }
