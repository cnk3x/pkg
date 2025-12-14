package jsonx

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

var (
	Unmarshal = json.Unmarshal
	Marshal   = json.Marshal
)

func UnmarshalFromFile(file string, v any) error {
	d, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return Unmarshal(d, v)
}

func MarshalToFile(file string, v any) error {
	d, err := Marshal(v)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return os.WriteFile(file, d, 0644)
}

func Decode(r io.Reader, v any, limitSize ...int64) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	if len(limitSize) > 0 && limitSize[0] > 0 {
		r = io.LimitReader(r, limitSize[0])
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return Unmarshal(data, v)
}
