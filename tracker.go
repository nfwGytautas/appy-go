package appy

import (
	"context"
	"net/http"
)

// A tracker is a in-housr or 3rd party service used to monitor and track the app health, e.g. Sentry
type Tracker interface {
	Initialize(*Appy, TrackerOptions) error

	// Open a scope with a given name
	OpenScope(string) TrackerScope
	OpenTransaction(context.Context, string) TrackerTransaction

	Flush()
	ForceFlush()
}

type TrackerOptions struct {
	Provider Tracker

	// Number of request to wait before flushing
	FlushInterval int
}

// A scope
type TrackerScope interface {
	SetTag(string, string)
	SetContext(string, map[string]interface{})

	SetRequest(*http.Request)

	AddBreadcrumb(string, string)
	AddWarning(string, string)

	CaptureError(error)
}

type TrackerTransaction interface {
	Span(string) TrackerTransaction
	Finish()
}
