package html

import (
	"io"

	"github.com/valyala/fasttemplate"
)

func noExTag(params map[string]string) func(w io.Writer, tag string) (n int, err error) {
	return func(w io.Writer, tag string) (n int, err error) {
		if params != nil {
			if v := params[tag]; v != "" {
				n, err = w.Write([]byte(v))
			}
		}
		return
	}
}

// ReplaceTemplate 模板替换
func ReplaceTemplate(template string, params map[string]string) (s string) {
	s, _ = fasttemplate.ExecuteFuncStringWithErr(template, "{", "}", noExTag(params))
	return
}
