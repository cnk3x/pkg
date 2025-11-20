package errx

import (
	"errors"
	"io"
	"log/slog"

	"golang.org/x/xerrors"
)

var (
	Define = errors.New     // 创建新错误
	New    = xerrors.New    // 创建新错误
	Is     = errors.Is      // 判断错误是否为指定类型
	As     = errors.As      // 错误类型断言，用于判断错误是否为指定类型
	Unwrap = errors.Unwrap  // 解包错误，用于获取错误链中的底层错误
	Errorf = xerrors.Errorf // 格式化错误
)

// Join 合并多个错误，返回第一个非 nil 错误。
// 如果所有错误都为 nil，则返回 nil。
func Join(errs ...error) (err error) {
	switch errs = Unwraps(errs); len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errors.Join(errs...)
	}
}

// Unwraps 递归解包错误链，返回所有底层错误。
func Unwraps(src []error) (errs []error) {
	for _, err := range src {
		if err != nil {
			if es, ok := err.(interface{ Unwrap() []error }); ok {
				errs = append(errs, Unwraps(es.Unwrap())...)
			} else {
				errs = append(errs, err)
			}
		}
	}
	return
}

// 如果 err != nil, 抛出异常(panic)
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// 如果 err != nil, 返回 zero(T), 只有当 err==nil的时候才返回v，(不管v是不是零值)
func May[T any](v T, err error) (r T) {
	if err != nil {
		errPrint("忽略错误返回", "err", err)
		return
	}
	return v
}

// 如果 err 不是零值，打印输出
func Ig(err error) {
	if err != nil {
		errPrint("忽略错误", "err", err)
	}
}

func Close(closer io.Closer, msg string, args ...any) {
	if err := closer.Close(); err != nil {
		errPrint(msg, append(args, "err", err)...)
	}
}

func Closes(closers ...io.Closer) {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			errPrint("忽略关闭错误", "err", err)
		}
	}
}

var errPrint = slog.With("service", "公共").Debug

func SetErrPrint(printer func(msg string, args ...any)) { errPrint = printer }
