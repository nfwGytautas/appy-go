package appy_driver_tracker

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/nfwGytautas/appy"
)

type sentryTransaction struct {
	tracker *sentryTracker
	tx      *sentry.Span
}

func newTransaction(tracker *sentryTracker, ctx context.Context, name string) *sentryTransaction {
	return &sentryTransaction{
		tracker: tracker,
		tx:      sentry.StartTransaction(ctx, name),
	}
}

func (t *sentryTransaction) Span(name string) appy.TrackerTransaction {
	return &sentryTransaction{
		tracker: t.tracker,
		tx:      t.tx.StartChild("default"),
	}
}

func (t *sentryTransaction) Finish() {
	t.tx.Finish()
}
