package filex

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cnk3x/pkg/x"
)

// ProcessFunc 定义了对已打开文件进行自定义处理的函数签名
type ProcessFunc func(file *os.File) error

// ProgressFunc 定义了在拷贝过程中报告已拷贝字节数的函数签名
type ProgressFunc func(n int64)

func ProgressJoin(progress ...ProgressFunc) ProgressFunc {
	var n bool
	for _, p := range progress {
		if p != nil {
			n = true
			break
		}
	}
	if !n {
		return nil
	}
	return func(n int64) {
		for _, p := range progress {
			if p != nil {
				p(n)
			}
		}
	}
}

// Cat 一次性读取文件全部内容并以字符串形式返回
func Cat(filePath string, trimSpace ...bool) string { return x.May(CatE(filePath, trimSpace...)) }

func CatE(filePath string, trimSpace ...bool) (string, error) {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file %s error: %w", filePath, err)
	}
	if len(trimSpace) > 0 && trimSpace[0] {
		bs = bytes.TrimSpace(bs)
	}
	return string(bs), nil
}

// WriteText 将文本写入指定文件，支持通过 opts 控制是否自动创建目录及文件权限
func WriteText(filePath string, text string, opts ...Option) (err error) {
	options := applyOptions(opts...)
	if options.createDirs {
		if err = os.MkdirAll(filepath.Dir(filePath), options.createDirsMode); err != nil {
			return
		}
	}
	err = os.WriteFile(filePath, []byte(text), options.createMode)
	return
}

func Copy(src, dst string) (err error) {
	return Process(src, func(r *os.File) error {
		return Process(dst, func(w *os.File) error {
			return CopyPipe(context.Background(), w, r, nil)
		})
	}, Readonly())
}

// Open 根据 opts 选项打开或创建文件，支持只读、覆盖、追加、排他创建等模式
func Open(filePath string, opts ...Option) (file *os.File, err error) {
	options := applyOptions(opts...)
	if options.createDirs {
		if err = os.MkdirAll(filepath.Dir(filePath), options.createDirsMode); err != nil {
			return
		}
	}

	if options.readonly {
		file, err = os.Open(filePath)
		return
	}

	createFlag := os.O_CREATE | os.O_RDWR
	switch {
	case options.overwrite:
		createFlag |= os.O_TRUNC
	case options.append:
		createFlag |= os.O_APPEND
	default:
		createFlag |= os.O_EXCL
	}

	file, err = os.OpenFile(filePath, createFlag, options.createMode)
	return
}

// Process 打开文件后立即执行传入的 processFunc，并确保文件关闭
func Process(filePath string, processFunc ProcessFunc, opts ...Option) (err error) {
	file, err := Open(filePath, opts...)
	if err != nil {
		return
	}
	defer file.Close()
	return processFunc(file)
}

// WriteFrom 返回一个 ProcessFunc，用于将外部 Reader 数据写入文件
func WriteFrom(ctx context.Context, r io.Reader, progress ...ProgressFunc) ProcessFunc {
	return func(file *os.File) error {
		return CopyPipe(ctx, file, r, ProgressJoin(progress...))
	}
}

// ReadTo 返回一个 ProcessFunc，用于将文件内容读出到外部 Writer
func ReadTo(ctx context.Context, w io.Writer, progress ...ProgressFunc) ProcessFunc {
	return func(file *os.File) error {
		return CopyPipe(ctx, w, file, ProgressJoin(progress...))
	}
}

// CalcMD5 返回一个 ProcessFunc，用于计算文件 MD5 值并写入 hexMd5 指针指向的字符串
func CalcMD5(hexMd5 *string) ProcessFunc {
	return func(file *os.File) (err error) {
		h := md5.New()
		_, err = io.Copy(h, file)
		if err != nil {
			return
		}
		*hexMd5 = hex.EncodeToString(h.Sum(nil))
		return
	}
}
