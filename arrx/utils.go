package arrx

// 可变参数固定
func Variadic[T any](arr ...T) []T {
	return arr
}

// Ni -- No Index
func Ni[T, R any](f func(T) R) func(T, int) R {
	return func(t T, _ int) R { return f(t) }
}

func NiB[T, R any](f func(T) (R, bool)) func(T, int) (R, bool) {
	return func(t T, _ int) (R, bool) { return f(t) }
}

func Ni3[T, U, R any](f func(T, U) R) func(T, U, int) R {
	return func(t T, u U, _ int) R { return f(t, u) }
}

func Len[S ~[]T, T any](s S) int { return len(s) }
