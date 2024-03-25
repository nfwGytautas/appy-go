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

func Initialize(options TrackerOptions) error {
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
