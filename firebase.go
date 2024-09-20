package appy

import (
	"context"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

const cFirebaseAuthEnvironmentKey = "FIREBASE_SERVICE_FILE"
const cFirebaseConnectionTimeout = 5 * time.Second

// FirebaseOptions are options for configuring the firebase connection
type FirebaseServicesOptions struct {
	// If true then messaging service (used for push notifications) will be configured
	PushNotifications bool

	// If true then firebase authentification will be used
	Auth bool
}

// FirebaseWrapper is a wrapper for firebase API
type firebaseWrapper struct {
	app *firebase.App

	messaging *messaging.Client
	auth      *auth.Client
}

var firebaseInstance firebaseWrapper = firebaseWrapper{}

// Get firebase app instance
func Firebase() *firebaseWrapper {
	return &firebaseInstance
}

// Configure firebase connection
func (fw *firebaseWrapper) Configure(opts FirebaseServicesOptions) error {
	// Get the auth key and decode it
	authKey, err := Environment().GetValue(cFirebaseAuthEnvironmentKey)
	if err != nil {
		return err
	}

	decodedKey, err := os.ReadFile(authKey)
	if err != nil {
		return err
	}

	// Create a timeout context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, cFirebaseConnectionTimeout)
	defer cancel()

	// Configure options
	firebaseOpts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	// Create app instance
	fw.app, err = firebase.NewApp(ctx, nil, firebaseOpts...)
	if err != nil {
		return err
	}

	// Now configure services
	err = fw.configureServices(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}

// configureServices configures the services for the firebase connection
func (fw *firebaseWrapper) configureServices(ctx context.Context, opts FirebaseServicesOptions) error {
	var err error

	// Push notifications
	if opts.PushNotifications {
		fw.messaging, err = fw.app.Messaging(ctx)
		if err != nil {
			return err
		}
	}

	// Auth
	if opts.Auth {
		fw.auth, err = fw.app.Auth(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// PushNotifications returns the native firebase messaging client for sending push notifications
func (fw *firebaseWrapper) PushNotifications() *messaging.Client {
	return fw.messaging
}

// Auth returns the native firebase auth client
func (fw *firebaseWrapper) Auth() *auth.Client {
	return fw.auth
}
