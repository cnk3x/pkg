package urlx

import "net/http"

// StatusRead 从响应读取状态码并做预处理
func StatusRead(processStatus func(status int) error) Process {
	return func(resp *http.Response) error {
		return processStatus(resp.StatusCode)
	}
}

func HeaderRead(read func(headers http.Header) error) Process {
	return func(resp *http.Response) error {
		return read(resp.Header)
	}
}

// CookieRead 从响应读取Cookie
func CookieRead(read func(cookies []*http.Cookie) error) Process {
	return func(resp *http.Response) error {
		return read(resp.Cookies())
	}
}
