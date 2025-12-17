package archive

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/cnk3x/pkg/filex"
)

type options struct {
	stripComponents int
	skipEmptyDir    bool
	progress        func(index int, name string, cur, total int64)
	filters         []string
}

const pathSeparator = string(filepath.Separator)

// Extract 创建一个解压处理器，用于将归档项提取到指定目录
//
// 参数:
//   - dir: 目标目录路径
//   - extractOptions: 解压选项列表
//
// 返回值:
//   - ProcessFunc: 处理函数，用于处理每个归档项
func Extract(dir string, extractOptions ...Option) ProcessFunc {
	var eop options
	for _, o := range extractOptions {
		o(&eop)
	}

	// 返回实际的处理函数
	return func(ctx context.Context, item Item) error {
		// 清理并获取文件路径
		fpath := filepath.Clean(item.Path())

		// 处理 stripComponents 选项，移除路径前缀
		if eop.stripComponents > 0 {
			paths := strings.Split(fpath, pathSeparator)
			if len(paths) <= eop.stripComponents {
				return nil
			}
			fpath = filepath.Join(paths[eop.stripComponents:]...)
		}

		// 跳过空路径
		if fpath == "" {
			return nil
		}

		// 应用过滤器筛选文件
		if len(eop.filters) > 0 {
			for _, f := range eop.filters {
				re, e := regexp.Compile(f)
				if e != nil {
					return e
				}
				if !re.MatchString(fpath) {
					return nil
				}
			}
		}

		// 构建目标文件路径
		target := filepath.Join(dir, fpath)

		// 处理目录创建
		if item.IsDir() && !eop.skipEmptyDir {
			err := os.MkdirAll(target, item.Mode())
			return err
		}

		// 打开归档项
		it, err := item.Open()
		if err != nil {
			return err
		}

		// 设置进度回调函数
		var p filex.ProgressFunc
		if eop.progress != nil {
			current, total, index := int64(0), item.Size(), item.Index()
			p = func(n int64) { eop.progress(index, fpath, atomic.AddInt64(&current, n), total) }
		}

		// 写入文件到目标位置
		return filex.OpenWrite(target, filex.WriteFrom(ctx, it, p), filex.CreateMode(item.Mode().Perm()))
	}
}

// Option 定义了解压选项的函数类型
type Option func(option *options)

// Filter 设置文件路径过滤器
//
// 参数:
//   - filters: 正则表达式字符串列表，用于匹配需要提取的文件路径
//
// 返回值:
//   - Option: 选项函数
func Filter(filters ...string) Option { return func(option *options) { option.filters = filters } }

// StripComponents 设置要剥离的路径层级数
//
// 参数:
//   - deep: 要从路径开头剥离的目录层数
//
// 返回值:
//   - Option: 选项函数
func StripComponents(deep int) Option {
	return func(option *options) { option.stripComponents = deep }
}

// SkipEmptyDir 设置是否跳过空目录
//
// 参数:
//   - skip: 是否跳过空目录的布尔值
//
// 返回值:
//   - Option: 选项函数
func SkipEmptyDir(skip bool) Option {
	return func(option *options) { option.skipEmptyDir = skip }
}

// Progress 设置进度回调函数
//
// 参数:
//   - progress: 进度回调函数，参数分别为索引、文件名、当前进度、总大小
//
// 返回值:
//   - Option: 选项函数
func Progress(progress func(index int, name string, cur, total int64)) Option {
	return func(option *options) { option.progress = progress }
}
