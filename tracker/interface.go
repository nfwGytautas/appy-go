package appy_tracker

import (
	"github.com/getsentry/sentry-go"
)

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
	tracker = Tracker{
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

func Get() *Tracker {
	return &tracker
}
