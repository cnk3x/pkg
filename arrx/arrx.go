package arrx

import (
	"errors"
	"io/fs"
)

func ReplaceIndex[T any](arr []T, fn func(T, int) T) {
	for i, v := range arr {
		arr[i] = fn(v, i)
	}
}

func Replace[T any](arr []T, fn func(T) T) {
	for i, v := range arr {
		arr[i] = fn(v)
	}
}

func Map[T any, R any](list []T, fn func(T) R) []R {
	res := make([]R, len(list))
	for i, v := range list {
		res[i] = fn(v)
	}
	return res
}

func MapIndex[T any, R any](list []T, fn func(T, int) R) []R {
	res := make([]R, len(list))
	for i, v := range list {
		res[i] = fn(v, i)
	}
	return res
}

func FilterIndex[T any](list []T, fn func(T, int) bool) []T {
	res := make([]T, 0)
	for i, v := range list {
		if fn(v, i) {
			res = append(res, v)
		}
	}
	return res
}

func Filter[T any](arr []T, fn func(T) bool) []T {
	res := make([]T, 0)
	for _, v := range arr {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}

func FilterMapIndex[T any, R any](list []T, fn func(T, int) (R, bool)) []R {
	res := make([]R, 0)
	for i, v := range list {
		if u, ok := fn(v, i); ok {
			res = append(res, u)
		}
	}
	return res
}

func FilterMap[T any, R any](arr []T, fn func(T) (R, bool)) []R {
	res := make([]R, 0)
	for _, v := range arr {
		if u, ok := fn(v); ok {
			res = append(res, u)
		}
	}
	return res
}

func ReduceIndex[T any, R any](list []T, fn func(R, T, int) R, init R) R {
	for i, v := range list {
		init = fn(init, v, i)
	}
	return init
}

func Reduce[T any, R any](lst []T, fn func(R, T) R, init R) R {
	for _, v := range lst {
		init = fn(init, v)
	}
	return init
}

// Union 返回所有不重复的对象
func Union[T comparable, S ~[]T](lists ...S) (r S) {
	r = make(S, 0, Sum(lists, Len))
	seen := make(map[T]struct{}, cap(r))
	for _, s := range lists {
		for _, v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				r = append(r, v)
			}
		}
	}
	return
}

// Union 返回所有不重复的对象
func UnionBy[K comparable, S ~[]T, T any](lists []S, gk func(T) K) (r S) {
	r = make(S, 0, Sum(lists, Len))
	seen := make(map[K]struct{}, cap(r))
	for _, s := range lists {
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

func Concat[T any, S ~[]T](lists ...S) S {
	result := make(S, 0, Sum(lists, func(items S) int { return len(items) }))
	Each(lists, func(items S) { result = append(result, items...) })
	return result
}

func Flatten[T any, S ~[]T](lists []S) S { return Concat(lists...) }

func Some[T any](list []T, fn func(T) bool) bool {
	for _, item := range list {
		if fn(item) {
			return false
		}
	}
	return true
}

func All[T any](list []T, fn func(T) bool) bool {
	for _, item := range list {
		if !fn(item) {
			return false
		}
	}
	return true
}

var Break = errors.New("user break walk")

func Walk[T any](list []T, fn func(T) error) (err error) {
	for _, v := range list {
		if err = fn(v); err != nil {
			break
		}
	}
	if err == fs.SkipAll || err == Break {
		err = nil
	}
	return
}

func Each[T any](list []T, fn func(T)) {
	for _, v := range list {
		fn(v)
	}
}

// Uniq returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the slice.
// Play: https://go.dev/play/p/DTzbeXZ6iEN
func Uniq[T comparable, S ~[]T](list S) S {
	result := make(S, 0, len(list))
	seen := make(map[T]struct{}, len(list))
	for i := range list {
		if _, ok := seen[list[i]]; ok {
			continue
		}
		seen[list[i]] = struct{}{}
		result = append(result, list[i])
	}
	return result
}

// UniqBy returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is
// invoked for each element in the slice to generate the criterion by which uniqueness is computed.
// Play: https://go.dev/play/p/g42Z3QSb53u
func UniqBy[T any, U comparable, S ~[]T](collection S, iteratee func(item T) U) S {
	result := make(S, 0, len(collection))
	seen := make(map[U]struct{}, len(collection))

	for i := range collection {
		key := iteratee(collection[i])
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, collection[i])
	}

	return result
}

// GroupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
// Play: https://go.dev/play/p/XnQBd_v6brd
func GroupBy[T any, U comparable, S ~[]T](collection S, iteratee func(item T) U) map[U]S {
	result := map[U]S{}
	for i := range collection {
		key := iteratee(collection[i])
		result[key] = append(result[key], collection[i])
	}
	return result
}

// GroupByMap returns an object composed of keys generated from the results of running each element of collection through iteratee.
// Play: https://go.dev/play/p/iMeruQ3_W80
func GroupByMap[T any, K comparable, V any](collection []T, iteratee func(item T) (K, V)) map[K][]V {
	result := map[K][]V{}
	for i := range collection {
		k, v := iteratee(collection[i])
		result[k] = append(result[k], v)
	}
	return result
}
