package jsonx

import "log/slog"

type (
	Strings     = List[string]
	Ints        = List[int64]
	List[T any] []T
)

func (s List[T]) MarshalJSON() ([]byte, error) {
	switch len(s) {
	case 0:
		return Marshal([]T{})
	case 1:
		return Marshal(s[0])
	default:
		return Marshal([]T(s))
	}
}

func (s *List[T]) UnmarshalJSON(data []byte) (err error) {
	slog.Debug("list^T", "data", string(data))
	var r []T
	if err = Unmarshal(data, &r); err == nil {
		*s = r
		return
	}

	var v T
	if err = Unmarshal(data, &v); err == nil {
		*s = List[T]{v}
	}
	return
}
