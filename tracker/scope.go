package appy_driver_tracker

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

type sentryScope struct {
	tracker *sentryTracker
	hub     *sentry.Hub
	scope   *sentry.Scope
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
