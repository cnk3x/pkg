package respond

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/cnk3x/gopkg/jsonx"
	"github.com/cnk3x/gopkg/syncx"
)

type ServerEventSource struct {
	// Retry 指定浏览器重新发起连接的时间间隔。
	//
	// 两种情况会导致浏览器重新发起连接
	//  - 时间间隔到期，
	//  - 由于网络错误等原因，导致连接出错。
	Retry int

	// 写入 Access-Control-Allow-Origin 头
	AllowOrigin string

	// 心跳间隔
	Heartbeat int

	// 数据
	Data <-chan any

	w io.Writer
}

func ServerEvent(w http.ResponseWriter, r *http.Request, source ServerEventSource) {
	source.Upgrade(w, r)
}

func (se *ServerEventSource) Upgrade(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/event-stream; charset=utf-8")
	h.Set("Cache-Control", "no-cache")
	if r.ProtoMajor == 1 {
		// An endpoint MUST NOT generate an HTTP/2 message containing connection-specific header fields.
		// Source: RFC7540
		h.Set("Connection", "keep-alive")
	}

	h.Set("Access-Control-Allow-Origin", se.AllowOrigin)
	h.Set("X-Content-Type-Options", "nosniff")
	h.Del("Content-Length")

	w.WriteHeader(http.StatusOK)

	if se.Retry > 0 {
		se.Write(`retry`, []byte(strconv.Itoa(se.Retry)))
		se.Flush()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		heartbeat, stop := syncx.Heartbeat(time.Duration(se.Heartbeat) * time.Second)
		defer stop()

		for ctx := r.Context(); ; {
			select {
			case <-ctx.Done():
				se.WriteEventData(`error`, M{"error": ctx.Err().Error()})
				return
			case data, ok := <-se.Data:
				if !ok {
					se.WriteEvent(`EOF`)
					return
				}
				se.WriteAnyData(data)
			case now := <-heartbeat:
				se.WriteEventData(`heartbeat`, now.Format(time.RFC3339))
			}
		}
	}()

	<-done
}

func (se *ServerEventSource) Flush() {
	se.w.Write([]byte("\n")) // nolint: errcheck
	if f, ok := se.w.(http.Flusher); ok {
		f.Flush()
	}
}

func (se *ServerEventSource) Write(name string, line []byte) {
	se.w.Write([]byte(name)) // nolint: errcheck
	se.w.Write([]byte(": ")) // nolint: errcheck
	se.w.Write(line)         // nolint: errcheck
	se.w.Write([]byte("\n")) // nolint: errcheck
}

func (se *ServerEventSource) WriteEvent(event string) {
	if event != "" {
		se.Write(KeyEvent, []byte(event))
	}
}

func (se *ServerEventSource) WriteID(id string) {
	if id != "" {
		se.Write(KeyID, []byte(id))
	}
}

func (se *ServerEventSource) WriteData(data ...[]byte) {
	for _, line := range data {
		if line = bytes.TrimSpace(line); len(line) > 0 {
			se.Write(KeyData, Escape(line))
		}
	}
}

func (se *ServerEventSource) WriteAnyData(data ...any) {
	// nolint: errcheck
	for _, line := range data {
		se.w.Write([]byte(`data: `)) // nolint: errcheck
		switch t := line.(type) {
		case string:
			se.w.Write([]byte(t))
		case []byte:
			se.w.Write(t)
		case jsonx.Raw:
			se.w.Write([]byte(t))
		default:
			json.NewEncoder(se.w).Encode(t)
		}
		se.w.Write([]byte("\n"))
	}
	se.Flush()
}

func (se *ServerEventSource) WriteFull(id, event string, data ...[]byte) {
	if id != "" {
		se.Write(KeyID, []byte(id))
	}

	if event != "" {
		se.Write(KeyEvent, []byte(event))
	}

	se.WriteData(data...)

	se.Flush()
}

func (se *ServerEventSource) WriteEventData(event string, data any) {
	se.WriteEvent(event)
	se.WriteAnyData(data)
}

const (
	KeyID    = "id"
	KeyEvent = "event"
	KeyData  = "data"
	KeyRetry = "retry"
)
