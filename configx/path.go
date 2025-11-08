package configx

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	workPath string
	envKey   string
)

func init() {
	pathInit()
}

func pathInit() {
	workPath = getEnv("WORK_PATH", ".")
	workPath, _ = filepath.Abs(workPath)
}

func SetEnvKey(key string) {
	envKey = key
	pathInit()
}

func getEnv(key string, defaultValue string) (value string) {
	if envKey != "" {
		key = envKey + "_" + key
	}
	if v, ok := os.LookupEnv(strings.ToUpper(key)); ok {
		value = v
	} else {
		value = defaultValue
	}
	return
}

func SetWorkPath(path string) { workPath = path }
func WorkPath() string        { return workPath }

func GetPath(paths ...string) string {
	paths = append(paths, workPath)
	for i, j := 0, 0; i < len(paths)-1; i++ {
		if p := strings.TrimSpace(paths[i]); p != "" {
			paths[j] = p
			j++
		}
	}
	paths = slices.DeleteFunc(slices.Insert(paths, 0, workPath), func(s string) bool { return s == "" })
	return filepath.Join(paths...)
}
