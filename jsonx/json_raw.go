package jsonx

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/tidwall/jsonc"
)

func From(v any) Raw {
	d, err := Marshal(v)
	if err != nil {
		return nil
	}
	return d
}

func FromYAML(yaml []byte) Raw {
	w := &bytes.Buffer{}
	if err := translateStream(bytes.NewReader(yaml), w); err != nil {
		slog.Debug("yaml to json error", "err", err)
		return nil
	}
	return w.Bytes()
}

func FromFile(file string) Raw {
	d, err := os.ReadFile(file)
	if err != nil {
		slog.Debug("load json file error", "file", file, "err", err)
		return nil
	}
	return jsonc.ToJSON(d)
}
