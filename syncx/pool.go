package syncx

import "sync"

type Pool[T any] interface {
	Get() *T
	Put(*T)
}

func NewPool[T any](ctor ...func() *T) Pool[T] {
	return &syncPool[T]{pool{
		New: func() any {
			if len(ctor) > 0 {
				return ctor[0]()
			}
			return new(T)
		},
	}}
}

type syncPool[T any] struct{ pool }

type pool = sync.Pool

func (p *syncPool[T]) Get() *T { return p.pool.Get().(*T) }

func (p *syncPool[T]) Put(o *T) {
	if reset, ok := any(o).(interface{ Reset() }); ok {
		reset.Reset()
	}
	p.pool.Put(o)
}
