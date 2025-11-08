package cmdx

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type LogOptions logOptions

type logOptions struct {
	Out string `json:"out,omitempty"`
	Err string `json:"err,omitempty"`
}

func (p LogOptions) MarshalJSON() ([]byte, error) {
	if p.Out == p.Err {
		return json.Marshal(p.Out)
	}
	return json.Marshal(logOptions(p))
}

func (p *LogOptions) UnmarshalJSON(data []byte) (err error) {
	data = bytes.TrimSpace(data)
	if len(data) < 2 {
		return
	}

	if data[0] == '{' && data[len(data)-1] == '}' {
		m := (logOptions)(*p)
		if err = json.Unmarshal(data, &m); err == nil {
			*p = LogOptions(m)
		}
		return
	}

	var s string
	if err = json.Unmarshal(data, &s); err == nil {
		p.Out, p.Err = s, s
	}
	return
}

func (p *LogOptions) Open() (stdout *os.File, stderr *os.File, closeIt func(), err error) {
	nOut, nErr := normalizeLog(p.Out, true), normalizeLog(p.Err, false)

	if stdout, err = createLog(nOut); err != nil {
		return
	}

	if nErr != nOut {
		if stderr, err = createLog(nErr); err != nil {
			if stdout != nil {
				stdout.Close()
			}
			return
		}
	} else {
		stderr = stdout
	}

	closeIt = func() {
		if stdout != nil {
			stdout.Close()
		}
		if stderr != nil {
			stderr.Close()
		}
	}

	return
}

func createLog(log string) (*os.File, error) {
	switch log {
	case "nul":
		return nil, nil
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	case "inline":
		r, w, err := os.Pipe()
		if err != nil {
			return nil, fmt.Errorf("create inline log error: %w", err)
		}
		go lineRead(r, func(line string) { slog.Info(line) })
		return w, nil
	default:
		if err := os.MkdirAll(filepath.Dir(log), 0755); err != nil {
			return nil, fmt.Errorf("create log dir error: %w", err)
		}
		return os.OpenFile(log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	}
}

func normalizeLog(log string, isOut bool) string {
	switch log {
	case "", "NUL", "null", "devnull", "discard":
		return "nul"
	case "std":
		return iif(isOut, "stdout", "stderr")
	default:
		return log
	}
}

func iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}

func lineRead(r io.Reader, lineFunc func(line string)) {
	for s := bufio.NewScanner(r); s.Scan(); {
		if line := s.Text(); line != "" {
			lineFunc(line)
		}
	}
}
