package appy_tracker

import (
	"context"

	"github.com/getsentry/sentry-go"
)

type Transaction struct {
	tracker *Tracker
	tx      *sentry.Span
}

func newTransaction(tracker *Tracker, ctx context.Context, name string) *Transaction {
	return &Transaction{
		tracker: tracker,
		tx:      sentry.StartTransaction(ctx, name),
	}
}

func (t *Transaction) Span(name string) *Transaction {
	return &Transaction{
		tracker: t.tracker,
		tx:      t.tx.StartChild("default"),
	}
}

func (t *Transaction) Finish() {
	t.tx.Finish()
}
