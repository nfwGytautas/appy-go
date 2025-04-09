package appy_mail

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

func (mw *mailerliteWrapper) CreateSubscriber(ctx context.Context, email string, fields MailerliteFields) (string, error) {
	subscriber := &mailerlite.Subscriber{
		Email:  email,
		Fields: fields,
	}

	sub, res, err := mw.client.Subscriber.Create(ctx, subscriber)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", errors.New("failed to create subscriber (mailerlite bad response): " + res.Status)
	}

	return sub.Data.ID, nil
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

func (mw *mailerliteWrapper) CreateGroup(ctx context.Context, name string) (*mailerlite.Group, error) {
	group, res, err := mw.client.Group.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New("failed to create group (mailerlite bad response): " + res.Status)
	}

	return &group.Data, nil
}

func (mw *mailerliteWrapper) GetGroups(ctx context.Context) ([]mailerlite.Group, error) {
	listOptions := &mailerlite.ListGroupOptions{
		Page:  1,
		Limit: 1000,
		Sort:  mailerlite.SortByName,
	}

	groups, res, err := mw.client.Group.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New("failed to get groups (mailerlite bad response): " + res.Status)
	}

	return groups.Data, nil
}

func (mw *mailerliteWrapper) AddToGroup(ctx context.Context, groupId string, subscriber string) error {
	_, res, err := mw.client.Group.Assign(ctx, groupId, subscriber)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("failed to add subscriber to group (mailerlite bad response): " + res.Status)
	}

	return nil
}

func (mw *mailerliteWrapper) RemoveFromGroup(ctx context.Context, groupId string, subscriber string) error {
	res, err := mw.client.Group.UnAssign(ctx, groupId, subscriber)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("failed to remove subscriber from group (mailerlite bad response): " + res.Status)
	}

	return nil
}
