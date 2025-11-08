package configx

import (
	"log/slog"
	"os"
)

func Load[T any](filePath string, exts ...string) func() (T, error) {
	return func() (v T, err error) {
		err = UnmarshalFile(&v, filePath, exts...)
		if os.IsNotExist(err) {
			slog.Warn("config file not exist", "filePath", filePath)
			err = nil
		}
		return
	}
}
