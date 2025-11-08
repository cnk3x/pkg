package strx

import (
	"cmp"
	"io"
	"maps"
	"strings"

	"github.com/valyala/fasttemplate"
)

func ReplaceTemplate(src string, args ...string) string {
	return fasttemplate.ExecuteFuncString(src, "{", "}", func(w io.Writer, tag string) (int, error) {
		nTag := strings.TrimSpace(tag)
		for i := 0; i < len(args)-1; i += 2 {
			if args[i] == nTag {
				return io.WriteString(w, args[i+1])
			}
		}
		return io.WriteString(w, "{"+tag+"}")
	})
}

type Template struct {
	sTag string
	eTag string
	vars map[string]string
}

func NewTemplate(args ...string) *Template {
	return (&Template{}).Vars(args...)
}

func (t *Template) New(args ...string) *Template {
	return NewTemplate().Tag(t.sTag, t.eTag).VarMap(t.vars).Vars(args...)
}

func (t *Template) Tag(startTag, endTag string) *Template {
	t.sTag, t.eTag = startTag, endTag
	return t
}

func (t *Template) StartTag(startTag string) *Template {
	t.sTag = startTag
	return t
}

func (t *Template) EndTag(endTag string) *Template {
	t.eTag = endTag
	return t
}

func (t *Template) VarMap(vars map[string]string) *Template {
	if t.vars == nil {
		t.vars = make(map[string]string, len(vars))
	}
	maps.Copy(t.vars, vars)
	return t
}

func (t *Template) VarClear() *Template {
	clear(t.vars)
	return t
}

func (t *Template) Vars(vars ...string) *Template {
	if t.vars == nil {
		t.vars = make(map[string]string, len(vars)/2)
	}
	for i := 0; i < len(vars)-1; i += 2 {
		t.vars[vars[i]] = vars[i+1]
	}
	return t
}

func (t *Template) Replace(src string, args ...string) string {
	startTag, endTag := cmp.Or(t.sTag, "{"), cmp.Or(t.eTag, "}")

	return fasttemplate.ExecuteFuncString(src, startTag, endTag, func(w io.Writer, tag string) (int, error) {
		nTag := strings.TrimSpace(tag)

		if len(args) > 1 {
			for i := 0; i < len(args)-1; i += 2 {
				if args[i] == nTag {
					return io.WriteString(w, args[i+1])
				}
			}
		}

		if len(t.vars) > 0 {
			if v, ok := t.vars[nTag]; ok {
				return io.WriteString(w, v)
			}
		}

		return io.WriteString(w, startTag+tag+endTag)
	})
}
