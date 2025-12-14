package respond

import (
	"encoding/xml"
	"io"
	"net/http"

	"github.com/cnk3x/gopkg/jsonx"
	"github.com/go-playground/form/v4"
)

func DecodeBody[I any](r *http.Request) (*I, error) {
	var in I
	if err := UnmarshalBody(r, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

func UnmarshalBody[I any](r *http.Request, in *I) (err error) {
	switch ctType := GetAcceptedContentType(r); ctType {
	case ContentTypeForm:
		if err = r.ParseForm(); err == nil {
			err = form.NewDecoder().Decode(in, r.Form)
		}
	case ContentTypeXML:
		defer io.Copy(io.Discard, r.Body) //nolint:errcheck
		err = xml.NewDecoder(r.Body).Decode(in)
	case ContentTypeJSON:
		err = jsonx.Decode(r.Body, in)
	default:
		err = form.NewDecoder().Decode(in, r.URL.Query())
	}
	return
}
