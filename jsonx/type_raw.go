package jsonx

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnk3x/pkg/x"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Raw 原始json数据, 封装了一些常用便捷的操作
type Raw string
type Result = gjson.Result

func Get[T ~string | ~[]byte](src T, key ...string) Raw {
	if len(key) == 0 || key[0] == "" {
		return Raw(src)
	}
	return Raw(gjson.Get(string(src), key[0]).Raw)
}

func New(key string, value any) Raw {
	return Raw(x.Must(sjson.Set(``, key, value)))
}

func (r *Raw) UnmarshalJSON(bytes []byte) error { *r = Raw(bytes); return nil }
func (r Raw) MarshalJSON() ([]byte, error)      { return []byte(r), nil }

func (r Raw) GetRet(key string) Result {
	if key == "" {
		// return gjson.Parse(string(r))
		key = "@this"
	}
	return gjson.Get(string(r), key)
}

func (r Raw) String() string { return string(r) }
func (r Raw) IsEmpty() bool  { return len(r) == 0 }

func (r Raw) Get(key string) Raw { return Raw(r.GetRet(key).Raw) }

func (r Raw) Pretty() Raw { return r.Get("@pretty") }
func (r Raw) Ugly() Raw   { return r.Get("@ugly") }

func (r Raw) Exists(key string) bool      { return r.GetRet(key).Exists() }
func (r Raw) GetString(key string) string { return r.GetRet(key).String() }
func (r Raw) GetInt(key string) int64     { return r.GetRet(key).Int() }
func (r Raw) GetFloat(key string) float64 { return r.GetRet(key).Float() }

// GetBool  获取bool值, 如果不存在, 则返回def
func (r Raw) GetBool(key string, def ...bool) bool {
	if gr := r.GetRet(key); gr.Exists() {
		return gr.Bool()
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func (r Raw) GetDuration(key string) time.Duration {
	s := r.GetString(key)
	d, _ := time.ParseDuration(s)
	return d
}

func (r Raw) GetTime(key string) (dst time.Time) {
	t := gjson.Get(string(r), key)
	switch t.Type {
	default:
		return
	case gjson.String:
		if d, err := time.Parse(time.RFC3339, t.Str); err == nil {
			return d
		}
		switch len(t.Str) {
		case 13:
			switch {
			case strings.Contains(t.Str, "."):
				dst, _ = time.Parse("20060102.1504", t.Str)
			case strings.Contains(t.Str, "-"):
				dst, _ = time.Parse("20060102-1504", t.Str)
			case strings.Contains(t.Str, "/"):
				dst, _ = time.Parse("20060102/1504", t.Str)
			}
		case 19:
			sep := " "
			if strings.Contains(t.Str, "T") {
				sep = "T"
			}
			switch {
			case strings.Contains(t.Str, "-"):
				dst, _ = time.Parse("2006-01-02"+sep+"15:04:05", t.Str)
			case strings.Contains(t.Str, "/"):
				dst, _ = time.Parse("2006/01/02"+sep+"15:04:05", t.Str)
			}
		case 16:
			sep := " "
			if strings.Contains(t.Str, "T") {
				sep = "T"
			}
			switch {
			case strings.Contains(t.Str, "-"):
				dst, _ = time.Parse("2006-01-02"+sep+"15:04", t.Str)
			case strings.Contains(t.Str, "/"):
				dst, _ = time.Parse("2006/01/02"+sep+"15:04", t.Str)
			}
		case 10:
			switch {
			case strings.Contains(t.Str, "-"):
				dst, _ = time.Parse("2006-01-02", t.Str)
			case strings.Contains(t.Str, "/"):
				dst, _ = time.Parse("2006/01/02", t.Str)
			case strings.Contains(t.Str, "."):
				dst, _ = time.Parse("2006.01.02", t.Str)
			}
		case 14:
			dst, _ = time.Parse("20060102150405", t.Str)
		case 12:
			dst, _ = time.Parse("200601021504", t.Str)
		case 15:
			dst, _ = time.Parse("20060102.150405", t.Str)
			switch {
			case strings.Contains(t.Str, "."):
				dst, _ = time.Parse("20060102.150405", t.Str)
			case strings.Contains(t.Str, "-"):
				dst, _ = time.Parse("20060102-150405", t.Str)
			case strings.Contains(t.Str, "/"):
				dst, _ = time.Parse("20060102/150405", t.Str)
			}
		}
		return
	case gjson.Number:
		if t.Num == 0 {
			return
		}

		//小数或者10位整数,毫秒时间戳
		if strings.Contains(t.Raw, ".") || len(t.Raw) == 10 {
			//1970年1月12日 21:46:40 - 1970年4月27日 01:46:39
			dst = time.UnixMilli(int64(t.Num * 1000))
			return
		}

		// 重叠时间段
		// 2572年2月2日 21:50:00 - 2762年3月23日 08:30:00
		// 1970年8月9日 05:48:21 - 1970年10月17日 16:28:21
		if t.Num >= 19000101000000 && t.Num < 25000101000000 {
			s := t.Raw
			if d, err := time.Parse("20060102150405", s); err == nil {
				return d
			}
			return
		}

		dst = time.UnixMilli(int64(t.Num * 1000))
		return
	}
}

func (r Raw) Index(key string) int     { return r.GetRet(key).Index }
func (r Raw) Indexes(key string) []int { return r.GetRet(key).Indexes }

func (r Raw) Gets(key string) []Raw          { return loMap(r.GetRet(key).Array(), toRaw) }
func (r Raw) GetStrings(key string) []string { return loMap(r.GetRet(key).Array(), toString) }
func (r Raw) GetInts(key string) []int64     { return loMap(r.GetRet(key).Array(), toInt) }
func (r Raw) GetFloats(key string) []float64 { return loMap(r.GetRet(key).Array(), toFloat) }
func (r Raw) GetBools(key string) []bool     { return loMap(r.GetRet(key).Array(), toBool) }

func toRaw(item gjson.Result) Raw       { return Raw(item.Raw) }
func toString(item gjson.Result) string { return item.String() }
func toInt(item gjson.Result) int64     { return item.Int() }
func toFloat(item gjson.Result) float64 { return item.Float() }
func toBool(item gjson.Result) bool     { return item.Bool() }

// 如果失败则不做更改
func (r Raw) Set(key string, value any) Raw {
	n, err := sjson.Set(string(r), key, value)
	if err != nil {
		slog.Warn("json set key error", "key", key, "value", value)
		return r
	}
	return Raw(n)
}

// 如果失败则不做更改
func (r Raw) Del(keys ...string) Raw {
	if len(keys) > 0 && len(r) > 0 {
		obj := gjson.Parse(string(r))
		for _, key := range keys {
			if key = strings.TrimSpace(key); key == "" || !obj.Get(key).Exists() {
				continue
			}
			if n, err := sjson.Delete(string(r), key); err != nil {
				slog.Warn("json delete key error", "key", key)
			} else {
				r = Raw(n)
			}
		}
	}
	return r
}

func (r Raw) Map(key string) (out map[string]Raw) {
	m := gjson.Get(string(r), key).Map()
	out = make(map[string]Raw, len(m))
	for k, v := range m {
		out[k] = Raw(v.Raw)
	}
	return
}

func (r Raw) Unmarshal(obj any) error { return Unmarshal([]byte(r), obj) }

func (r Raw) Save(file string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
		return err
	}
	return os.WriteFile(file, []byte(r), 0666)
}

func (r Raw) WriteTo(w io.Writer) (int64, error) {
	n, e := io.WriteString(w, string(r))
	return int64(n), e
}

func loMap[T, R any](collection []T, iteratee func(item T) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item)
	}

	return result
}

func FieldIs(field, value string) func(item Raw) bool {
	return func(item Raw) bool { return item.GetString(field) == value }
}

func KeyIs(field string) func(item Raw) string {
	return func(node Raw) string { return node.GetString(field) }
}
