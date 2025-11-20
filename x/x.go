package x

func Iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}
