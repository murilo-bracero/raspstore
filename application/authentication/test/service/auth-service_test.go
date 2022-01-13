package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/pb"
	mg "raspstore.github.io/authentication/repository"
	sv "raspstore.github.io/authentication/service"
	"raspstore.github.io/authentication/token"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(ctx, cfg)
	assert.NoError(t, err)
	defer conn.Close(ctx)
	userRepo := mg.NewUsersRepository(ctx, conn)
	credRepo := mg.NewCredentialsRepository(ctx, conn)
	tokenManager := token.NewTokenManager(cfg)
	as := sv.NewAuthService(userRepo, credRepo, tokenManager)
	us := sv.NewUserService(userRepo, credRepo)

	createUserRequest := &pb.CreateUserRequest{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361320",
		Password:    "penispintorola212",
	}

	_, errService := us.CreateUser(ctx, createUserRequest)

	assert.NoError(t, errService)

	loginRequest := &pb.LoginRequest{
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	}

	res, err := as.Login(ctx, loginRequest)

	assert.NoError(t, err)

	assert.NotEmpty(t, res.Token)
}

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(ctx, cfg)
	assert.NoError(t, err)
	defer conn.Close(ctx)
	userRepo := mg.NewUsersRepository(ctx, conn)
	credRepo := mg.NewCredentialsRepository(ctx, conn)
	tokenManager := token.NewTokenManager(cfg)
	as := sv.NewAuthService(userRepo, credRepo, tokenManager)
	us := sv.NewUserService(userRepo, credRepo)

	createUserRequest := &pb.CreateUserRequest{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361320",
		Password:    "penispintorola212",
	}

	_, errService := us.CreateUser(ctx, createUserRequest)

	assert.NoError(t, errService)

	loginRequest := &pb.LoginRequest{
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	}

	res, err := as.Login(ctx, loginRequest)

	assert.NoError(t, err)

	assert.NotEmpty(t, res.Token)

	tokenReq := &pb.AuthenticateRequest{Token: res.Token}

	tokenRes, err := as.Authenticate(ctx, tokenReq)

	assert.NoError(t, err)

	assert.NotEmpty(t, tokenRes.Uid)
}
