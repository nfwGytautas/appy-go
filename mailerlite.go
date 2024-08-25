package appy

import (
	"context"
	"errors"

	"github.com/mailerlite/mailerlite-go"
)

type MailerliteFields map[string]interface{}

// MailerLiteOptions are options for configuring the mailerlite connection
type MailerLiteOptions struct {
	Key string
}

type mailerliteWrapper struct {
	client *mailerlite.Client
}

var mailerliteInstance mailerliteWrapper = mailerliteWrapper{}

// Get mailerlite instance
func Mailerlite() *mailerliteWrapper {
	return &mailerliteInstance
}

// Configure mailerlite connection
func (mw *mailerliteWrapper) Configure(opts MailerLiteOptions) error {
	mw.client = mailerlite.NewClient(opts.Key)
	if mw.client == nil {
		return errors.New("failed to create mailerlite client")
	}

	return nil
}

func (mw *mailerliteWrapper) CreateSubscriber(ctx context.Context, email string, fields MailerliteFields) error {
	subscriber := &mailerlite.Subscriber{
		Email:  email,
		Fields: fields,
	}

	_, res, err := mw.client.Subscriber.Create(ctx, subscriber)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("failed to create subscriber (mailerlite bad response): " + res.Status)
	}

	return nil
}

func (mw *mailerliteWrapper) DeleteSubscriber(ctx context.Context, email string) error {
	subscriber, err := mw.GetSubscriber(ctx, email)
	if err != nil {
		return err
	}

	res, err := mw.client.Subscriber.Delete(ctx, subscriber.ID)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("failed to delete subscriber (mailerlite bad response): " + res.Status)
	}

	return nil
}

func (mw *mailerliteWrapper) GetSubscriber(ctx context.Context, email string) (*mailerlite.Subscriber, error) {
	getOptions := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, res, err := mw.client.Subscriber.Get(ctx, getOptions)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New("failed to get subscriber (mailerlite bad response): " + res.Status)
	}

	return &subscriber.Data, nil
}
