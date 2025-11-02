package configx

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cnk3x/gopkg/errx"
)

// FindFile 查找文件，支持多个扩展名。
func FindFile(filePath string, exts ...string) (string, error) {
	existFiles, err := findFile(filePath, exts, true)
	if err != nil {
		return "", err
	}
	return existFiles[0], nil
}

// FindFiles 查找文件，支持多个扩展名。
func FindFiles(filePath string, exts ...string) ([]string, error) {
	return findFile(filePath, exts, false)
}

func findFile(filePath string, exts []string, single bool) (existFiles []string, err error) {
	files := append(make([]string, 0, 1+len(exts)), filePath)

	if len(exts) > 0 {
		srcExt := path.Ext(filePath)
		noExtPath := strings.TrimSuffix(filePath, srcExt)
		for _, ext := range exts {
			if ext == srcExt {
				continue
			}
			files = append(files, noExtPath+ext)
		}
	}

	for _, fn := range files {
		if stat, e := os.Stat(fn); e == nil && stat.Mode().IsRegular() {
			existFiles = append(existFiles, fn)
			if single {
				return
			}
		}
	}
	return
}

// UnmarshalFiles 从多个文件中加载配置数据，并将其反序列化为指定的 Go 值。
// 如果 errContinue 为 true，则在遇到错误时继续加载下一个文件；
// 如果为 false，则在遇到错误时立即返回。
func UnmarshalFiles(value any, files []string, errContinue bool) (err error) {
	for _, fn := range files {
		if e := unmarshalFile(value, fn); e != nil {
			err = errx.Join(err, fmt.Errorf("unmarshal file %s failed: %w", fn, e))
			if errContinue {
				continue
			} else {
				return
			}
		}
	}
	return
}

// UnmarshalFile 从文件中加载配置数据，并将其反序列化为指定的 Go 值。
func UnmarshalFile(value any, fn string, exts ...string) (err error) {
	if len(exts) > 0 {
		fns, e := FindFiles(fn, exts...)
		if err = e; err != nil {
			return
		}
		return UnmarshalFiles(value, fns, true)
	}
	return unmarshalFile(value, fn)
}

// UnmarshalFile 从文件中加载配置数据，并将其反序列化为指定的 Go 值。
func unmarshalFile(value any, fn string) (err error) {
	configData, e := os.ReadFile(fn)
	if err = e; err != nil {
		return
	}

	switch ext := path.Ext(fn); ext {
	case ".json", ".jsonc", "json":
		err = UnmarshalJSONC(configData, value)
	case ".yaml", ".yml", "yaml":
		err = UnmarshalYAML(configData, value)
	default:
		err = fmt.Errorf("unsupported file ext: %s", ext)
	}
	return
}
