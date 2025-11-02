package configx

import (
	"encoding/json"

	"github.com/tidwall/jsonc"
)

var (
	UnmarshalJSON  = json.Unmarshal  // 用于将 JSON 数据反序列化为 Go 值，默认使用标准库 json.Unmarshal。
	MarshalJSON    = json.Marshal    // 用于将 Go 值序列化为 JSON 数据，默认使用标准库 json.Marshal。
	CompactJSON    = json.Compact    // 用于将 JSON 数据紧凑化，默认使用标准库 json.Compact。
	IndentJSON     = json.Indent     // 用于将 JSON 数据缩进化，默认使用标准库 json.Indent。
	HTMLEscape     = json.HTMLEscape // 用于将 JSON 数据转换为 HTML 安全的字符串，默认使用标准库 json.HTMLEscape。
	NewJSONDecoder = json.NewDecoder // 用于创建 JSON 解码器，默认使用标准库 json.NewDecoder。
	NewJSONEncoder = json.NewEncoder // 用于创建 JSON 编码器，默认使用标准库 json.NewEncoder。
)

// UnmarshalJSONC 将 JSONC（带注释的 JSON）数据解析为任意值。
// 它先将 JSONC 转换为标准 JSON，再调用 UnmarshalJSON 完成解析。
func UnmarshalJSONC(data []byte, v any) error {
	return UnmarshalJSON(jsonc.ToJSON(data), v)
}

// YAMLToJSON 将 YAML 格式的字节数组解析为任意值，再将其序列化为 JSON 字节数组返回。
func YAMLToJSON(data []byte) ([]byte, error) {
	var v any
	if err := UnmarshalYAML(data, &v); err != nil {
		return nil, err
	}
	return MarshalJSON(v)
}
