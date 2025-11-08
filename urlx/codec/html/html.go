package html

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Process = func(resp *http.Response) error

func Html(process func(s *goquery.Selection) error) Process {
	return func(resp *http.Response) error {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return fmt.Errorf("read as html: %w", err)
		}
		return process(doc.Selection)
	}
}

// Struct 将HTML解析为Struct
func Struct(out any, selection string, options StructOptions) Process {
	return Html(func(s *goquery.Selection) error {
		return BindStruct(s.Find(selection), out, options)
	})
}

// Map 将HTML解析为Map, out 必须为 *any
func Map(out *[]any, field MapField, params map[string]string) Process {
	return Html(func(s *goquery.Selection) (err error) {
		*out, err = BindMapField(s, field, params)
		return
	})
}
