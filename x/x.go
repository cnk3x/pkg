package x

func Iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}

func CloneBy[T any](r *T, mod ...func(r *T)) *T {
	n := new(T)
	*n = *r
	for _, m := range mod {
		m(n)
	}
	return n
}
