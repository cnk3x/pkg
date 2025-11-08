package types

import (
	"encoding/json/v2"
)

type (
	Strings     = List[string]
	Ints        = List[int64]
	List[T any] []T
)

func (s List[T]) MarshalJSON() ([]byte, error) {
	switch len(s) {
	case 0:
		return json.Marshal([]T{})
	case 1:
		return json.Marshal(s[0])
	default:
		return json.Marshal([]T(s))
	}
}

func (s *List[T]) UnmarshalJSON(data []byte) (err error) {
	var r []T
	if err = json.Unmarshal(data, &r); err == nil {
		*s = r
	}

	var v T
	if err = json.Unmarshal(data, &v); err == nil {
		*s = List[T]{v}
	}
	return
}
