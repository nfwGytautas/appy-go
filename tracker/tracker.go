package appy_tracker

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

func Begin(ctx context.Context, name string) (context.Context, Scope, Transaction) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}

	scope := Scope{
		scope: hub.Scope(),
	}

	scope.SetTag("name", name)

	tx := Transaction{
		tx: sentry.StartTransaction(ctx,
			"default",
		),
	}

	return tx.tx.Context(), scope, tx
}

func CaptureError(ctx context.Context, err error) {
	hub := sentry.GetHubFromContext(ctx)
	hub.CaptureException(err)
}

func Flush() {
	sentry.Flush(5 * time.Second)
}
