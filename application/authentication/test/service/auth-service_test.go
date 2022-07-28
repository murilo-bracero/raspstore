package service

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/pb"
	mg "raspstore.github.io/authentication/repository"
	sv "raspstore.github.io/authentication/service"
	"raspstore.github.io/authentication/token"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	as, conn, cr := bootstrap(ctx)

	defer conn.Close(ctx)

	cred := createCredential()

	cr.Save(cred)

	loginRequest := &pb.LoginRequest{
		Email:    cred.Email,
		Password: "testpassword",
	}

	res, err := as.Login(ctx, loginRequest)

	assert.NoError(t, err)

	assert.NotEmpty(t, res.Token)
}

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	as, conn, cr := bootstrap(ctx)

	defer conn.Close(ctx)

	cred := createCredential()

	cr.Save(cred)

	loginRequest := &pb.LoginRequest{
		Email:    cred.Email,
		Password: "testpassword",
		MfaToken: "",
	}

	res, err := as.Login(ctx, loginRequest)

	assert.NoError(t, err)

	assert.NotEmpty(t, res.Token)

	tokenReq := &pb.AuthenticateRequest{Token: res.Token}

	tokenRes, err := as.Authenticate(ctx, tokenReq)

	assert.NoError(t, err)

	assert.NotEmpty(t, tokenRes.Uid)
}

func createCredential() *model.Credential {
	id := uuid.NewString()

	pass, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)

	return &model.Credential{
		Id:            primitive.NewObjectID(),
		UserId:        id,
		Email:         id + "_test@email.com",
		Secret:        "",
		Hash:          string(pass[:]),
		Has2FAEnabled: false,
	}
}

func bootstrap(ctx context.Context) (pb.AuthServiceServer, db.MongoConnection, mg.CredentialsRepository) {

	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(ctx, cfg)

	if err != nil {
		log.Panicln(err)
	}

	credRepo := mg.NewCredentialsRepository(ctx, conn)
	tokenManager := token.NewTokenManager(cfg)
	return sv.NewAuthService(credRepo, tokenManager), conn, credRepo
}
