package strx

func Anys[E ~[]T, T any](s E, appends ...any) (r []any) {
	r = make([]any, len(s), len(s)+len(appends))
	for i, v := range s {
		r[i] = v
	}
	r = append(r, appends...)
	return
}

func AnyMap(vars ...any) (r map[string]any) {
	r = make(map[string]any, len(vars)/2)
	for i := 0; i < len(vars)-1; i += 2 {
		k, ok := vars[i].(string)
		if ok {
			r[k] = vars[i+1]
		}
	}
	return
}

func PairsIndex(pairs []any, tag string) int {
	for i := 0; i < len(pairs)-1; i += 2 {
		if key, ok := pairs[i].(string); ok && key == tag {
			return i
		}
	}
	return -1
}

func PairsFind(pairs []any, tag string) (value any, found bool) {
	if i := PairsIndex(pairs, tag); i > -1 {
		value, found = pairs[i+1], true
	}
	return
}
