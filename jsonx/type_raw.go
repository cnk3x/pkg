package jsonx

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Raw 原始json数据, 封装了一些常用便捷的操作
type Raw []byte

func (r *Raw) UnmarshalJSON(bytes []byte) error { *r = bytes; return nil }
func (r Raw) MarshalJSON() ([]byte, error)      { return []byte(r), nil }

func (r Raw) get(key string) gjson.Result {
	if key == "" {
		return gjson.ParseBytes(r)
	}
	return gjson.GetBytes(r, key)
}

func (r Raw) IsEmpty() bool                { return len(r) == 0 }
func (r Raw) Exists(key string) bool       { return r.get(key).Exists() }
func (r Raw) Get(key string) Raw           { return Raw(r.get(key).Raw) }
func (r Raw) GetString(key string) string  { return r.get(key).String() }
func (r Raw) GetInt(key string) int64      { return r.get(key).Int() }
func (r Raw) GetFloat(key string) float64  { return r.get(key).Float() }
func (r Raw) GetBool(key string) bool      { return r.get(key).Bool() }
func (r Raw) GetTime(key string) time.Time { return r.get(key).Time() }

func (r Raw) Index(key string) int     { return r.get(key).Index }
func (r Raw) Indexes(key string) []int { return r.get(key).Indexes }

func toRaw(item gjson.Result, _ int) Raw       { return Raw(item.Raw) }
func toString(item gjson.Result, _ int) string { return item.String() }
func toInt(item gjson.Result, _ int) int64     { return item.Int() }
func toFloat(item gjson.Result, _ int) float64 { return item.Float() }
func toBool(item gjson.Result, _ int) bool     { return item.Bool() }

func (r Raw) Gets(key string) []Raw          { return loMap(r.get(key).Array(), toRaw) }
func (r Raw) GetStrings(key string) []string { return loMap(r.get(key).Array(), toString) }
func (r Raw) GetInts(key string) []int64     { return loMap(r.get(key).Array(), toInt) }
func (r Raw) GetFloats(key string) []float64 { return loMap(r.get(key).Array(), toFloat) }
func (r Raw) GetBools(key string) []bool     { return loMap(r.get(key).Array(), toBool) }

// 如果失败则不做更改
func (r Raw) Set(key string, value any) Raw {
	n, err := sjson.SetBytes(r, key, value)
	if err != nil {
		slog.Warn("json set key error", "key", key, "value", value)
		return r
	}
	return n
}

// 如果失败则不做更改
func (r Raw) Del(keys ...string) Raw {
	if len(keys) > 0 && len(r) > 0 {
		obj := gjson.ParseBytes(r)
		for _, key := range keys {
			if key = strings.TrimSpace(key); key == "" || !obj.Get(key).Exists() {
				continue
			}
			if n, err := sjson.DeleteBytes(r, key); err != nil {
				slog.Warn("json delete key error", "key", key)
			} else {
				r = n
			}
		}
	}
	return r
}

func (r Raw) Map(key string) (out map[string]Raw) {
	m := gjson.GetBytes(r, key).Map()
	out = make(map[string]Raw, len(m))
	for k, v := range m {
		out[k] = Raw(v.Raw)
	}
	return
}

func (r Raw) Unmarshal(obj any) error { return Unmarshal(r, obj) }

func (r Raw) Save(file string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
		return err
	}
	return os.WriteFile(file, r, 0666)
}

func (r Raw) WriteTo(w io.Writer) (int64, error) {
	n, e := w.Write(r)
	return int64(n), e
}

func loMap[T, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i := range collection {
		result[i] = iteratee(collection[i], i)
	}

	return result
}
