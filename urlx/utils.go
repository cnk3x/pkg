package urlx

import (
	"context"
	"fmt"
	"io"
)

// 将无错误返回的函数包装成返回 error 的函数，便于统一处理
func FuncAddE[T any](f func(T)) func(T) error { return func(t T) error { f(t); return nil } }

// 将无错误返回的函数包装成返回 error 的函数，便于统一处理
func Func2AddE[T, T2 any](f func(T, T2)) func(T, T2) error {
	return func(t T, t2 T2) error { f(t, t2); return nil }
}

// 将返回 error 的函数包装成无错误返回的函数，便于统一处理
func FuncDelE[T any](f func(T) error) func(T) { return func(t T) { errIg(f(t)) } }

// 将返回 error 的函数包装成无错误返回的函数，便于统一处理
func Func2DelE[T, T2 any](f func(T, T2) error) func(T, T2) {
	return func(t T, t2 T2) { errIg(f(t, t2)) }
}

func closes(closer io.Closer, log ...Logger) {
	if err := closer.Close(); err != nil {
		for _, l := range log {
			l(context.Background(), fmt.Sprintf("close %T failed", closer), "err", err)
		}
	}
}

func errIg(_ error) {}

func errDel[T any](t T, _ error) T { return t }
