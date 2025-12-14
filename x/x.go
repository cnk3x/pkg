package x

func Iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}

func Clone[T any](r *T, mod ...func(r *T)) *T {
	n := new(T)
	*n = *r
	for _, m := range mod {
		m(n)
	}
	return n
}

func ReduceIf[T, R any](items []T, fn func(r R, item T) (v R, next bool), init R) (v R) {
	var next bool
	for _, item := range items {
		if init, next = fn(init, item); !next {
			return
		}
	}
	return init
}
