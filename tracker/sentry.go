package appy_driver_tracker

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/nfwGytautas/appy"
)

type sentryTracker struct {
	app        *appy.Appy
	dsn        string
	sampleRate float32

	options appy.TrackerOptions

	currentFlush int
}

func NewTracker(dsn string, sampleRate float32) appy.Tracker {
	return &sentryTracker{
		dsn:        dsn,
		sampleRate: sampleRate,
	}
}

func (s *sentryTracker) Initialize(a *appy.Appy, options appy.TrackerOptions) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              s.dsn,
		EnableTracing:    true,
		TracesSampleRate: float64(s.sampleRate),
	})
	if err != nil {
		return err
	}

	s.app = a
	s.options = options

	return nil
}

func (s *sentryTracker) Flush() {
	if s.currentFlush == 0 {
		sentry.Flush(5 * time.Second)
		s.currentFlush = s.options.FlushInterval
		return
	}
	s.currentFlush -= 1
}

func (s *sentryTracker) ForceFlush() {
	sentry.Flush(5 * time.Second)
}

func (s *sentryTracker) OpenScope(name string) appy.TrackerScope {
	scope := newScope(s)
	scope.SetTag("name", name)
	return scope
}

func (s *sentryTracker) OpenTransaction(ctx context.Context, name string) appy.TrackerTransaction {
	return newTransaction(s, ctx, name)
}
