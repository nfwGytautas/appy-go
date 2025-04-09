package appy_hooks

import "sync"

type EventFunc[T any] func(T)

type EventChannel[T any] struct {
	mtx         sync.Mutex
	subscribers []EventFunc[T]
}

func NewEventChannel[T any]() *EventChannel[T] {
	return &EventChannel[T]{
		mtx: sync.Mutex{},
	}
}

func (ec *EventChannel[T]) Subscribe(subscriber EventFunc[T]) {
	ec.subscribers = append(ec.subscribers, subscriber)
}

func (ec *EventChannel[T]) Publish(event T) {
	ec.mtx.Lock()
	defer ec.mtx.Unlock()

	for _, subscriber := range ec.subscribers {
		subscriber(event)
	}
}

func (ec *EventChannel[T]) PublishAsync(event T) {
	go ec.Publish(event)
}
