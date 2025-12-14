package urlx

import (
	"cmp"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/cnk3x/pkg/errx"
	"github.com/cnk3x/pkg/jsonx"
	"golang.org/x/net/http/httpproxy"
)

/* 实现分块下载 */

// func (c *Request) Chuck(ctx context.Context, start, end int64) error {
// 	options := Options(HeaderSet("Range", fmt.Sprintf("bytes=%d-%d", start, end)))

// 	return c.Process(ctx, Process(func(resp *http.Response) error {
// 		if resp.Header.Get("Accept-Ranges") != "bytes" {
// 			return errors.New("不支持分块下载")
// 		}

// 		if resp.StatusCode == 200 {
// 			return errors.New("不支持分块下载")
// 		}

// 		if resp.StatusCode != http.StatusPartialContent {
// 			return errors.New(http.StatusText(resp.StatusCode))
// 		}

// 		// total := resp.ContentLength

// 		return nil
// 	}), options)
// }

type PowerDownload struct {
	url   string           // 下载链接
	dir   string           // 保存目录
	name  string           // 保存文件名
	extra jsonx.Raw        // 扩展信息
	proxy httpproxy.Config // 代理信息

	// maxThreads   int   // 最大线程数
	// maxChunkSize int64 // 最大块大小（字节）

	fileSave string //保存文件
	fileTemp string //临时文件
	fileInfo string //信息文件

	fTemp *os.File // 临时文件句柄

	fileSize       int64  // 文件大小（字节）
	supportsResume bool   // 是否支持断点续传
	location       string // 重定向后的真实下载链接，如果没有重定向则与 url 相同

	downloaded int64 // 已下载字节数
}

func Power(url, dir, name string) *PowerDownload {
	return &PowerDownload{url: url, dir: dir, name: name}
}

// Prepare 准备下载，检查文件是否存在，是否支持断点续传等
func (d *PowerDownload) Prepare(ctx context.Context) error {
	d.fileSave = filepath.Join(d.dir, d.name)
	d.fileTemp = d.fileSave + ".downloading"
	d.fileInfo = d.fileSave + ".info"

	req, err := d.buildReq(ctx)
	if err != nil {
		return err
	}

	client, err := d.buildClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint: errcheck

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errx.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	d.fileSize = resp.ContentLength
	d.supportsResume = resp.Header.Get("Accept-Ranges") == "bytes" && d.fileSize > 0
	d.location = resp.Request.URL.String()
	return nil
}

func (d *PowerDownload) Start(ctx context.Context) error {
	if err := d.Prepare(ctx); err != nil {
		return err
	}

	if d.supportsResume {
		return d.mDownload(ctx)
	}

	os.Remove(d.fileTemp) //nolint: errcheck
	return d.sDownload(ctx)
}

func (d *PowerDownload) createTemp() (f *os.File, err error) {
	if f, err = os.Create(d.fileTemp); err != nil {
		return
	}
	if err = f.Truncate(d.fileSize); err != nil {
		f.Close() //nolint: errcheck
		f = nil
		return
	}
	return
}

func (d *PowerDownload) mDownload(context.Context) error {
	return nil
}

func (d *PowerDownload) sDownload(ctx context.Context) (err error) {
	if d.fTemp, err = d.createTemp(); err != nil {
		return
	}
	defer d.fTemp.Close() //nolint: errcheck

	var req *http.Request
	if req, err = d.buildReq(ctx); err != nil {
		return
	}

	var client *http.Client
	if client, err = d.buildClient(); err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint: errcheck

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errx.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	// w := ProgressWriter(d.fTemp, float64(d.fileSize), func(state ProgressState) {})
	errInvalidWrite := errors.New("invalid write result")
	src := resp.Body
	dst := d.fTemp
	buf := make([]byte, 32*1024)
	var written int64
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
			written += int64(nw)
			atomic.AddInt64(&d.downloaded, int64(nw))
		}

		if er != nil && er != io.EOF {
			err = er
			break
		}
	}

	return
}

func (d *PowerDownload) buildClient() (client *http.Client, err error) {
	client = &http.Client{Transport: http.DefaultTransport}
	if d.proxy.HTTPProxy != "" || d.proxy.HTTPSProxy != "" {
		proxyFunc := d.proxy.ProxyFunc()
		client.Transport.(*http.Transport).Proxy = func(req *http.Request) (*url.URL, error) { return proxyFunc(req.URL) }
	}
	return client, nil
}

func (d *PowerDownload) buildReq(ctx context.Context) (req *http.Request, err error) {
	data := d.extra.GetString("data")
	var params io.Reader
	if data != "" {
		params = strings.NewReader(data)
	}

	method := cmp.Or(d.extra.GetString("method"), http.MethodGet)
	if req, err = http.NewRequestWithContext(ctx, method, d.url, params); err != nil {
		return
	}
	headerSets(d.extra.GetStrings("header"))(req.Header)
	return
}

func headerSets(lines []string) func(header http.Header) {
	return func(header http.Header) {
		for _, line := range lines {
			if line = strings.TrimSpace(line); line != "" {
				switch line[0] {
				case '-':
					for _, k := range strings.Split(line[1:], ",") {
						header.Del(strings.TrimSpace(k))
					}
				case '+':
					if k, v, ok := strings.Cut(line[1:], ":"); ok {
						header.Add(strings.TrimSpace(k), strings.TrimSpace(v))
					}
				default:
					if k, v, ok := strings.Cut(line, ":"); ok {
						header.Set(strings.TrimSpace(k), strings.TrimSpace(v))
					}
				}
			}
		}
	}
}

// func Chunk() {
// 	process := Process(func(resp *http.Response) error {
// 		// 检查状态码
// 		if resp.StatusCode != http.StatusPartialContent {
// 			return fmt.Errorf("unexpected status code: %s", resp.Status)
// 		}
// 	})
// 	return func(c *Request) error {
// 		c.HeaderSet("Range", fmt.Sprintf("bytes=%d-%d", start, end))
// 		return nil
// 	}
// }
