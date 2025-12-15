package fsw

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

func TestWatcher(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	w := New(Root("../"), Filter(`!/\.(.*)`, "!modules"))

	w.Handle("log", Match(`!(.*)_test\.go$`), Handle(func(ctx context.Context, ev []fsnotify.Event) {
		t.Log("handle", lo.Uniq(lo.Map(ev, func(e fsnotify.Event, _ int) string { return e.Name })))
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := w.Run(ctx); err != nil {
		t.Fatal(err)
	}
}
