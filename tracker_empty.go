package appy

import (
	"context"
	"net/http"
)

type nilTracker struct{}

func (n *nilTracker) Initialize(a *Appy, options TrackerOptions) error {
	return nil
}

func (n *nilTracker) OpenScope(name string) TrackerScope {
	return &nilTrackerScope{}
}

func (n *nilTracker) OpenTransaction(ctx context.Context, name string) TrackerTransaction {
	return &nilTrackerTransaction{}
}

func (n *nilTracker) Flush() {
}

func (n *nilTracker) ForceFlush() {
}

type nilTrackerScope struct{}

func (n *nilTrackerScope) SetTag(key, value string) {
}

func (n *nilTrackerScope) SetContext(key string, value map[string]interface{}) {
}

func (n *nilTrackerScope) SetRequest(r *http.Request) {
}

func (n *nilTrackerScope) AddBreadcrumb(message, category string) {
}

func (n *nilTrackerScope) AddWarning(message, category string) {
}

func (n *nilTrackerScope) CaptureError(err error) {
}

type nilTrackerTransaction struct{}

func (n *nilTrackerTransaction) Span(name string) TrackerTransaction {
	return &nilTrackerTransaction{}
}

func (n *nilTrackerTransaction) Finish() {
}
