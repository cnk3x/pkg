package arrx

import "slices"

// Find 在数组中查找第一个满足条件的元素
//
// 参数:
//   - s: 要搜索的数组
//   - fn:  判断元素是否符合条件的函数
//   - fallback:  可选的默认值，当未找到满足条件的元素时返回
//
// 返回值:
//   - T:     第一个满足条件的元素，如果没有找到则返回零值
//   - bool:  是否找到了满足条件的元素
func Find[T any](s []T, fn func(T) bool, fallback ...T) (r T) {
	for _, v := range s {
		if fn(v) {
			return v
		}
	}
	if len(fallback) > 0 {
		r = fallback[0]
	}
	return r
}

// IndexOf 查找元素在数组中的索引位置
//
// 参数:
//   - s: 要搜索的数组
//   - v: 要查找的元素
//
// 返回值:
//   - int: 元素在数组中的索引，如果不存在则返回-1
func IndexOf[T comparable](s []T, v T) int {
	for i, vv := range s {
		if vv == v {
			return i
		}
	}
	return -1
}

// IndexBy 根据条件函数查找第一个满足条件的元素索引
//
// 参数:
//   - s:  要搜索的数组
//   - fn: 判断元素是否符合条件的函数
//
// 返回值:
//   - int: 第一个满足条件的元素索引，如果不存在则返回-1
func IndexBy[T any](s []T, fn func(T) bool) int {
	for i, v := range s {
		if fn(v) {
			return i
		}
	}
	return -1
}

// Reverse 反转切片中的元素顺序
//
// 参数:
//   - s: 需要被反转的切片
//   - inPlace: 可选的布尔值参数，决定是否就地反转切片。
//     如果不提供或为false，则创建并返回一个新的反转后的切片；
//     如果为true，则直接修改原切片并返回
//
// 返回值:
//   - S: 反转后的新切片或原切片（取决于inPlace参数）
func Reverse[S ~[]E, E any](s S, inPlace ...bool) S {
	if len(inPlace) == 0 || !inPlace[0] {
		s = append(make(S, 0, len(s)), s...)
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Sort 使用给定的比较函数对切片进行稳定排序
//
// 参数:
//   - s: 需要排序的切片
//   - less: 比较函数，用于确定元素的顺序。
//     对于两个元素i和j，如果less(i, j)返回true，则i排在j前面
//   - inPlace: 可选的布尔值参数，决定是否就地排序切片。
//     如果不提供或为false，则创建并返回一个新的排序后的切片；
//     如果为true，则直接修改原切片并返回
//
// 返回值:
//   - S: 排序后的新切片或原切片（取决于inPlace参数）
func Sort[S ~[]E, E any](s S, less func(i, j E) bool, inPlace ...bool) S {
	if len(inPlace) == 0 || !inPlace[0] {
		s = append(make(S, 0, len(s)), s...)
	}
	slices.SortStableFunc(s, func(i, j E) int {
		switch {
		case less(i, j):
			return -1
		case less(j, i):
			return 1
		default:
			return 0
		}
	})
	return s
}
