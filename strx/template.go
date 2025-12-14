package strx

import (
	"io"
	"os"
	"strings"

	"github.com/valyala/fasttemplate"
)

type TagFind func(tag string) (value any, ok bool)

const sTag = "{{"
const eTag = "}}"

func EnvFind(tag string) (value any, ok bool) {
	if ok = len(tag) > 0 && tag[0] == '$'; ok && len(tag) > 1 {
		value = Atob(os.Getenv(tag[1:]))
	}
	return
}

func Replace(src string, tagFinds ...TagFind) string {
	return fasttemplate.ExecuteFuncString(src, sTag, eTag, func(w io.Writer, tag string) (int, error) {
		nTag := strings.TrimSpace(tag)

		if len(nTag) > 1 && nTag[0] == '$' {
			return w.Write(Atob(os.Getenv(nTag[1:])))
		}

		for _, tagFind := range tagFinds {
			if tagFind == nil {
				continue
			}
			if v, ok := tagFind(nTag); ok {
				if written, n, err := writeValue(w, nTag, v); written {
					return n, err
				}
			}
		}

		return writeBack(w, sTag, eTag, tag)
	})
}

func Replaces(src []string, tagFinds ...TagFind) (dst []string) {
	dst = make([]string, len(src))
	for i, s := range src {
		dst[i] = Replace(s, tagFinds...)
	}
	return dst
}

func ReplaceInline(src []string, tagFinds ...TagFind) {
	for i, s := range src {
		src[i] = Replace(s, tagFinds...)
	}
}

func ReplaceWithTagArg(src string, args ...any) string {
	return Replace(src, TagArg(args...))
}

func writeValue(w io.Writer, tag string, v any) (written bool, n int, err error) {
	written = true
	if v != nil {
		switch value := v.(type) {
		case []byte:
			n, err = w.Write(value)
		case string:
			n, err = w.Write([]byte(value))
		case fasttemplate.TagFunc:
			n, err = value(w, tag)
		case func() string:
			n, err = w.Write([]byte(value()))
		default:
			written = false
		}
	}
	return
}

func writeBack(w io.Writer, sTag, eTag, tag string) (int, error) {
	if _, err := w.Write(Atob(sTag)); err != nil {
		return 0, err
	}

	if _, err := w.Write(Atob(tag)); err != nil {
		return 0, err
	}

	if _, err := w.Write(Atob(eTag)); err != nil {
		return 0, err
	}

	return len(sTag) + len(tag) + len(eTag), nil
}

func TagArgs(args []string) TagFind {
	if len(args) < 2 {
		return nil
	}
	return func(tag string) (any, bool) {
		for i := 0; i < len(args)-1; i += 2 {
			if args[i] == tag {
				return args[i+1], true
			}
		}
		return nil, false
	}
}

func TagArg(args ...any) TagFind {
	if len(args) < 2 {
		return nil
	}
	return func(tag string) (any, bool) {
		for i := 0; i < len(args)-1; i += 2 {
			if args[i] == tag {
				return args[i+1], true
			}
		}
		return nil, false
	}
}

func TagMap(tagMap map[string]any) TagFind {
	if len(tagMap) == 0 {
		return nil
	}
	return func(tag string) (any, bool) {
		v, ok := tagMap[tag]
		return v, ok
	}
}

func TagUpper(finds ...TagFind) TagFind {
	return func(tag string) (value any, ok bool) {
		tag = strings.ToUpper(tag)
		for _, find := range finds {
			if value, ok = find(tag); ok {
				break
			}
		}
		return
	}
}
