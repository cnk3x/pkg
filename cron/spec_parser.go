package cron

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// ParseOption 创建解析器的配置选项
// 大多数选项指定应包括哪些字段，而其他选项启用功能
// 如果未包含某个字段，解析器将假定一个默认值
// 这些选项不会更改解析字段的顺序
type ParseOption int

const (
	Second         ParseOption = 1 << iota // 秒字段，默认值为0
	SecondOptional                         // 可选秒字段，默认值为0
	Minute                                 // 分字段，默认值为0
	Hour                                   // 时字段，默认值为0
	Dom                                    // 月中的天字段，默认值为*
	Month                                  // 月字段，默认值为*
	Dow                                    // 周中的天字段，默认值为*
	DowOptional                            // 可选周中的天字段，默认值为*
	Descriptor                             // 允许使用描述符，如@monthly、@weekly等
)

var places = []ParseOption{
	Second,
	Minute,
	Hour,
	Dom,
	Month,
	Dow,
}

var defaults = []string{
	"0",
	"0",
	"0",
	"*",
	"*",
	"*",
}

// Parser 可配置的自定义解析器
type Parser struct {
	options ParseOption
}

// New 使用自定义选项创建解析器
//
// 如果配置了多个可选项，它会panic，因为在一般情况下无法正确推断提供了哪个可选项或缺少哪个可选项
//
// 示例
//
//	// 不带描述符的标准解析器
//	specParser := New(Minute | Hour | Dom | Month | Dow)
//	sched, err := specParser.Parse("0 0 15 */3 *")
//
//	// 与上述相同，只是排除了时间字段
//	specParser := New(Dom | Month | Dow)
//	sched, err := specParser.Parse("15 */3 *")
//
//	// 与上述相同，只是使Dow成为可选的
//	specParser := New(Dom | Month | DowOptional)
//	sched, err := specParser.Parse("15 */3")
//
// options: 解析器选项
//
// 返回值: 配置好的Parser实例
func New(options ParseOption) Parser {
	optionals := 0
	if options&DowOptional > 0 {
		optionals++
	}
	if options&SecondOptional > 0 {
		optionals++
	}
	if optionals > 1 {
		panic("multiple optionals may not be configured")
	}
	return Parser{options}
}

// Parse 返回表示给定规范的新crontab计划
// 如果规范无效，它将返回描述性错误
// 它接受crontab规范和由NewParser配置的功能
//
//   - spec: 计划规范字符串
//
// 返回值: 解析后的Schedule接口和可能的错误
func (p Parser) Parse(spec string) (Schedule, error) {
	if len(spec) == 0 {
		return nil, fmt.Errorf("empty spec string")
	}

	// 如果存在，提取时区
	var loc = time.Local
	if strings.HasPrefix(spec, "TZ=") || strings.HasPrefix(spec, "CRON_TZ=") {
		var err error
		i := strings.Index(spec, " ")
		eq := strings.Index(spec, "=")
		if loc, err = time.LoadLocation(spec[eq+1 : i]); err != nil {
			return nil, fmt.Errorf("provided bad location %s: %v", spec[eq+1:i], err)
		}
		spec = strings.TrimSpace(spec[i:])
	}

	// 如果配置了，处理命名计划（描述符）
	if strings.HasPrefix(spec, "@") {
		if p.options&Descriptor == 0 {
			return nil, fmt.Errorf("parser does not accept descriptors: %v", spec)
		}
		return parseDescriptor(spec, loc)
	}

	// 按空格分割
	fields := strings.Fields(spec)

	// 验证并填充任何省略或可选的字段
	var err error
	fields, err = normalizeFields(fields, p.options)
	if err != nil {
		return nil, err
	}

	field := func(field string, r bounds) uint64 {
		if err != nil {
			return 0
		}
		var bits uint64
		bits, err = getField(field, r)
		return bits
	}

	var (
		second     = field(fields[0], seconds)
		minute     = field(fields[1], minutes)
		hour       = field(fields[2], hours)
		dayOfMonth = field(fields[3], dom)
		month      = field(fields[4], months)
		dayOfWeek  = field(fields[5], dow)
	)
	if err != nil {
		return nil, err
	}

	return &SpecSchedule{
		Second:   second,
		Minute:   minute,
		Hour:     hour,
		Dom:      dayOfMonth,
		Month:    month,
		Dow:      dayOfWeek,
		Location: loc,
	}, nil
}

// normalizeFields 接受时间字段的一个子集，并返回填充了默认值（零）的完整字段集
// 作为执行此函数的一部分，它还会验证提供的字段是否与配置的选项兼容
//
//   - fields: 字段数组
//   - options: 解析选项
//
// 返回值: 标准化后的字段数组和可能的错误
func normalizeFields(fields []string, options ParseOption) ([]string, error) {
	// 验证可选项并将它们的字段添加到选项中
	optionals := 0
	if options&SecondOptional > 0 {
		options |= Second
		optionals++
	}
	if options&DowOptional > 0 {
		options |= Dow
		optionals++
	}
	if optionals > 1 {
		return nil, fmt.Errorf("multiple optionals may not be configured")
	}

	// 计算我们需要多少个字段
	max := 0
	for _, place := range places {
		if options&place > 0 {
			max++
		}
	}
	min := max - optionals

	// 验证字段数量
	if count := len(fields); count < min || count > max {
		if min == max {
			return nil, fmt.Errorf("expected exactly %d fields, found %d: %s", min, count, fields)
		}
		return nil, fmt.Errorf("expected %d to %d fields, found %d: %s", min, max, count, fields)
	}

	// 如果未提供可选字段，则填充它
	if min < max && len(fields) == min {
		switch {
		case options&DowOptional > 0:
			fields = append(fields, defaults[5]) // TODO: improve access to default
		case options&SecondOptional > 0:
			fields = append([]string{defaults[0]}, fields...)
		default:
			return nil, fmt.Errorf("unknown optional field")
		}
	}

	// 使用默认值填充不属于选项的所有字段
	n := 0
	expandedFields := make([]string, len(places))
	copy(expandedFields, defaults)
	for i, place := range places {
		if options&place > 0 {
			expandedFields[i] = fields[n]
			n++
		}
	}
	return expandedFields, nil
}

var Standard = New(Minute | Hour | Dom | Month | Dow | Descriptor)

// ParseStandard 返回表示给定标准规范的新crontab计划
// 它需要5个条目，分别代表：分、时、月中天、月和周中天，按此顺序排列
// 如果规范无效，它将返回描述性错误
//
// 它接受
//   - 标准crontab规范，例如 "* * * * ?"
//   - 描述符，例如 "@midnight"、"@every 1h30m"
//
// standardSpec: 标准计划规范字符串
//
// 返回值: 解析后的Schedule接口和可能的错误
func ParseStandard(standardSpec string) (Schedule, error) {
	return Standard.Parse(standardSpec)
}

// getField 返回一个设置了代表字段所表示的所有时间的位的整数，或解析字段值时出错
// "字段"是以逗号分隔的"范围"列表
//
// field: 字段字符串
// r: 边界值
//
// 返回值: 位集和可能的错误
func getField(field string, r bounds) (uint64, error) {
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })
	for _, expr := range ranges {
		bit, err := getRange(expr, r)
		if err != nil {
			return bits, err
		}
		bits |= bit
	}
	return bits, nil
}

// getRange 返回由给定表达式指示的位：
//
//	number | number "-" number [ "/" number ]
//
// 或解析范围时出错
//
//   - expr: 表达式字符串
//   - r: 边界值
//
// 返回值: 位集和可能的错误
func getRange(expr string, r bounds) (uint64, error) {
	var (
		start, end, step uint
		rangeAndStep     = strings.Split(expr, "/")
		lowAndHigh       = strings.Split(rangeAndStep[0], "-")
		singleDigit      = len(lowAndHigh) == 1
		err              error
	)

	var extra uint64
	if lowAndHigh[0] == "*" || lowAndHigh[0] == "?" {
		start = r.min
		end = r.max
		extra = starBit
	} else {
		start, err = parseIntOrName(lowAndHigh[0], r.names)
		if err != nil {
			return 0, err
		}
		switch len(lowAndHigh) {
		case 1:
			end = start
		case 2:
			end, err = parseIntOrName(lowAndHigh[1], r.names)
			if err != nil {
				return 0, err
			}
		default:
			return 0, fmt.Errorf("too many hyphens: %s", expr)
		}
	}

	switch len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		step, err = mustParseInt(rangeAndStep[1])
		if err != nil {
			return 0, err
		}

		// 特殊处理："N/step"表示"N-max/step"。
		if singleDigit {
			end = r.max
		}
		if step > 1 {
			extra = 0
		}
	default:
		return 0, fmt.Errorf("too many slashes: %s", expr)
	}

	if start < r.min {
		return 0, fmt.Errorf("beginning of range (%d) below minimum (%d): %s", start, r.min, expr)
	}
	if end > r.max {
		return 0, fmt.Errorf("end of range (%d) above maximum (%d): %s", end, r.max, expr)
	}
	if start > end {
		return 0, fmt.Errorf("beginning of range (%d) beyond end of range (%d): %s", start, end, expr)
	}
	if step == 0 {
		return 0, fmt.Errorf("step of range should be a positive number: %s", expr)
	}

	return getBits(start, end, step) | extra, nil
}

// parseIntOrName 返回expr中包含的（可能命名的）整数
//
//   - expr: 表达式字符串
//   - names: 名称映射
//
// 返回值: 解析出的整数和可能的错误
func parseIntOrName(expr string, names map[string]uint) (uint, error) {
	if names != nil {
		if namedInt, ok := names[strings.ToLower(expr)]; ok {
			return namedInt, nil
		}
	}
	return mustParseInt(expr)
}

// mustParseInt 将给定表达式解析为int或返回错误
//
//   - expr: 表达式字符串
//
// 返回值: 解析出的整数和可能的错误
func mustParseInt(expr string) (uint, error) {
	num, err := strconv.Atoi(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int from %s: %s", expr, err)
	}
	if num < 0 {
		return 0, fmt.Errorf("negative number (%d) not allowed: %s", num, expr)
	}

	return uint(num), nil
}

// getBits 设置范围内[min, max]的所有位，模给定的步长
//
//   - min: 起始值
//   - max: 结束值
//   - step: 步长
//
// 返回值: 位集
func getBits(min, max, step uint) uint64 {
	var bits uint64

	// 如果步长为1，使用移位。
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	// 否则，使用简单循环。
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

// all 返回给定边界内的所有位（加上星号位）
//
//   - r: 边界值
//
// 返回值: 位集
func all(r bounds) uint64 {
	return getBits(r.min, r.max, 1) | starBit
}

// parseDescriptor 返回表达式的预定义计划，如果没有匹配项则返回错误
//
//   - descriptor: 描述符字符串
//   - loc: 时区位置
//
// 返回值: 解析后的Schedule接口和可能的错误
func parseDescriptor(descriptor string, loc *time.Location) (Schedule, error) {
	switch descriptor {
	case "@yearly", "@annually":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      1 << dom.min,
			Month:    1 << months.min,
			Dow:      all(dow),
			Location: loc,
		}, nil

	case "@monthly":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      1 << dom.min,
			Month:    all(months),
			Dow:      all(dow),
			Location: loc,
		}, nil

	case "@weekly":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      all(dom),
			Month:    all(months),
			Dow:      1 << dow.min,
			Location: loc,
		}, nil

	case "@daily", "@midnight":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      all(dom),
			Month:    all(months),
			Dow:      all(dow),
			Location: loc,
		}, nil

	case "@hourly":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     all(hours),
			Dom:      all(dom),
			Month:    all(months),
			Dow:      all(dow),
			Location: loc,
		}, nil

	}

	const every = "@every "
	if strings.HasPrefix(descriptor, every) {
		duration, err := time.ParseDuration(descriptor[len(every):])
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration %s: %s", descriptor, err)
		}
		return Every(duration), nil
	}

	return nil, fmt.Errorf("unrecognized descriptor: %s", descriptor)
}
