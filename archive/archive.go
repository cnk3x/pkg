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

var Done = errors.New("done")

type Item interface {
	Context() context.Context
	Index() int
	Open() (io.ReadCloser, error)
	Path() string
	os.FileInfo
}

type ProcessFunc func(item Item) (err error)

func Read(ctx context.Context, source string, itemProcess ProcessFunc) (err error) {
	switch filepath.Ext(source) {
	case ".zip":
		err = readZip(ctx, source, itemProcess)
	case ".tar.gz":
		err = readTgz(ctx, source, itemProcess)
	}

	if err == fs.SkipAll || err == io.EOF || err == Done {
		err = nil
	}

	return
}

func readZip(ctx context.Context, source string, itemProcess ProcessFunc) (err error) {
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
			err = itemProcess(&simpleItem{ctx, i, zr.File[i].Name, zr.File[i].FileInfo(), zr.File[i].Open})
		}
	}
	return
}

func readTgz(ctx context.Context, source string, itemProcess ProcessFunc) (err error) {
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
			err = itemProcess(&simpleItem{ctx, i, it.Name, it.FileInfo(), func() (io.ReadCloser, error) { return io.NopCloser(tr), nil }})
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
