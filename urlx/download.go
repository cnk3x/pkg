package urlx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cnk3x/pkg/x"
)

const DownloadTempExt = ".downloading"

// Download 下载到文件
func (c *Request) Download(ctx context.Context, fn string, overwrite ...bool) (err error) {
	return c.Process(ctx, func(resp *http.Response) (err error) {
		return downloadFile(resp, fn, len(overwrite) > 0 && overwrite[0])
	})
}

// 下载文件
func downloadFile(resp *http.Response, fn string, overwrite bool) (err error) {
	if err = os.MkdirAll(filepath.Dir(fn), 0755); err != nil {
		return
	}
	var fi os.FileInfo
	if fi, err = os.Stat(fn); err != nil && !os.IsNotExist(err) {
		return
	}
	if fi != nil {
		if fi.IsDir() {
			err = fmt.Errorf("保存路径是文件夹(%s)", fn)
			return
		}
		if !overwrite {
			err = fmt.Errorf("%w: %s", os.ErrExist, fn)
			return
		}
	}

	tempFn := fn + DownloadTempExt
	if err = writeFile(tempFn, resp.Body); err != nil {
		return
	}

	err = os.Rename(tempFn, fn)
	return
}

// 将 body 写入到文件
func writeFile(path string, body io.Reader) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer x.Close(f, "关闭写入的文件")
	_, err = f.ReadFrom(body)
	return err
}
