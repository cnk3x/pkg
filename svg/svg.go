package svg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

func Sprite(dstPath string, files []string, opts ...Option) error {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	slog.Debug("svg sprite generate", "target", dstPath, "pretty", options.Pretty)

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
		options.ID = NameFunc("")
	}

	var spaceRe = regexp.MustCompile(`\s+`)
	var nodeSpaceRe = regexp.MustCompile(`\s*(<|>)\s*`)
	var symbols []Symbol

	for _, file := range files {
		// 获取图标 ID
		iconId := options.ID(file)

		// 去重
		for constId, i := iconId, 1; ; i++ {
			if !lo.SomeBy(symbols, func(s Symbol) bool { return s.ID == iconId }) {
				break
			}
			n := fmt.Sprintf("%s-%d", constId, i)
			slog.Warn(fmt.Sprintf("[svg sprite] id %q replicated, rename to %q", iconId, n))
			iconId = n
		}

		// 读取 SVG 内容
		svgData, err := os.ReadFile(file)
		if err != nil {
			slog.Warn(fmt.Sprintf("[svg sprite] read %s fail: %v, skip", file, err))
			continue
		}

		// 解析 SVG，提取 viewBox 和内部内容
		var svg SVG
		if err := xml.Unmarshal(svgData, &svg); err != nil {
			slog.Warn(fmt.Sprintf("[svg sprite] unmarshal %s fail: %v, skip", file, err))
			continue
		}

		// 必须有 viewBox（否则复用后尺寸异常）
		if svg.ViewBox == "" {
			slog.Warn(fmt.Sprintf("[svg sprite] %s miss viewBox, skip", file))
			continue
		}

		svg.Content = spaceRe.ReplaceAll(svg.Content, []byte(" "))
		svg.Content = nodeSpaceRe.ReplaceAll(svg.Content, []byte(`$1`))

		if len(svg.Content) == 0 {
			slog.Warn(fmt.Sprintf("[svg sprite] %s miss content, skip", file))
			continue
		}

		// 构建 symbol（移除 svg 标签，保留内部内容）
		symbols = append(symbols, Symbol{ID: iconId, ViewBox: svg.ViewBox, Content: svg.Content})
	}

	if len(symbols) == 0 {
		return fmt.Errorf("svg sprite no valid symbols find")
	}

	slog.Debug("[svg sprite] find symbols", "count", len(symbols))

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
		return fmt.Errorf("svg sprite write header fail: %w", err)
	}

	// 序列化 XML
	encoder := xml.NewEncoder(&buf)
	if options.Pretty {
		encoder.Indent("", "  ")
	}
	if err := encoder.Encode(sprite); err != nil {
		return fmt.Errorf("svg sprite xml encode fail: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("svg sprite xml encode close fail: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(dstPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("svg sprite write file fail: %w", err)
	}

	names := lo.Map(symbols, func(item Symbol, _ int) string { return item.ID })

	symbolsFile := dstPath + ".symbols.json"
	if err := os.WriteFile(symbolsFile, fmt.Appendf(nil, `["%s"]`, strings.Join(names, `", "`)), 0644); err != nil {
		slog.Warn(fmt.Sprintf("[svg sprite] %s write symbols name fail: %v", symbolsFile, err))
	}

	slog.Debug("[svg sprite] generated completed", "file", dstPath, "count", len(symbols), "symbols", strings.Join(names, ", "))
	return nil
}
