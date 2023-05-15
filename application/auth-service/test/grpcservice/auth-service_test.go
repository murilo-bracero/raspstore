package service

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"github.com/stretchr/testify/assert"
	gs "raspstore.github.io/auth-service/grpcservice"
	"raspstore.github.io/auth-service/token"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	ctx := context.Background()
	as, tm := bootstrap(ctx)

	id := uuid.NewString()

	token, err := tm.Generate(id)

	assert.NoError(t, err)

	tokenReq := &pb.AuthenticateRequest{Token: fmt.Sprintf("Bearer %s", token)}

	tokenRes, err := as.Authenticate(ctx, tokenReq)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenRes.Uid)
	assert.Equal(t, id, tokenRes.Uid)
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

func bootstrap(ctx context.Context) (pb.AuthServiceServer, token.TokenManager) {
	tokenManager := token.NewTokenManager()
	return gs.NewAuthService(tokenManager), tokenManager
}
