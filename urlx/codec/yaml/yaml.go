package yaml

import (
	"net/http"
)

type Process = func(resp *http.Response) error // 响应处理器

// Decode 处理yaml响应
func Decode(out any) Process {
	return func(resp *http.Response) error {
		return NewDecoder(resp.Body).Decode(out)
	}
}
