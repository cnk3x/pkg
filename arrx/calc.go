package arrx

// Sum 叠加
func Sum[T any, R Num](arr []T, toNum func(T) R) R {
	return Reduce(arr, func(r R, t T) R { return r + toNum(t) }, 0)
}

// Sum 叠加
func SumIndex[T any, R Num](arr []T, toNum func(T, int) R) R {
	return ReduceIndex(arr, func(r R, t T, i int) R { return r + toNum(t, i) }, 0)
}
