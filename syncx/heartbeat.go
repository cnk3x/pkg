package syncx

import (
	"sync"
	"time"
)

// 创建一个心跳计时器，如果 interval<=0，这个计时器永远不会触发，通道会一直阻塞，直到调用 stop
func Heartbeat(interval time.Duration) (ticker <-chan time.Time, stop func()) {
	if interval > 0 {
		t := time.NewTicker(interval)
		return t.C, t.Stop
	}
	t := make(chan time.Time)
	return t, sync.OnceFunc(func() { close(t) })
}
