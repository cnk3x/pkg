package urlx

import (
	"net/http"
	"time"
)

// ProgressReport 进度报告方法
type ProgressReport = func(total float64, cur float64, speed float64)

// Progress 下载进度, reportInterval 报告的时间间隔, 最小1秒, 默认2秒
func Progress(report ProgressReport, reportInterval ...time.Duration) ProcessMw {
	var interval time.Duration

	for _, ri := range reportInterval {
		if ri > time.Second {
			interval = ri
			break
		}
	}

	if interval < time.Second {
		interval = time.Second * 2
	}

	return func(next Process) Process {
		return func(resp *http.Response) error {
			var (
				body  = resp.Body
				total = float64(resp.ContentLength)
				cur   = float64(0) //当前已读取

				cTime  time.Time //上一次计算速度的时间
				cBytes float64   //从上一次计算速度的时间之后新读取的字节数
			)

			calc := func(n ...time.Time) {
				var now time.Time
				if len(n) > 0 {
					now = n[0]
				} else {
					now = time.Now()
				}

				if d := now.Sub(cTime); len(n) == 0 || d >= interval {
					report(total, cur, cBytes/d.Seconds())
					cTime = now
					cBytes = 0
				}
			}

			reader := func(p []byte) (n int, err error) {
				if cTime.IsZero() {
					cTime = time.Now()
				}

				if n, err = body.Read(p); err != nil {
					return
				}

				fv := float64(n)
				cur += fv
				cBytes += fv

				if report != nil {
					calc(time.Now())
				}
				return
			}

			defer calc()
			resp.Body = rFunc(reader)
			return next(resp)
		}
	}
}

type rFunc func(p []byte) (n int, err error)

func (f rFunc) Read(p []byte) (n int, err error) { return f(p) }
func (rFunc) Close() error                       { return nil }
