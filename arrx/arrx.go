package arrx

import (
	"errors"
	"io/fs"
	"slices"
)

// Inserts 将切片v中的所有元素插入到切片s的指定位置i处(原切片被修改)
//
// 参数:
//   - s: 原始切片
//   - i: 插入位置索引
//   - v: 要插入的元素切片
//
// 返回值:
//   - 插入元素后的新切片
func Inserts[S ~[]T, T any](s S, i int, v []T) S {
	return slices.Insert(s, i, v...)
}

// Insert 将可变参数v中的所有元素插入到切片s的指定位置i处(原切片被修改)
//
// 参数:
//   - s: 原始切片
//   - i: 插入位置索引
//   - v: 要插入的元素（可变参数）
//
// 返回值:
//   - 插入元素后的新切片
func Insert[S ~[]T, T any](s S, i int, v ...T) S {
	return Inserts(s, i, v)
}

// Replace 对数组中的每个元素应用函数fn，并就地更新数组(原切片被修改)
//
// 参数:
//   - s: 要处理的数组
//   - fn: 应用于每个元素的函数
func Replace[T any](s []T, fn func(T) T) {
	for i, v := range s {
		s[i] = fn(v)
	}
}

// ReplaceIndex 对数组中的每个元素及其索引应用函数fn，并就地更新数组(原切片被修改)
//
// 参数:
//   - s: 要处理的数组
//   - fn: 应用于每个元素及索引的函数
func ReplaceIndex[T any](s []T, fn func(T, int) T) {
	for i, v := range s {
		s[i] = fn(v, i)
	}
}

// Map 对数组中的每个元素应用函数fn并返回结果数组
//
// 参数:
//   - s: 要处理的数组
//   - fn: 应用于每个元素的函数
//
// 返回值:
//   - 包含应用函数后结果的新数组
func Map[T any, R any](s []T, fn func(T) R) (r []R) {
	r = make([]R, len(s))
	for i, v := range s {
		r[i] = fn(v)
	}
	return
}

// MapIndex 对数组中的每个元素及其索引应用函数fn并返回结果数组
//
// 参数:
//   - s: 要处理的数组
//   - fn: 应用于每个元素及索引的函数
//
// 返回值:
//   - 包含应用函数后结果的新数组
func MapIndex[T any, R any](s []T, fn func(T, int) R) (r []R) {
	r = make([]R, len(s))
	for i, v := range s {
		r[i] = fn(v, i)
	}
	return
}

// Filter 根据谓词函数过滤数组元素
//
// 参数:
//   - s: 要过滤的数组
//   - fn: 谓词函数，返回true表示保留该元素
//
// 返回值:
//   - 过滤后的新数组
func Filter[S ~[]T, T any](s S, fn func(T) bool) S {
	r := make([]T, 0)
	for _, v := range s {
		if fn(v) {
			r = append(r, v)
		}
	}
	return r
}

// FilterIndex 根据带索引的谓词函数过滤数组元素
//
// 参数:
//   - s: 要过滤的数组
//   - fn: 带索引的谓词函数，返回true表示保留该元素
//
// 返回值:
//   - 过滤后的新数组
func FilterIndex[S ~[]T, T any](s S, fn func(T, int) bool) S {
	r := make([]T, 0)
	for i, v := range s {
		if fn(v, i) {
			r = append(r, v)
		}
	}
	return r
}

// FilterMapIndex 根据带索引的函数过滤并映射数组元素
//
// 参数:
//   - s: 要处理的数组
//   - fn: 带索引的函数，返回映射值和布尔值，布尔值为true时保留该元素
//
// 返回值:
//   - 过滤并映射后的新数组
func FilterMapIndex[T any, R any](s []T, fn func(T, int) (R, bool)) []R {
	r := make([]R, 0)
	for i, v := range s {
		if u, ok := fn(v, i); ok {
			r = append(r, u)
		}
	}
	return r
}

// FilterMap 根据函数过滤并映射数组元素
//
// 参数:
//   - s: 要处理的数组
//   - fn: 函数，返回映射值和布尔值，布尔值为true时保留该元素
//
// 返回值:
//   - 过滤并映射后的新数组
func FilterMap[T any, R any](s []T, fn func(T) (R, bool)) []R {
	r := make([]R, 0)
	for _, v := range s {
		if u, ok := fn(v); ok {
			r = append(r, u)
		}
	}
	return r
}

// ReduceIndex 使用带索引的累积函数对数组进行归约操作
//
// 参数:
//   - s: 要归约的数组
//   - fn: 带索引的累积函数，接收累积值、当前元素和索引，返回新的累积值
//   - init: 初始累积值
//
// 返回值:
//   - 归约后的最终值
func ReduceIndex[T any, R any](s []T, fn func(R, T, int) R, init R) R {
	for i, v := range s {
		init = fn(init, v, i)
	}
	return init
}

// Reduce 使用累积函数对数组进行归约操作
//
// 参数:
//   - s: 要归约的数组
//   - fn: 累积函数，接收累积值和当前元素，返回新的累积值
//   - init: 初始累积值
//
// 返回值:
//   - 归约后的最终值
func Reduce[T any, R any](s []T, fn func(R, T) R, init R) R {
	for _, v := range s {
		init = fn(init, v)
	}
	return init
}

// ReduceIf 使用累积函数对数组进行归约操作，当函数返回false时停止
//
// 参数:
//   - s: 要归约的数组
//   - fn: 累积函数，接收累积值和当前元素，返回新的累积值和是否继续的布尔值
//   - init: 初始累积值
//
// 返回值:
//   - 归约后的最终值
func ReduceIf[T any, R any](s []T, fn func(R, T) (R, bool), init R) R {
	for _, v := range s {
		acc, next := fn(init, v)
		if !next {
			break
		}
		init = acc
	}
	return init
}

// Union 合并多个切片并去除重复元素
//
// 参数:
//   - ss: 要合并的切片列表
//
// 返回值:
//   - 去重后的合并切片
func Union[S ~[]T, T comparable](ss ...S) (r S) {
	r = make(S, 0, SumBy(ss, Len))
	seen := make(map[T]struct{}, cap(r))
	for _, s := range ss {
		for _, v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				r = append(r, v)
			}
		}
	}
	return
}

// UnionBy 根据指定的键函数合并多个切片并去除重复元素
//
// 参数:
//   - ss: 要合并的切片列表
//   - gk: 生成元素键值的函数
//
// 返回值:
//   - 去重后的合并切片
func UnionBy[K comparable, S ~[]T, T any](ss []S, gk func(T) K) (r S) {
	r = make(S, 0, SumBy(ss, Len))
	seen := make(map[K]struct{}, cap(r))
	for _, s := range ss {
		for _, v := range s {
			k := gk(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				r = append(r, v)
			}
		}
	}
	return
}

// Concat 连接多个切片
//
// 参数:
//   - ss: 要连接的切片列表
//
// 返回值:
//   - 连接后的新切片
func Concat[S ~[]T, T any](ss ...S) S {
	result := make(S, 0, SumBy(ss, Len))
	return Reduce(ss, func(r S, items S) S {
		return append(r, items...)
	}, result)
}

// Flatten 展平切片的切片为单个切片
//
// 参数:
//   - ss: 切片的切片
//
// 返回值:
//   - 展平后的一维切片
func Flatten[S ~[]T, T any](ss []S) S { return Concat(ss...) }

// Some 检查数组中是否存在至少一个元素满足谓词函数
// 当找到第一个满足条件的元素时立即停止迭代
//
// 参数:
//   - s: 要检查的元素切片
//   - fn: 谓词函数，接收一个元素并返回布尔值
//
// 返回值:
//   - 如果至少有一个元素满足谓词函数则返回true，否则返回false
func Some[T any](s []T, fn func(T) bool) bool {
	return slices.ContainsFunc(s, fn)
}

// Every 检查数组中是否所有元素都满足谓词函数
//
// 参数:
//   - s: 要检查的数组
//   - fn: 谓词函数，返回true表示元素满足条件
//
// 返回值:
//   - 如果所有元素都满足条件则返回true，否则返回false
func Every[T any](s []T, fn func(T) bool) bool {
	for _, item := range s {
		if !fn(item) {
			return false
		}
	}
	return true
}

// All 检查数组中是否所有元素都满足谓词函数（与Every功能相同）
//
// 参数:
//   - s: 要检查的数组
//   - fn: 谓词函数，返回true表示元素满足条件
//
// 返回值:
//   - 如果所有元素都满足条件则返回true，否则返回false
func All[T any](s []T, fn func(T) bool) bool {
	for _, item := range s {
		if !fn(item) {
			return false
		}
	}
	return true
}

// ErrBreak 用于中断Walk函数遍历的错误
var ErrBreak = errors.New("user break walk")

// Walk 遍历数组并对每个元素执行函数，遇到错误时停止
//
// 参数:
//   - arr: 要遍历的数组
//   - fn: 对每个元素执行的函数，返回error表示是否中断遍历
//
// 返回值:
//   - 遍历过程中产生的错误（如果有的话）
func Walk[T any](s []T, fn func(T) error) (err error) {
	for _, v := range s {
		if err = fn(v); err != nil {
			break
		}
	}
	if err == fs.SkipAll || err == ErrBreak {
		err = nil
	}
	return
}

// Each 遍历数组并对每个元素执行函数
//
// 参数:
//   - s: 要遍历的数组
//   - fn: 对每个元素执行的函数
func Each[T any](s []T, fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

// EachIndex 遍历数组并对每个元素及其索引执行函数
//
// 参数:
//   - s: 要遍历的数组
//   - fn: 对每个元素及索引执行的函数
func EachIndex[T any](s []T, fn func(T, int)) {
	for i, v := range s {
		fn(v, i)
	}
}

// Uniq 去除切片中的重复元素，只保留首次出现的元素
//
// 参数:
//   - s: 包含可能重复元素的切片
//
// 返回值:
//   - 去重后的新切片
func Uniq[S ~[]T, T comparable](s S) S {
	result := make(S, 0, len(s))
	seen := make(map[T]struct{}, len(s))
	for i := range s {
		if _, ok := seen[s[i]]; ok {
			continue
		}
		seen[s[i]] = struct{}{}
		result = append(result, s[i])
	}
	return result
}

// UniqBy 根据指定的函数对切片元素去重，只保留首次出现的元素
//
// 参数:
//   - s: 包含可能重复元素的切片
//   - fn: 生成元素键值的函数
//
// 返回值:
//   - 去重后的新切片
func UniqBy[U comparable, S ~[]T, T any](s S, fn func(item T) U) S {
	result := make(S, 0, len(s))
	seen := make(map[U]struct{}, len(s))

	for i := range s {
		key := fn(s[i])
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, s[i])
	}

	return result
}

// GroupBy 根据指定的函数对切片元素进行分组
//
// 参数:
//   - collection: 要分组的切片
//   - fn: 生成元素分组键的函数
//
// 返回值:
//   - 以分组键为键、对应元素切片为值的映射
func GroupBy[U comparable, S ~[]T, T any](s S, fn func(item T) U) map[U]S {
	result := map[U]S{}
	for i := range s {
		key := fn(s[i])
		result[key] = append(result[key], s[i])
	}
	return result
}

// GroupByMap 根据指定的函数对切片元素进行分组并映射值
//
// 参数:
//   - s: 要分组的切片
//   - fn: 生成分组键和映射值的函数
//
// 返回值:
//   - 以分组键为键、对应映射值切片为值的映射
func GroupByMap[T any, K comparable, V any](s []T, fn func(item T) (K, V)) map[K][]V {
	result := map[K][]V{}
	for i := range s {
		k, v := fn(s[i])
		result[k] = append(result[k], v)
	}
	return result
}
