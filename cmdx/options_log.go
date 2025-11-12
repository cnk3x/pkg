package cmdx

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type logConfig struct {
	Out string `json:"out,omitempty"`
	Err string `json:"err,omitempty"`

	processInline func(s string, stdout bool)
}

type LogConfig logConfig

func (p LogConfig) MarshalJSON() ([]byte, error) {
	if p.Out == p.Err {
		return json.Marshal(p.Out)
	}
	return json.Marshal(logConfig(p))
}

func (p *LogConfig) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return
	}

	if data[0] == '"' {
		if err = json.Unmarshal(data, &p.Out); err == nil {
			p.Err = p.Out
		}
		return
	}

	m := (*logConfig)(p)
	if err = json.Unmarshal(data, m); err == nil {
		*p = LogConfig(*m)
	}
	return
}

func (p *LogConfig) Open() (stdout *os.File, stderr *os.File, closeIt func(), err error) {
	nOut, nErr := normalizeLog(p.Out, true), normalizeLog(p.Err, false)

	if stdout, err = createLog(nOut, func(line string) {
		if p.processInline != nil {
			p.processInline(line, true)
		} else {
			slog.Info(line)
		}
	}); err != nil {
		return
	}

	if nErr != nOut {
		if stderr, err = createLog(nErr, func(line string) {
			if p.processInline != nil {
				p.processInline(line, false)
			} else {
				slog.Warn(line)
			}
		}); err != nil {
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

func createLog(log string, processInline func(line string)) (*os.File, error) {
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
		go lineRead(r, func(line string) {
			if processInline != nil {
				processInline(line)
			} else {
				slog.Info(line)
			}
		})
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
