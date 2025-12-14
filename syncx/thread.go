package syncx

import (
	"context"
	"sync"
)

func ContextGroup(ctx context.Context) interface {
	Go(f func(context.Context) error)
	Run(f func(context.Context) error)
	Wait() error
} {
	g := &contextGroup{}
	g.ctx, g.cancel = context.WithCancelCause(ctx)
	return g
}

type contextGroup struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel func(error)
	once   sync.Once
}

func (g *contextGroup) Go(f func(context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(g.ctx); err != nil {
			g.once.Do(func() { g.cancel(err) })
			return
		}
	}()

	// go1.25
	// g.wg.Go(func() {
	// 	if err := f(g.ctx); err != nil {
	// 		g.once.Do(func() { g.cancel(err) })
	// 		return
	// 	}
	// })
}

func (g *contextGroup) Run(f func(context.Context) error) {
	select {
	case <-g.ctx.Done():
		return
	default:
		if err := f(g.ctx); err != nil {
			g.once.Do(func() { g.cancel(err) })
			return
		}
	}
}

func (g *contextGroup) Wait() error {
	g.wg.Wait()
	return context.Cause(g.ctx)
}

func AllDone(done ...<-chan struct{}) {
	var wg sync.WaitGroup
	for _, d := range done {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-d
		}()
	}
	wg.Wait()
}
