package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrDone 是一个用于标记处理完成的错误标识
var ErrDone = errors.New("done")

// Item 接口定义了归档文件中单个项目的规范，包含了获取上下文、索引、打开文件、路径以及文件信息的方法
type Item interface {
	// Context 返回与此项目关联的上下文
	Context() context.Context
	// Index 返回项目在归档文件中的索引位置
	Index() int
	// Open 打开项目并返回可读取的流
	Open() (io.ReadCloser, error)
	// Path 返回项目的路径
	Path() string
	os.FileInfo
}

// ProcessFunc 定义了处理归档文件中项目的函数类型
type ProcessFunc func(ctx context.Context, item Item) (err error)

// Read 根据文件扩展名识别归档文件格式并读取其中的内容，支持 .zip 和 .tar.gz 格式
// 参数:
//   - ctx: 上下文，用于控制处理流程和超时取消
//   - source: 归档文件的路径
//   - process: 处理每个归档项的函数
//
// 返回值:
//   - 处理过程中发生的错误，如果正常结束或遇到EOF、SkipAll、ErrDone则返回nil
func Read(ctx context.Context, source string, process ProcessFunc) (err error) {
	switch filepath.Ext(source) {
	case ".zip":
		err = readZip(ctx, source, process)
	case ".tar.gz":
		err = readTgz(ctx, source, process)
	}

	if err == fs.SkipAll || err == io.EOF || err == ErrDone {
		err = nil
	}

	return
}

func readZip(ctx context.Context, source string, process ProcessFunc) (err error) {
	zr, ze := zip.OpenReader(source)
	if err = ze; err != nil {
		return
	}
	defer zr.Close()
	for i := 0; err == nil && i < len(zr.File); i++ {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			err = process(ctx, &simpleItem{ctx, i, zr.File[i].Name, zr.File[i].FileInfo(), zr.File[i].Open})
		}
	}
	return
}

func readTgz(ctx context.Context, source string, process ProcessFunc) (err error) {
	fr, fe := os.Open(source)
	if err = fe; err != nil {
		return
	}
	defer fr.Close()

	zr, ze := gzip.NewReader(fr)
	if err = ze; err != nil {
		return
	}
	defer zr.Close()

	for tr, i := tar.NewReader(zr), 0; err == nil; i++ {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			it, e := tr.Next()
			if err = e; err != nil {
				return
			}
			err = process(ctx, &simpleItem{ctx, i, it.Name, it.FileInfo(), func() (io.ReadCloser, error) { return io.NopCloser(tr), nil }})
		}
	}

	return
}

type simpleItem struct {
	ctx   context.Context
	index int
	path  string
	os.FileInfo
	open func() (io.ReadCloser, error)
}

func (item *simpleItem) Context() context.Context     { return item.ctx }
func (item *simpleItem) Index() int                   { return item.index }
func (item *simpleItem) Open() (io.ReadCloser, error) { return item.open() }
func (item *simpleItem) Path() string                 { return item.path }
