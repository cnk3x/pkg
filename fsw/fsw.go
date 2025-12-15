package fsw

import (
	"cmp"
	"context"
	"errors"
	"fmt"
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

type Watcher struct {
	root     []string
	filter   func(string) bool
	allowOp  fsnotify.Op
	throttle time.Duration

	routes  []*Route
	watches []string
	fw      *fsnotify.Watcher
	ctx     context.Context
	mu      sync.Mutex
}

func New(options ...Option) *Watcher {
	w := &Watcher{}
	for _, opt := range options {
		opt(w)
	}
	return w
}

type Options struct {
	Root     []string
	Filter   []string
	Event    string
	Throttle time.Duration
}

func WithOptions(options Options) *Watcher {
	return New(
		Root(options.Root...),
		Filter(options.Filter...),
		AllowOp(options.Event),
		Throttle(options.Throttle),
	)
}

func (w *Watcher) Handle(name string, options ...HandlerOption) {
	r := &Route{Name: name}
	for _, opt := range options {
		opt(r)
	}
	r.timer = time.AfterFunc(max(w.throttle, time.Second*3), func() { r.Run(w.ctx) })
	r.timer.Stop()
	w.routes = append(w.routes, r)
}

func (w *Watcher) Run(ctx context.Context) (err error) {
	slog.Info("watcher run", "err", err)
	defer slog.Info("watcher done", "err", err)

	w.ctx = ctx
	if w.fw, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	for _, root := range w.root {
		if err = w.addRecursive(root); err != nil {
			return
		}
	}

	for _, f := range w.watches {
		slog.Debug("watches", "path", f)
	}

	err = func() (err error) {
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("事件处理通道关闭: %w", context.Cause(ctx))
			case e, ok := <-w.fw.Errors:
				if !ok {
					return
				}
				return e
			case ev, ok := <-w.fw.Events:
				if !ok {
					return
				}

				if !w.filter(ev.Name) {
					continue
				}

				switch {
				case ev.Op.Has(fsnotify.Remove | fsnotify.Rename):
					if e := w.Remove(ev.Name); e != nil {
						slog.Debug("delRecursive", "err", e, "event", ev.Op.String())
					}
				case ev.Op.Has(fsnotify.Create):
					if stat, e := os.Stat(ev.Name); e == nil && stat.IsDir() {
						if e := w.Add(ev.Name); e != nil {
							slog.Debug("addRecursive", "err", e, "event", ev.Op.String())
						}
					}
				}

				slog.Debug(fmt.Sprintf("handleEvent %s %s", ev.Op.String(), ev.Name))
				w.handleEvent(ev)
			}
		}
	}()

	return
}

func (w *Watcher) handleEvent(ev fsnotify.Event) {
	if w.allowOp&ev.Op == 0 {
		slog.Debug(fmt.Sprintf("skipEvent %s %s", ev.Op.String(), ev.Name), "allowOp", w.allowOp.String())
		return
	}

	for _, r := range w.routes {
		if r.Op != 0 && r.Op&ev.Op == 0 {
			continue
		}
		if r.Match == nil || r.Match(ev.Name) {
			slog.Debug("route match", "name", r.Name, "event", ev.Op.String(), "path", ev.Name)
			r.mu.Lock()
			r.events = append(r.events, ev)
			r.timer.Reset(max(w.throttle, time.Second*3))
			r.mu.Unlock()
		}
	}
}

// 递归删除
func (w *Watcher) Remove(fullPath string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, f := range w.watches {
		if strings.HasPrefix(fullPath, f) {
			err = errors.Join(err, w.fw.Remove(f))
		}
	}

	w.watches = w.fw.WatchList()
	return
}

// 递归添加
func (w *Watcher) Add(fullPath string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if stat, e := os.Stat(fullPath); e != nil || stat.IsDir() {
		return cmp.Or(e, fmt.Errorf("路径不是目录，只能添加目录为监控目标： %s", fullPath))
	}

	if err = w.addRecursive(fullPath); err != nil {
		return
	}

	w.watches = w.fw.WatchList()
	return
}

func (w *Watcher) addRecursive(fullPath string) (err error) {
	e := FileWalk(fullPath, WalkFilter(w.filter), WithSkipFile, Walk(func(subPath string) {
		if !lo.ContainsBy(w.watches, func(p string) bool { return strings.HasPrefix(subPath, p) }) {
			err = errors.Join(err, w.fw.Add(subPath))
		}
	}))
	if e != nil {
		err = errors.Join(err, e)
	}
	return
}

type Option func(w *Watcher)

func Root(root ...string) Option {
	return func(o *Watcher) {
		o.root = lo.FilterMap(root, func(root string, _ int) (string, bool) {
			root, err := filepath.Abs(root)
			return root, err == nil
		})
	}
}

func Filter(filter ...string) Option {
	return func(o *Watcher) {
		o.filter = rex.CompileMatch(filter...)
	}
}

func AllowOp(allowOp string) Option {
	return func(o *Watcher) {
		o.allowOp = convOp(allowOp)
	}
}

func Throttle(throttle time.Duration) Option {
	return func(o *Watcher) {
		o.throttle = max(throttle, time.Second*3)
	}
}

type (
	HandlerFunc   func(ctx context.Context, events []fsnotify.Event)
	HandlerOption func(*Route)
	Route         struct {
		Name    string
		Match   func(string) bool
		Op      fsnotify.Op
		Handler HandlerFunc

		timer  *time.Timer
		events []fsnotify.Event

		mu sync.Mutex
	}
)

func (r *Route) Run(ctx context.Context) {
	r.mu.Lock()
	events := r.events
	r.events = r.events[:0]
	r.mu.Unlock()

	if len(events) == 0 {
		return
	}

	events = lo.UniqBy(events, func(item fsnotify.Event) string { return item.String() })
	slog.Debug(
		"route run",
		"name", r.Name,
		"event", lo.Reduce(events, func(acc fsnotify.Op, item fsnotify.Event, _ int) fsnotify.Op { return acc | item.Op }, 0).String(),
		"path", lo.Map(events, func(it fsnotify.Event, _ int) string { return it.Name }))

	if r.Handler == nil {
		return
	}

	r.Handler(ctx, events)
}

func Name(name string) HandlerOption {
	return func(r *Route) {
		r.Name = name
	}
}

func Match(match ...string) HandlerOption {
	return func(r *Route) {
		r.Match = rex.CompileMatch(match...)
	}
}

func Handle(handler HandlerFunc) HandlerOption {
	return func(r *Route) {
		r.Handler = handler
	}
}

func Events(eventOp string) HandlerOption {
	return func(r *Route) {
		r.Op = convOp(eventOp)
	}
}

func convOp(sop string) fsnotify.Op {
	if sop != "" {
		return lo.Reduce([]rune(sop), func(agg fsnotify.Op, item rune, _ int) fsnotify.Op {
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
	return fsnotify.Create | fsnotify.Write | fsnotify.Remove | fsnotify.Rename
}
