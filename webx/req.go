package webx

func Clone[T any](r *T, mod ...func(r *T)) *T {
	n := new(T)
	*n = *r
	for _, m := range mod {
		m(n)
	}
	return n
}
