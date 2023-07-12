package grpc_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/auth-service/internal/grpc"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/service"
	"github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	ctx := context.Background()
	as, tm := bootstrap(ctx)

	user := &model.User{
		UserId:      uuid.NewString(),
		Permissions: []string{},
	}

	token, err := tm.Generate(user)

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
	as, _ := bootstrap(ctx)

	token := "testffailtokeninvalid"

	tokenReq := &pb.AuthenticateRequest{Token: fmt.Sprintf("Bearer %s", token)}

	_, err := as.Authenticate(ctx, tokenReq)

	assert.Error(t, err)
}

func TestAuthenticateFailWithEmptyToken(t *testing.T) {
	ctx := context.Background()
	as, _ := bootstrap(ctx)

	tokenReq := &pb.AuthenticateRequest{Token: ""}

	_, err := as.Authenticate(ctx, tokenReq)

	assert.Error(t, err)
}

func TestAuthenticateFailWithInsufficientParts(t *testing.T) {
	ctx := context.Background()
	as, _ := bootstrap(ctx)

	tokenReq := &pb.AuthenticateRequest{Token: "Bearer "}

	_, err := as.Authenticate(ctx, tokenReq)

	assert.Error(t, err)
}

func bootstrap(ctx context.Context) (pb.AuthServiceServer, service.TokenService) {
	ts := service.NewTokenService()
	return grpc.NewAuthService(ts), ts
}
