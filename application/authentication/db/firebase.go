package db

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseConnection interface {
	Client() *auth.Client
}

type firebaseConnection struct {
	client *auth.Client
}

func NewFirebaseConnection(ctx context.Context) (FirebaseConnection, error) {

	var (
		err    error
		app    *firebase.App
		client *auth.Client
	)

	app, err = firebase.NewApp(ctx, nil)

	if err != nil {
		return nil, err
	}

	if client, err = app.Auth(ctx); err != nil {
		return nil, err
	} else {
		return &firebaseConnection{client: client}, nil
	}
}

func (f *firebaseConnection) Client() *auth.Client {
	return f.client
}
