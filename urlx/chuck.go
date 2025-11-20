package urlx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

/* 实现分块下载 */

func (c *Request) Chuck(ctx context.Context, start, end int64) error {
	options := Options(HeaderSet("Range", fmt.Sprintf("bytes=%d-%d", start, end)))

	return c.Process(ctx, Process(func(resp *http.Response) error {
		if resp.Header.Get("Accept-Ranges") != "bytes" {
			return errors.New("不支持分块下载")
		}

		if resp.StatusCode == 200 {
			return errors.New("不支持分块下载")
		}

		if resp.StatusCode != http.StatusPartialContent {
			return errors.New(http.StatusText(resp.StatusCode))
		}

		// total := resp.ContentLength

		return nil
	}), options)
}
