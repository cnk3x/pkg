package svg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

func Sprite(dstPath string, files []string, opts ...Option) error {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	// SVG 基础结构定义
	type SVG struct {
		XMLName xml.Name `xml:"svg"`
		XMLNS   string   `xml:"xmlns,attr"`
		ViewBox string   `xml:"viewBox,attr,omitempty"`
		Width   string   `xml:"width,attr,omitempty"`
		Height  string   `xml:"height,attr,omitempty"`
		Content []byte   `xml:",innerxml"` // 保留 SVG 内部内容
	}

	// Symbol 表示雪碧图中的单个图标符号
	type Symbol struct {
		XMLName xml.Name `xml:"symbol"`
		ID      string   `xml:"id,attr"`
		ViewBox string   `xml:"viewBox,attr"`
		Content []byte   `xml:",innerxml"`
	}

	// Sprite 雪碧图整体结构
	type Sprite struct {
		XMLName xml.Name `xml:"svg"`
		XMLNS   string   `xml:"xmlns,attr"`
		Width   string   `xml:"width,attr"`
		Height  string   `xml:"height,attr"`
		Symbols []Symbol `xml:"symbol"`
	}

	if options.ID == nil {
		options.ID = defaultIdGen
	}

	var symbols []Symbol

	for _, file := range files {
		// 获取图标 ID
		iconId := options.ID(file)

		// 去重
		if lo.SomeBy(symbols, func(s Symbol) bool { return s.ID == iconId }) {
			slog.Warn(fmt.Sprintf("svg sprite id %s replicated, skip", iconId))
			continue
		}

		// 读取 SVG 内容
		svgData, err := os.ReadFile(file)
		if err != nil {
			slog.Error(fmt.Sprintf("读取 %s 失败: %v", file, err))
			continue
		}

		// 解析 SVG，提取 viewBox 和内部内容
		var svg SVG
		if err := xml.Unmarshal(svgData, &svg); err != nil {
			slog.Error(fmt.Sprintf("解析 %s 失败: %v", file, err))
			continue
		}

		// 必须有 viewBox（否则复用后尺寸异常）
		if svg.ViewBox == "" {
			slog.Error(fmt.Sprintf("%s 缺少 viewBox 属性", file))
			continue
		}

		// 构建 symbol（移除 svg 标签，保留内部内容）
		symbols = append(symbols, Symbol{ID: iconId, ViewBox: svg.ViewBox, Content: spaceRe.ReplaceAll(svg.Content, []byte(" "))})
	}

	slog.Debug("symbols", "count", len(symbols))

	// 构建雪碧图 SVG
	sprite := Sprite{
		XMLNS:   "http://www.w3.org/2000/svg",
		Width:   "0", // 隐藏容器
		Height:  "0",
		Symbols: symbols,
	}

	// 补充 XML 声明
	var buf bytes.Buffer
	if _, err := buf.WriteString(xml.Header); err != nil {
		return fmt.Errorf("写入雪碧图文件头失败: %w", err)
	}

	// 序列化 XML
	encoder := xml.NewEncoder(&buf)
	if options.Pretty {
		encoder.Indent("", "  ")
	}
	if err := encoder.Encode(sprite); err != nil {
		return fmt.Errorf("序列化雪碧图失败: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("序列化雪碧图失败: %w", err)
	}

	// 写入文件
	return os.WriteFile(dstPath, buf.Bytes(), 0644)
}

func defaultIdGen(file string) string {
	return cleanRe.ReplaceAllString(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)), "-")
}
