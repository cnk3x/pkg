package arrx

import "slices"

func Find[T any](arr []T, fn func(T) bool) (T, bool) {
	for _, v := range arr {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

func FindOrZero[T any](arr []T, fn func(T) bool) T {
	for _, v := range arr {
		if fn(v) {
			return v
		}
	}
	var zero T
	return zero
}

func IndexOf[T comparable](arr []T, v T) int {
	for i, vv := range arr {
		if vv == v {
			return i
		}
	}
	return -1
}

// Reverse 反转切片, inPlace 为 true 时, 会修改原切片然后返回(原切片); 为 false 时, 会返回新切片
func Reverse[S ~[]E, E any](arr S, inPlace bool) S {
	if !inPlace {
		arr = append(make(S, 0, len(arr)), arr...)
	}
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// Sort 排序切片, inPlace 为 true 时, 会修改原切片然后返回(原切片); 为 false 时, 会返回新切片
func Sort[S ~[]E, E any](arr S, less func(i, j E) bool, inPlace bool) S {
	if !inPlace {
		arr = append(make(S, 0, len(arr)), arr...)
	}
	slices.SortStableFunc(arr, func(i, j E) int {
		if less(i, j) {
			return -1
		}
		if less(j, i) {
			return 1
		}
		return 0
	})
	return arr
}
