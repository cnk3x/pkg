package jsonx

import (
	"encoding/json"
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
