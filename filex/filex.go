package filex

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

// ProcessFunc 定义了对已打开文件进行自定义处理的函数签名
type ProcessFunc func(file *os.File) error

// ProgressFunc 定义了在拷贝过程中报告已拷贝字节数的函数签名
type ProgressFunc func(n int64)

// Cat 一次性读取文件全部内容并以字符串形式返回（内部使用 unsafe 转换，零拷贝）
func Cat(filePath string) string {
	bs, _ := os.ReadFile(filePath)
	return unsafe.String(unsafe.SliceData(bs), len(bs))
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
func WriteFrom(r io.Reader, progress ...ProgressFunc) ProcessFunc {
	return func(file *os.File) error {
		return ProgressCopy(file, r, progress...)
	}
}

// ReadTo 返回一个 ProcessFunc，用于将文件内容读出到外部 Writer
func ReadTo(w io.Writer, progress ...ProgressFunc) ProcessFunc {
	return func(file *os.File) error {
		return ProgressCopy(w, file, progress...)
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
