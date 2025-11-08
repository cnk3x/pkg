//go:build linux || dragonfly || openbsd || solaris
// +build linux dragonfly openbsd solaris

package filex

import (
	"os"
	"syscall"
	"time"
)

func FileInfoAccessTime(fi os.FileInfo) time.Time {
	ts := fi.Sys().(*syscall.Stat_t).Atim
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
