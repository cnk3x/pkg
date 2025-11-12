package configx

import (
	"runtime"

	"github.com/cnk3x/gopkg/strx"
)

const (
	IsWindows = runtime.GOOS == "windows"
	GOOS      = runtime.GOOS
	GOARCH    = runtime.GOARCH
	OSARCH    = GOOS + "-" + GOARCH
)

var (
	BinaryExt  = iif(IsWindows, ".exe", "")
	ArchiveExt = iif(IsWindows, ".zip", ".tar.gz")
)

var SysEnvs = map[string]any{
	"BASE":        workPath,
	"OS":          GOOS,
	"ARCH":        GOARCH,
	"OSARCH":      OSARCH,
	"ARCHIVE_EXT": ArchiveExt,
	"BIN_EXT":     BinaryExt,
}

func iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}

func ReplaceEnv(src string, args ...any) string {
	return strx.Replace(src, strx.TagUpper(strx.TagMap(SysEnvs)), strx.TagArg(args...))
}
