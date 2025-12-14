package webx

import (
	"cmp"
	"context"
	"fmt"
	"net/http"

	"github.com/cnk3x/gopkg/webx/respond"
)

func HandleSimple(handle func()) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle()
	})
}

func HandleForm(handle func(value string), field string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(r.FormValue(field))
	})
}

// func HandleFunc1(handle func() error) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if err := handle(); err != nil {
// 			Respond(w, r, E(err))
// 		}
// 	})
// }

// func HandleContext(handle func(context.Context)) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		handle(r.Context())
// 	})
// }

// func HandleContextE(handle func(context.Context) error) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if err := handle(r.Context()); err != nil {
// 			Respond(w, r, E(err))
// 		}
// 	})
// }

func HandleFunc[T func() | func() error | func(context.Context) | func(context.Context) error](handle T) http.Handler {
	var h func(context.Context) error
	switch t := any(handle).(type) {
	case func():
		h = func(context.Context) error { t(); return nil }
	case func() error:
		h = func(context.Context) error { return t() }
	case func(context.Context):
		h = func(ctx context.Context) error { t(ctx); return nil }
	case func(context.Context) error:
		h = func(ctx context.Context) error { return t(ctx) }
	default:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("%T", handle), http.StatusNotImplemented)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(r.Context()); err != nil {
			Respond(w, r, E(err))
		}
	})
}

func Handle[I, O any](handle func(context.Context, *I) (O, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := respond.DecodeBody[I](r)
		if err != nil {
			respond.Status(r, http.StatusBadRequest)
			respond.Respond(w, r, E(err, "RESOLVE_INPUT"))
			return
		}

		out, err := handle(r.Context(), in)
		if err != nil {
			respond.Status(r, http.StatusInternalServerError)
			respond.Respond(w, r, E(err, "PROCESS"))
			return
		}

		respond.Respond(w, r, out)
	})
}

func HandleSSE[I, O any](handle func(context.Context, *I) (<-chan O, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in I
		if err := respond.UnmarshalBody(r, &in); err != nil {
			respond.Status(r, http.StatusBadRequest)
			respond.Respond(w, r, E(err, "RESOLVE_INPUT"))
			return
		}

		out, err := handle(r.Context(), &in)
		if err != nil {
			respond.Status(r, http.StatusInternalServerError)
			respond.Respond(w, r, E(err, "PROCESS"))
			return
		}

		respond.ServerEvent(w, r, respond.ServerEventSource[O]{Heartbeat: 30, Data: out})
		respond.Respond(w, r, out)
	})
}

func Redirect(u string, permanent ...bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cmp.Or(permanent...) {
			http.Redirect(w, r, u, http.StatusPermanentRedirect)
		} else {
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		}
	})
}
