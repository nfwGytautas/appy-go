package appy_tracker

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

type Tracker struct {
	dsn        string
	sampleRate float32

	options TrackerOptions

	currentFlush int
}

func (s *Tracker) Flush() {
	if s.currentFlush == 0 {
		sentry.Flush(5 * time.Second)
		s.currentFlush = s.options.FlushInterval
		return
	}
	s.currentFlush -= 1
}

func (s *Tracker) ForceFlush() {
	sentry.Flush(5 * time.Second)
}

func (s *Tracker) OpenScope(name string) *Scope {
	scope := newScope(s)
	scope.SetTag("name", name)
	return scope
}

func (s *Tracker) OpenTransaction(ctx context.Context, name string) *Transaction {
	return newTransaction(s, ctx, name)
}
