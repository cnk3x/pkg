package webx

import (
	"net/http"

	"github.com/cnk3x/pkg/webx/respond"
)

type M = respond.M

func StatusSet(r *http.Request, status int)                    { respond.Status(r, status) }
func E(err error, code ...string) M                            { return respond.E(err, code...) }
func Respond(w http.ResponseWriter, r *http.Request, data any) { respond.Respond(w, r, data) }

func Blob(w http.ResponseWriter, r *http.Request, contentType string, data []byte) {
	respond.Blob(w, r, contentType, data)
}
