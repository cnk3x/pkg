package logx

import (
	"context"
	"fmt"
	"strings"
)

func Printer(prefix string, level Level) iPrinter {
	return &printer{log: With(prefix), level: level}
}

type iPrinter = interface {
	Print(...any)
	Printf(string, ...any)
}

type printer struct {
	log   *Logger
	level Level
}

func (r *printer) Print(v ...any) { r.log.Log(bg, r.level, fmt.Sprint(v...)) }
func (r *printer) Printf(s string, v ...any) {
	r.log.Log(bg, r.level, strings.ReplaceAll(strings.TrimSuffix(fmt.Sprintf(s, v...), "\n"), "\t\t", "\t"))
}

var bg = context.Background()
