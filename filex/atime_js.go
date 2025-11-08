package filex

import (
	"os"
	"syscall"
	"time"
)

func FileInfoAccessTime(fi os.FileInfo) time.Time {
	sys := fi.Sys().(*syscall.Stat_t)
	return time.Unix(sys.Atime, sys.AtimeNsec)
}
