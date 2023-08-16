package utils

import "sync"

func NewFifo[T any]() *Fifo[T] {
	return &Fifo[T]{
		mutex: &sync.Mutex{},
		items: make([]T, 0),
	}
}

type Fifo[T any] struct {
	mutex *sync.Mutex
	items []T
}

func (f *Fifo[T]) Push(items ...T) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.items = append(f.items, items...)
}

func (f *Fifo[T]) Pop() *T {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.items) == 0 {
		return nil
	}

	item := f.items[0]

	if len(f.items) == 1 {
		f.items = make([]T, 0)
		return &item
	}

	f.items = f.items[1:]
	return &item
}

func (f *Fifo[T]) Len() int {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return len(f.items)
}
