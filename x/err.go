package x

import (
	"io"
	"log/slog"
)

// 如果 err 不是零值，打印输出
func Ig(err error) {
	if err != nil {
		errPrint("忽略错误", "err", err)
	}
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
