package appy_tracker

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

type Tracker struct {
	scope *sentry.Scope
	spans []*sentry.Span
}

func Begin(ctx context.Context, name string) (context.Context, Tracker) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}

	scope := hub.Scope()

	scope.SetTag("name", name)

	tx := sentry.StartTransaction(ctx,
		name,
	)

	return tx.Context(), Tracker{
		scope: scope,
		spans: []*sentry.Span{
			tx,
		},
	}
}

func (t *Tracker) Finish() {
	t.spans[0].Finish()
}

func CaptureError(ctx context.Context, err error) {
	hub := sentry.GetHubFromContext(ctx)
	hub.CaptureException(err)
}

func Flush() {
	sentry.Flush(5 * time.Second)
}
