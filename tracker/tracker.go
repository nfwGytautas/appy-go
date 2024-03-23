package appy_tracker

import (
	"context"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type sentryTracker struct {
	dsn        string
	sampleRate float32

	options TrackerOptions

	currentFlush int
}

type sentryScope struct {
	tracker *sentryTracker
	hub     *sentry.Hub
	scope   *sentry.Scope
}

type sentryTransaction struct {
	tracker *sentryTracker
	tx      *sentry.Span
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

func (s *sentryTracker) OpenScope(name string) TrackerScope {
	scope := newScope(s)
	scope.SetTag("name", name)
	return scope
}

func (s *sentryTracker) OpenTransaction(ctx context.Context, name string) TrackerTransaction {
	return newTransaction(s, ctx, name)
}

func newScope(tracker *sentryTracker) *sentryScope {
	hub := sentry.CurrentHub().Clone()
	return &sentryScope{
		tracker: tracker,
		hub:     hub,
		scope:   hub.Scope(),
	}
}

func (s *sentryScope) SetTag(key, value string) {
	s.scope.SetTag(key, value)
}

func (s *sentryScope) SetContext(key string, value map[string]interface{}) {
	s.scope.SetContext(key, value)
}

func (s *sentryScope) SetUser(id, username string) {
	s.scope.SetUser(sentry.User{
		ID:       id,
		Username: username,
	})
}

func (s *sentryScope) SetRequest(r *http.Request) {
	s.scope.SetRequest(r)
}

func (s *sentryScope) AddBreadcrumb(message, category string) {
	s.scope.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Level:    sentry.LevelInfo,
	}, 0)
}

func (s *sentryScope) AddWarning(message, category string) {
	s.scope.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Level:    sentry.LevelWarning,
	}, 0)
}

func (s *sentryScope) CaptureError(err error) {
	s.hub.CaptureException(err)
}

func newTransaction(tracker *sentryTracker, ctx context.Context, name string) *sentryTransaction {
	return &sentryTransaction{
		tracker: tracker,
		tx:      sentry.StartTransaction(ctx, name),
	}
}

func (t *sentryTransaction) Span(name string) TrackerTransaction {
	return &sentryTransaction{
		tracker: t.tracker,
		tx:      t.tx.StartChild("default"),
	}
}

func (t *sentryTransaction) Finish() {
	t.tx.Finish()
}
