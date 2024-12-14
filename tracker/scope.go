package appy_tracker

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func (t *Tracker) SetTag(key, value string) {
	t.scope.SetTag(key, value)
}

func (t *Tracker) SetContext(key string, value map[string]interface{}) {
	t.scope.SetContext(key, value)
}

func (t *Tracker) SetUser(id uint64, username string) {
	t.scope.SetUser(sentry.User{
		ID:       fmt.Sprintf("%v", id),
		Username: username,
	})
}

func (t *Tracker) SetRequest(r *http.Request) {
	t.scope.SetRequest(r)
}

func (t *Tracker) AddBreadcrumb(message, category string) {
	t.scope.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Level:    sentry.LevelInfo,
	}, 0)
}

func (t *Tracker) AddWarning(message, category string) {
	t.scope.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Level:    sentry.LevelWarning,
	}, 0)
}
