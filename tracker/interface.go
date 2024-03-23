package appy_tracker

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// A tracker is a in-house or 3rd party service used to monitor and track the app health, e.g. Sentry
type Tracker interface {
	// Open a scope with a given name
	OpenScope(string) TrackerScope
	OpenTransaction(context.Context, string) TrackerTransaction

	Flush()
	ForceFlush()
}

// A scope
type TrackerScope interface {
	SetTag(string, string)
	SetContext(string, map[string]interface{})
	SetUser(string, string)

	SetRequest(*http.Request)

	AddBreadcrumb(string, string)
	AddWarning(string, string)

	CaptureError(error)
}

// Transaction entry for a tracker
type TrackerTransaction interface {
	Span(string) TrackerTransaction
	Finish()
}

// Options to describe a tracker
type TrackerOptions struct {
	// URL of the tracker
	DSN string

	// Sample rate of the tracker
	SampleRate float32

	// Number of request to wait before flushing
	FlushInterval int
}

var tracker Tracker

func Initialize(options TrackerOptions) error {
	tracker = &sentryTracker{
		dsn:        options.DSN,
		sampleRate: options.SampleRate,
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              options.DSN,
		EnableTracing:    true,
		TracesSampleRate: float64(options.SampleRate),
	})
	if err != nil {
		return err
	}

	return nil
}

func Get() Tracker {
	return tracker
}
