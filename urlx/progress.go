package urlx

import (
	"io"
	"net/http"
	"time"
)

type ProgressState struct {
	Total   float64 // 总字节数
	Current float64 // 当前已读/写取字节数
	Speed   float64 // 当前读/写速度（字节/秒）
}

type ProgressReport func(state ProgressState)

// Progress 下载进度
func Progress(report ProgressReport) Process {
	return func(resp *http.Response) error {
		body := resp.Body
		resp.Body = io.NopCloser(ProgressReader(body, float64(resp.ContentLength), report))
		return nil
	}
}

func ProgressReader(r io.Reader, total float64, report ProgressReport) io.Reader {
	return fReader(ProgressStream(r.Read, total, report))
}

func ProgressWriter(w io.Writer, total float64, report ProgressReport) io.Writer {
	return fWriter(ProgressStream(w.Write, total, report))
}

func ProgressStream(process func([]byte) (int, error), total float64, report ProgressReport) func([]byte) (int, error) {
	var (
		cur       float64   // 当前已读/写取字节数
		rpt_bytes float64   // 上一次报告的字节数
		rpt_time  time.Time // 上一次报告的时间
	)

	return func(p []byte) (n int, err error) {
		if n, err = process(p); err != nil {
			return
		}
		cur += float64(n)

		now := time.Now()
		if rpt_time.IsZero() {
			rpt_time = now
			return
		}

		if d := now.Sub(rpt_time); d >= time.Second || cur >= total {
			speed := (cur - rpt_bytes) / d.Seconds()
			rpt_bytes, rpt_time = cur, now
			report(ProgressState{Total: total, Current: cur, Speed: speed})
		}

		return
	}
}

type fReader func(p []byte) (n int, err error)

func (f fReader) Read(p []byte) (n int, err error) { return f(p) }

type fWriter func(p []byte) (n int, err error)

func (f fWriter) Write(p []byte) (n int, err error) { return f(p) }
