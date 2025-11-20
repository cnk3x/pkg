package logx

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func With(prefix string) *Logger {
	return sWith("service", prefix)
}

func Init(level Level, addSource bool, prefix string) {
	SetDefault(New(level, addSource).With("service", prefix))
}

func New(level Level, addSource bool) *Logger {
	tOpt := &TintOptions{
		Level:      level,
		TimeFormat: "15:04:05",
		NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
		AddSource:  addSource,
	}
	return sNew(Prefix(Tint(colorable.NewColorableStdout(), tOpt), &PrefixOptions{PrefixKeys: []string{"service"}}))
}
