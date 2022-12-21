package consumer

import (
	"context"
	"sync"
)

type Queue[T any] struct {
	m  sync.Map
	ch chan T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		ch: make(chan T),
	}
}

func (q *Queue[T]) Put(v T) {
	_, loaded := q.m.LoadOrStore(v, nil)
	if !loaded {
		q.ch <- v
	}
}

func (q *Queue[T]) Subscribe(ctx context.Context, fn func(v T)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case v := <-q.ch:
				q.m.Delete(v)
				fn(v)
			}
		}
	}()
}
