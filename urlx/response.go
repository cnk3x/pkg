package urlx

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cnk3x/gopkg/errx"
	"github.com/cnk3x/gopkg/logx"
)

/*处理响应*/

type (
	Process   = func(resp *http.Response) error // 响应处理器
	ProcessMw = func(next Process) Process      // 响应预处理器
)

// Use 在处理之前的预处理
func (c *Request) Use(mws ...ProcessMw) *Request {
	c.uses = append(c.uses, mws...)
	return c
}

// Status .
func Status(processStatus func(status int) Process) ProcessMw {
	return func(next Process) Process {
		return func(resp *http.Response) error {
			if process := processStatus(resp.StatusCode); process != nil {
				return process(resp)
			}
			return next(resp)
		}
	}
}

// Process 处理响应
func (c *Request) Process(ctx context.Context, process ...Process) error {
	log := logx.With("请求")
	if c.client == nil {
		c.client = &http.Client{}
	}

	for _, apply := range c.options {
		if err := apply(c); err != nil {
			return err
		}
	}

	if c.method == "" {
		c.method = http.MethodGet
	}

	requestUrl := c.url
	if c.query != "" {
		if strings.Contains(requestUrl, "?") {
			requestUrl += "&" + c.query
		} else {
			requestUrl += "?" + c.query
		}
	}

	if c.buildBody == nil {
		c.buildBody = func(context.Context) (body io.Reader, contentType string, err error) {
			return nil, "", nil
		}
	}

	log.Debug("发起请求", "url", requestUrl, "method", c.method)
	if requestUrl == "" {
		return errx.Define("请求地址为空")
	}

	var resp *http.Response
	for i := 0; i < len(c.tryTimes)+1; i++ {
		body, contentType, err := c.buildBody(ctx)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, c.method, requestUrl, body)
		if err != nil {
			return err
		}

		if contentType != "" {
			req.Header.Set(HeaderContentType, contentType)
		}

		for _, headerOption := range c.headers {
			headerOption(req.Header)
		}

		if resp, err = c.client.Do(req); err != nil {
			var ne net.Error
			if i < len(c.tryTimes) && errors.As(err, &ne) {
				log.Debug("返回错误", "try", i+1, "err", err, "delay", c.tryTimes[i])
				select {
				case <-ctx.Done():
					return err
				case <-time.After(c.tryTimes[i]):
					continue
				}
			}
			log.Debug("返回错误", "try", i+1, "err", err)
			return err
		}
		break
	}

	body := resp.Body
	defer errx.Close(body, "close http body")

	proc := Process(func(resp *http.Response) error {
		for _, proc := range process {
			if err := proc(resp); err != nil {
				return err
			}
		}
		return nil
	})

	for _, apply := range c.uses {
		proc = apply(proc)
	}

	return proc(resp)
}

// Bytes 处理响应字节
func (c *Request) Bytes(ctx context.Context) (data []byte, err error) {
	err = c.Process(ctx, func(resp *http.Response) (err error) {
		data, err = io.ReadAll(resp.Body)
		return
	})
	return
}
