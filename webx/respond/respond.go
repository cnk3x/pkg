package respond

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/fs"
	"net/http"
)

type M map[string]any

type Stringer interface{ ~string | ~[]byte }

func Respond(w http.ResponseWriter, r *http.Request, data any) {
	switch GetAcceptedContentType(r) {
	case ContentTypeJSON:
		JSON(w, r, data)
	case ContentTypeXML:
		XML(w, r, data)
	default:
		JSON(w, r, data)
	}
}

func Blob(w http.ResponseWriter, r *http.Request, contentType string, data ...[]byte) {
	h := w.Header()
	h.Del("Content-Length")
	h.Set("X-Content-Type-Options", "nosniff")

	if contentType != "" {
		h.Set("Content-Type", contentType)
	}

	if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	for _, b := range data {
		w.Write(b) //nolint: errcheck
	}
}

func JSON(w http.ResponseWriter, r *http.Request, data any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Blob(w, r, "application/json; charset=utf-8", buf.Bytes())
}

func XML(w http.ResponseWriter, r *http.Request, data any) {
	var contentType = "application/xml; charset=utf-8"

	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(data); err != nil {
		buf.Reset()
		buf.WriteString(`<error>`)                //nolint: errcheck
		xml.EscapeText(&buf, []byte(err.Error())) //nolint: errcheck
		buf.WriteString(`</error>`)               //nolint: errcheck
		return
	}

	b := buf.Bytes()
	if !bytes.Contains(b[:min(100, buf.Len())], []byte(`<?xml`)) {
		Blob(w, r, contentType, []byte(xml.Header), b)
	} else {
		Blob(w, r, contentType, b)
	}
}

func PlainText[S Stringer](w http.ResponseWriter, r *http.Request, data S) {
	Blob(w, r, "text/plain; charset=utf-8", []byte(data))
}

func File(w http.ResponseWriter, r *http.Request, name string) {
	http.ServeFile(w, r, name)
}

func FileFS(w http.ResponseWriter, r *http.Request, fsys fs.FS, name string) {
	http.ServeFileFS(w, r, fsys, name)
}

func E(err error, code ...string) M {
	errText := cmp.Or(err.Error(), "error")
	codeText := cmp.Or(cmp.Or(code...), "error")
	return M{"err": errText, "code": codeText}
}

// StatusCtxKey is a context key to record a future HTTP response status code.
var StatusCtxKey = &contextKey{"Status"}

// Status sets a HTTP response status code hint into request context at any point
// during the request life-cycle. Before the Responder sends its response header
// it will check the StatusCtxKey
func Status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusCtxKey, status))
}

type contextKey struct{ name string }

func (k *contextKey) String() string { return "chi render context value " + k.name }
