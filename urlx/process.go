package urlx

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

/*处理响应*/

type Process = func(resp *http.Response) error // 响应处理器

// Use 添加中间件（Process前置处理）
func (c *Request) Use(processes ...Process) *Request { c.uses = append(c.uses, processes...); return c }

// Process 处理响应
func (c *Request) Process(ctx context.Context, process Process, options ...Option) error {

	for _, apply := range c.options {
		if err := apply(c); err != nil {
			return err
		}
	}

	for _, apply := range options {
		if err := apply(c); err != nil {
			return err
		}
	}

	client := &http.Client{Transport: transportDefault()}

	for _, cliOpt := range c.clientOptions {
		if err := cliOpt(client); err != nil {
			return err
		}
	}

	method := c.method
	if method == "" {
		method = http.MethodGet
	}

	requestUrl := c.url

	log := c.logger()
	log(ctx, "发起请求", "url", requestUrl, "method", method)
	if requestUrl == "" {
		return errors.New("请求地址为空")
	}

	var resp *http.Response
	for i := 0; i < len(c.tryTimes)+1; i++ {
		var body io.Reader
		var contentType string
		var err error

		if c.body != nil {
			if body, contentType, err = c.body(ctx); err != nil {
				return err
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, requestUrl, body)
		if err != nil {
			return err
		}

		if contentType != "" {
			req.Header.Set(HeaderContentType, contentType)
		}

		for _, headerProcess := range c.headers {
			headerProcess(req.Header)
		}

		if resp, err = client.Do(req); err != nil {
			var ne net.Error
			if i < len(c.tryTimes) && errors.As(err, &ne) {
				log(ctx, "返回错误", "try", i+1, "delay", c.tryTimes[i], "err", err)
				select {
				case <-ctx.Done():
					return err
				case <-time.After(c.tryTimes[i]):
					continue
				}
			}
			log(ctx, "返回错误", "try", i+1, "err", err)
			return err
		}
		break
	}

	body := resp.Body
	defer closes(body, log)

	for _, proc := range c.uses {
		if err := proc(resp); err != nil {
			return err
		}
	}

	return process(resp)
}

// Bytes 处理响应字节
func (c *Request) Bytes(ctx context.Context) (data []byte, err error) {
	readBytes := func(resp *http.Response) (err error) { data, err = io.ReadAll(resp.Body); return }
	err = c.Process(ctx, readBytes)
	return
}

func ReadBytes(read func(data []byte) error) Process {
	return func(resp *http.Response) error {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return read(data)
	}
}

func ReadHeader(read func(headers http.Header) error) Process {
	return func(resp *http.Response) error {
		return read(resp.Header)
	}
}
