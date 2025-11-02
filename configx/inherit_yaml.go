package configx

import (
	"go.yaml.in/yaml/v4"
)

var (
	UnmarshalYAML  = yaml.Unmarshal  // 用于将 YAML 数据反序列化为 Go 值，默认使用 `go.yaml.in/yaml/v4`.Unmarshal。
	MarshalYAML    = yaml.Marshal    // 用于将 Go 值序列化为 YAML 数据，默认使用 `go.yaml.in/yaml/v4`.Marshal。
	NewYAMLDecoder = yaml.NewDecoder // 用于创建 YAML 解码器，默认使用 `go.yaml.in/yaml/v4`.NewDecoder。
	NewYAMLEncoder = yaml.NewEncoder // 用于创建 YAML 编码器，默认使用 `go.yaml.in/yaml/v4`.NewEncoder。
)
