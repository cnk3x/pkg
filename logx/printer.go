package logx

import (
	"context"
	"fmt"
)

func Printer(prefix string, level Level) interface{ Print(v ...any) } {
	return &printer{log: With(prefix), level: level}
}

type printer struct {
	log   *Logger
	level Level
}

func (r *printer) Print(v ...any) { r.log.Log(context.Background(), r.level, fmt.Sprint(v...)) }
