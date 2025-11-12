package webx

import "net/http"

type Middleware func(next http.Handler) http.Handler

var Nop Middleware = func(next http.Handler) http.Handler { return next }
