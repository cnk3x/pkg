package fsw

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

func TestWatcher(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fw := New(Options{
		Root:    []string{"../"},
		Exclude: []string{`^[\._-]`, "modules"},
	})

	//  Match(`!(.*)_test\.go$`),
	fw.Handle("log", Handle(func(ctx context.Context, ev []fsnotify.Event) {
		t.Log("handle log", strings.Join(lo.Map(ev, func(e fsnotify.Event, _ int) string { return fmt.Sprintf("\n\t-- [%s] %s", e.Op.String(), e.Name) }), ""))
	}))

	if err := fw.Run(ctx); err != nil {
		t.Fatal(err)
	}
}
