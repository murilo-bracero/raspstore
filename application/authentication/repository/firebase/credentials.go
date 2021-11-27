package repository

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
	rp "raspstore.github.io/authentication/repository"
)

type fireCredentialsRepository struct {
	ctx    context.Context
	client *auth.Client
}

func NewFireCredentials(ctx context.Context, conn db.FirebaseConnection) rp.CredentialsRepository {
	return &fireCredentialsRepository{ctx: ctx, client: conn.Client()}
}

func (f *fireCredentialsRepository) Save(user *model.User, password string) error {
	params := (&auth.UserToCreate{}).
		UID(user.UserId).
		Email(user.Email).
		EmailVerified(false).
		PhoneNumber(user.PhoneNumber).
		Password(password).
		DisplayName(user.Username).
		Disabled(false)

	_, err := f.client.CreateUser(f.ctx, params)

	return err
}

func (f *fireCredentialsRepository) Update(user *model.User) error {
	params := (&auth.UserToUpdate{}).
		Email(user.Email).
		PhoneNumber(user.PhoneNumber).
		DisplayName(user.Username)

	_, err := f.client.UpdateUser(f.ctx, user.UserId, params)

	return err
}

func (f *fireCredentialsRepository) Authenticate(token string) (string, error) {
	idToken, err := f.client.VerifyIDTokenAndCheckRevoked(f.ctx, token)

	if err != nil {
		return "", err
	}

	return idToken.UID, nil
}

func (f *fireCredentialsRepository) Delete(id string) error {
	return f.client.DeleteUser(f.ctx, id)
}
