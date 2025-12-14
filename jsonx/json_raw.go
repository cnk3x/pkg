package jsonx

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/tidwall/jsonc"
)

func FromE(v any) (Raw, error) {
	r, err := Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal json error: %w", err)
	}
	return Raw(r), nil
}

func From(v any) Raw { return eDebug(FromE(v)) }

func ParseYamlE(yaml []byte) (Raw, error) {
	w := &bytes.Buffer{}
	if err := translateStream(bytes.NewReader(yaml), w); err != nil {
		return "", fmt.Errorf("parse yaml to json error: %w", err)
	}
	return Raw(w.String()), nil
}

func ParseYaml(yaml []byte) Raw { return eDebug(ParseYamlE(yaml)) }

func LoadE(file string) (Raw, error) {
	d, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("load json file error: %w", err)
	}
	return Raw(jsonc.ToJSON(d)), nil
}

func Load(file string) Raw { return eDebug(LoadE(file)) }

func LoadYamlE(file string) (Raw, error) {
	d, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("load yaml file error: %w", err)
	}
	return ParseYamlE(d)
}

func LoadYaml(file string) Raw { return eDebug(LoadYamlE(file)) }

func eDebug(r Raw, err error) Raw {
	if err != nil {
		slog.Debug(err.Error())
		return ""
	}
	return r
}

func Read(r io.Reader) (Raw, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return Raw(data), nil
}
