package appy_tracker

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type Scope struct {
	scope *sentry.Scope
}

func (s *Scope) SetTag(key, value string) {
	s.scope.SetTag(key, value)
}

func (s *Scope) SetContext(key string, value map[string]interface{}) {
	s.scope.SetContext(key, value)
}

func (s *Scope) SetUser(id uint64, username string) {
	s.scope.SetUser(sentry.User{
		ID:       fmt.Sprintf("%v", id),
		Username: username,
	})
}

func (s *Scope) SetRequest(r *http.Request) {
	s.scope.SetRequest(r)
}

func (s *Scope) AddBreadcrumb(message, category string) {
	s.scope.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Level:    sentry.LevelInfo,
	}, 0)
}

func (s *Scope) AddWarning(message, category string) {
	s.scope.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Level:    sentry.LevelWarning,
	}, 0)
}
