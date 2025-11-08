//go:build darwin || freebsd || netbsd
// +build darwin freebsd netbsd

package filex

import (
	"os"
	"syscall"
	"time"
)

func FileInfoAccessTime(fi os.FileInfo) time.Time {
	ts := fi.Sys().(*syscall.Stat_t).Atimespec
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
