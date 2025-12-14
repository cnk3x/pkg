package rex

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fasttemplate"
)

func Template(re *regexp.Regexp, src string, template string) string {
	if re == nil {
		return ""
	}

	if match := re.FindStringSubmatch(src); match != nil {
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			result[strconv.Itoa(i)] = match[i]
			if i > 0 && name != "" {
				result[name] = match[i]
			}
		}

		return fasttemplate.ExecuteFuncString(template, "${", "}", func(w io.Writer, tag string) (int, error) {
			if val, ok := result[strings.TrimSpace(tag)]; ok {
				return w.Write([]byte(val))
			}
			return w.Write([]byte(tag))
		})
	}

	return ""
}
