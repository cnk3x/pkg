package filex

import "os"

// Option 定义用于配置 options 的函数类型
type Option func(*options)

// options 保存文件操作的各种可选配置
type options struct {
	createDirs     bool        // 是否自动创建缺失的目录
	createDirsMode os.FileMode // 创建目录时使用的权限模式
	createMode     os.FileMode // 创建文件时使用的权限模式
	append         bool        // 是否以追加模式写入文件
	overwrite      bool        // 是否覆盖已存在的文件
	readonly       bool        // 是否以只读方式打开文件
}

// applyOptions 将传入的 Option 应用到默认 options 上并返回最终配置
func applyOptions(opts ...Option) *options {
	options := &options{
		createDirs:     true, // 默认自动创建目录
		createDirsMode: 0755, // 默认目录权限：755
		createMode:     0644, // 默认文件权限：644
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// Options 将多个 Option 合并为一个 Option，便于批量应用
func Options(opts ...Option) Option {
	return func(opt *options) {
		for _, apply := range opts {
			apply(opt)
		}
	}
}

// CreateDirs 返回一个 Option，用于设置是否自动创建目录及目录权限
func CreateDirs(createDirs bool, mode os.FileMode) Option {
	return func(opts *options) {
		opts.createDirs = createDirs
		opts.createDirsMode = mode
	}
}

// CreateMode 返回一个 Option，用于设置创建文件时的权限模式
func CreateMode(mode os.FileMode) Option { return func(opts *options) { opts.createMode = mode } }

// Readonly 返回一个 Option，用于设置以只读方式打开文件
func Readonly() Option { return func(opts *options) { opts.readonly = true } }

// Overwrite 返回一个 Option，用于设置覆盖已存在的文件
func Overwrite() Option { return func(opts *options) { opts.overwrite = true } }

// Append 返回一个 Option，用于设置以追加模式写入文件
func Append() Option { return func(opts *options) { opts.append = true } }
