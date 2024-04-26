package appy_tracker

import "github.com/getsentry/sentry-go"

func (t *Tracker) PushSpan(name string) {
	t.spans = append(
		t.spans,
		t.spans[len(t.spans)-1].StartChild(
			"default",
			sentry.WithTransactionName(name),
		),
	)
}

func (t *Tracker) PopSpan() {
	// Leave root
	if len(t.spans) == 1 {
		return
	}

	t.spans[len(t.spans)-1].Finish()
	t.spans = t.spans[:len(t.spans)-1]
}
