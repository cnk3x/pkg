package webx

import (
	"cmp"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func Func(handle func() error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handle(); err != nil {
			Respond(w, r, E(err))
		}
	})
}

func FuncContext(handle func(context.Context) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handle(r.Context()); err != nil {
			Respond(w, r, E(err))
		}
	})
}

func JSON[T *I, I, O any](handle func(ctx context.Context, params T) (O, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error
		var in I
		var out O

		if data, err = io.ReadAll(r.Body); err != nil {
			StatusSet(r, http.StatusBadRequest)
			Respond(w, r, E(err, "READ_BODY"))
			return
		}

		if err = json.Unmarshal(data, in); err != nil {
			StatusSet(r, http.StatusBadRequest)
			Respond(w, r, E(err, "PARSER_BODY"))
			return
		}

		if out, err = handle(r.Context(), &in); err != nil {
			StatusSet(r, http.StatusInternalServerError)
			Respond(w, r, E(err, "PROCESS"))
			return
		}

		Respond(w, r, out)
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
