package filex

import (
	"os"
	"syscall"
	"time"
)

func FileInfoAccessTime(fi os.FileInfo) time.Time {
	sec := fi.Sys().(*syscall.Dir).Atime
	return time.Unix(int64(sec), 0)
}
