package errx

import (
	"errors"

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
