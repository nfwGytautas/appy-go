package appy_tracker

import (
	"github.com/getsentry/sentry-go"
)

type Transaction struct {
	tx *sentry.Span
}

func (t *Transaction) Span(name string) *Transaction {
	return &Transaction{
		tx: t.tx.StartChild("default"),
	}
}

func (t *Transaction) Finish() {
	t.tx.Finish()
}
