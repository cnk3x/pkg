package fsw

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cnk3x/pkg/rex"
	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

const throttleMin = time.Second * 2

type Watcher struct {
	root     []string
	exclude  func(string) bool
	allowOp  fsnotify.Op
	throttle time.Duration

	routes  []*Route
	watches []string
	fw      *fsnotify.Watcher
	ctx     context.Context
	mu      sync.Mutex
}

type Route struct {
	Name    string
	Match   func(string) bool
	Op      fsnotify.Op
	Handler HandlerFunc

	timer  *time.Timer
	events []fsnotify.Event

	mu sync.Mutex
}

type (
	HandlerFunc   func(ctx context.Context, events []fsnotify.Event)
	HandlerOption func(*Route)
)

type Options struct {
	Root     []string
	Exclude  []string
	Event    string
	Throttle time.Duration
}

func New(options Options) *Watcher {
	return &Watcher{
		root:     options.Root,
		exclude:  rex.Compile(options.Exclude...),
		allowOp:  Op(options.Event),
		throttle: options.Throttle,
	}
}

func (w *Watcher) Handle(name string, options ...HandlerOption) {
	r := &Route{Name: name}
	for _, opt := range options {
		opt(r)
	}
	r.timer = time.AfterFunc(max(w.throttle, throttleMin), func() { r.Run(w.ctx) })
	r.timer.Stop()
	w.routes = append(w.routes, r)
}

func (w *Watcher) Run(ctx context.Context) (err error) {
	slog.Info("watcher run")
	w.ctx = ctx

	if w.fw, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	for _, f := range w.root {
		w.Add(f)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("watcher context done: %w", context.Cause(ctx))
		case err = <-w.fw.Errors:
			return fmt.Errorf("watcher error: %w", err)
		case ev := <-w.fw.Events:
			if w.exclude != nil && w.exclude(filepath.Base(ev.Name)) {
				continue
			}

			switch {
			case ev.Op.Has(fsnotify.Remove | fsnotify.Rename):
				w.Remove(ev.Name)
			case ev.Op.Has(fsnotify.Create):
				w.Add(ev.Name)
			}

			if w.allowOp != 0 && w.allowOp&ev.Op == 0 {
				slog.Debug(fmt.Sprintf("event skip op %s %s", ev.Op.String(), ev.Name), "allowOp", w.allowOp.String())
				continue
			}

			for _, r := range w.routes {
				if (r.Op == 0 || r.Op&ev.Op != 0) && (r.Match == nil || r.Match(ev.Name)) {
					r.mu.Lock()
					r.events = append(r.events, ev)
					r.timer.Reset(max(w.throttle, throttleMin))
					r.mu.Unlock()
				}
			}
		}
	}
}

// 递归删除
func (w *Watcher) Remove(fullPath string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, f := range w.watches {
		if f == fullPath || strings.HasPrefix(f, fullPath+string(os.PathSeparator)) {
			if err := w.fw.Remove(f); err != nil {
				slog.Error("watch remove fail", "path", f, "err", err)
			}
		}
	}
	w.watches = w.fw.WatchList()
}

// 递归添加
func (w *Watcher) Add(dir string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		return
	}

	if dir, err = filepath.Abs(dir); err != nil {
		slog.Error("watch add fail", "path", dir, "err", err)
		return
	}

	if err = filepath.WalkDir(dir, func(fullPath string, d fs.DirEntry, init error) error {
		if init != nil || !d.IsDir() {
			return init
		}

		if w.exclude != nil && w.exclude(d.Name()) {
			return fs.SkipDir
		}

		if lo.SomeBy(w.watches, func(p string) bool { return p == fullPath }) {
			return nil
		}

		if e := w.fw.Add(fullPath); e != nil {
			slog.Error("watch add fail", "path", fullPath, "err", e)
			return nil
		}

		w.watches = append(w.watches, fullPath)
		return nil
	}); err != nil {
		slog.Error("watch add fail", "path", dir, "err", err)
	}
}

func (r *Route) Run(ctx context.Context) {
	if len(r.events) == 0 || r.Handler == nil {
		return
	}

	r.mu.Lock()
	events := r.events
	r.events = r.events[:0]
	r.mu.Unlock()
	r.Handler(ctx, events)
}

func Match(match ...string) HandlerOption      { return func(r *Route) { r.Match = rex.Compile(match...) } }
func Handle(handler HandlerFunc) HandlerOption { return func(r *Route) { r.Handler = handler } }
func Events(eventOp string) HandlerOption      { return func(r *Route) { r.Op = Op(eventOp) } }

func Op(s string) fsnotify.Op {
	if s == "" {
		return 0
	}

	return lo.Reduce([]rune(s), func(agg fsnotify.Op, item rune, _ int) fsnotify.Op {
		switch item {
		case 'c':
			return agg | fsnotify.Create
		case 'w':
			return agg | fsnotify.Write
		case 'd':
			return agg | fsnotify.Remove
		case 'm':
			return agg | fsnotify.Rename
		default:
			return agg
		}
	}, 0)
}
