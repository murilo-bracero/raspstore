package grpc_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/auth-service/internal/grpc"
	"github.com/murilo-bracero/raspstore/auth-service/internal/infra"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"
	"github.com/stretchr/testify/assert"
)

var config *infra.Config

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}

	config = infra.NewConfig()
}

func TestAuthenticateSuccess(t *testing.T) {
	ctx := context.Background()
	as := grpc.NewAuthService(config)

	user := &model.User{
		UserId:      uuid.NewString(),
		Permissions: []string{},
	}

	token, err := token.Generate(config, user)

	assert.NoError(t, err)

	tokenReq := &pb.AuthenticateRequest{Token: fmt.Sprintf("Bearer %s", token)}

	tokenRes, err := as.Authenticate(ctx, tokenReq)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenRes.Uid)
	assert.Equal(t, user.UserId, tokenRes.Uid)
	assert.Equal(t, user.Permissions, tokenRes.Roles)
}

func TestAuthenticateFailWithInvalidToken(t *testing.T) {
	ctx := context.Background()
	as := grpc.NewAuthService(config)

	token := "testffailtokeninvalid"

	tokenReq := &pb.AuthenticateRequest{Token: fmt.Sprintf("Bearer %s", token)}

	_, err := as.Authenticate(ctx, tokenReq)

	assert.Error(t, err)
}

func TestAuthenticateFailWithEmptyToken(t *testing.T) {
	ctx := context.Background()
	as := grpc.NewAuthService(config)

	tokenReq := &pb.AuthenticateRequest{Token: ""}

	_, err := as.Authenticate(ctx, tokenReq)

	assert.Error(t, err)
}

func TestAuthenticateFailWithInsufficientParts(t *testing.T) {
	ctx := context.Background()
	as := grpc.NewAuthService(config)

	tokenReq := &pb.AuthenticateRequest{Token: "Bearer "}

	_, err := as.Authenticate(ctx, tokenReq)

	assert.Error(t, err)
}
