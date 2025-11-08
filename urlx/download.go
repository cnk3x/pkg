package urlx

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const DownloadTempExt = ".uxdt"

// Download 下载到文件
func (c *Request) Download(fn *string, overwrite ...bool) (err error) {
	return c.Process(func(resp *http.Response) (err error) {
		return downloadFile(resp, fn, len(overwrite) > 0 && overwrite[0])
	})
}

// 下载文件
func downloadFile(resp *http.Response, fn *string, overwrite bool) (err error) {
	if err = os.MkdirAll(filepath.Dir(*fn), 0755); err != nil {
		return
	}
	var fi os.FileInfo
	if fi, err = os.Stat(*fn); err != nil && !os.IsNotExist(err) {
		return
	}
	if fi != nil {
		if fi.IsDir() {
			*fn = filepath.Join(*fn, filepath.Base(*fn))
			return downloadFile(resp, fn, overwrite)
		}
		if !overwrite {
			err = fmt.Errorf("%w: %s", os.ErrExist, *fn)
			return
		}
	}

	tempFn := *fn + DownloadTempExt
	if err = writeFile(tempFn, resp.Body); err != nil {
		return
	}

	err = os.Rename(tempFn, *fn)
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
	defer closes(f)
	_, err = f.ReadFrom(body)
	return err
}
